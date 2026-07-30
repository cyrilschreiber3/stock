package routes

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func (g *RouteGroup) FullPrefix() string {
	if g == nil || g.Parent == nil {
		return g.Prefix
	}

	joinedPrefix := joinURL(g.Parent.FullPrefix(), g.Prefix)

	if !strings.HasPrefix(joinedPrefix, "/") {
		joinedPrefix = "/" + joinedPrefix
	}

	return joinedPrefix
}

func (g *RouteGroup) FullMiddlewares() []gin.HandlerFunc {
	if g == nil || g.Parent == nil {
		return g.Middlewares
	}
	return append(g.Parent.FullMiddlewares(), g.Middlewares...)
}

func (r Route[P]) RouteName() string {
	return r.Name
}

func (r StaticRoute) RouteName() string {
	return r.Route.Name
}

func (r Route[P]) MethodType() Method {
	return r.Method
}

func (r StaticRoute) MethodType() Method {
	return r.Route.Method
}

func (r Route[P]) ParamSpecs() []ParamSpec {
	return r.Params
}

func (r StaticRoute) ParamSpecs() []ParamSpec {
	return r.Route.Params
}

func (r Route[P]) Pattern() string {
	fullPattern := strings.TrimSpace(r.Path)
	fullPattern = strings.TrimPrefix(fullPattern, "/")

	if fullPattern == "" {
		fullPattern = "/"
	}

	if r.Group != nil {
		fullPattern = joinURL(r.Group.FullPrefix(), fullPattern)
	}

	if !strings.HasPrefix(fullPattern, "/") {
		fullPattern = "/" + fullPattern
	}

	return fullPattern
}

func (r StaticRoute) Pattern() string {
	return r.Route.Pattern()
}

func (r Route[P]) PatternWithParams(params map[string]string) string {
	pattern := r.Pattern()
	result := replaceParams(pattern, params)
	return result
}

func (r Route[P]) PatternWithQuery(queryParams url.Values) string {
	pattern := r.Pattern()
	if len(queryParams) > 0 {
		pattern = fmt.Sprintf("%s?%s", pattern, queryParams.Encode())
	}
	return pattern
}

func (r StaticRoute) PatternWithQuery(queryParams url.Values) string {
	return r.Route.PatternWithQuery(queryParams)
}

func (r Route[P]) URLSafe(p P) (string, error) {
	raw := r.Pattern()
	url := replaceParams(raw, p.Values())

	if strings.Contains(url, ":") {
		return "", fmt.Errorf("unresolved params in route %s: %s", r.Path, url)
	}
	return url, nil
}

func (r Route[P]) URLSafeWithQuery(p P, queryParams url.Values) (string, error) {
	url, err := r.URLSafe(p)
	if err != nil {
		return "", err
	}

	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, queryParams.Encode())
	}

	return url, nil
}

func (r Route[P]) URL(p P) string {
	url, err := r.URLSafe(p)
	if err != nil {
		panic(err)
	}
	return url
}

func (r StaticRoute) URL() string {
	return r.Route.URL(NoParams{})
}

func (r Route[P]) URLWithQuery(p P, queryParams url.Values) string {
	url, err := r.URLSafeWithQuery(p, queryParams)
	if err != nil {
		panic(err)
	}
	return url
}

func (r StaticRoute) URLWithQuery(queryParams url.Values) string {
	return r.Route.URLWithQuery(NoParams{}, queryParams)
}

func (r Route[P]) URLWithReturnToCurrent(p P, c *gin.Context) string {
	queryParams := url.Values{}
	queryParams.Set("from", c.Request.URL.Path)
	return r.URLWithQuery(p, queryParams)
}

func (r StaticRoute) URLWithReturnToCurrent(c *gin.Context) string {
	return r.Route.URLWithReturnToCurrent(NoParams{}, c)
}

func (r Route[P]) URLWithReturnToCurrentWithQuery(p P, c *gin.Context, additionalQueryParams url.Values) string {
	queryParams := url.Values{}
	queryParams.Set("from", c.Request.URL.Path)

	for key, values := range additionalQueryParams {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	return r.URLWithQuery(p, queryParams)
}

func (r StaticRoute) URLWithReturnToCurrentWithQuery(c *gin.Context, additionalQueryParams url.Values) string {
	return r.Route.URLWithReturnToCurrentWithQuery(NoParams{}, c, additionalQueryParams)
}

func (r Route[P]) URLWithReturn(p P, c *gin.Context) string {
	redirectUrl := utils.ResolveReturnPath(c, "")
	if redirectUrl != "" {
		return r.URLWithQuery(p, url.Values{"from": []string{redirectUrl}})
	}
	return r.URL(p)
}

func (r StaticRoute) URLWithReturn(c *gin.Context) string {
	return r.Route.URLWithReturn(NoParams{}, c)
}

func (r Route[P]) URLWithReturnWithQuery(p P, c *gin.Context, additionalQueryParams url.Values) string {
	redirectUrl := utils.ResolveReturnPath(c, "")
	queryParams := url.Values{}

	if redirectUrl != "" {
		queryParams.Set("from", redirectUrl)
	}

	for key, values := range additionalQueryParams {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	return r.URLWithQuery(p, queryParams)
}

func (r StaticRoute) URLWithReturnWithQuery(c *gin.Context, additionalQueryParams url.Values) string {
	return r.Route.URLWithReturnWithQuery(NoParams{}, c, additionalQueryParams)
}

func (r Route[P]) ReturnOrURL(p P, c *gin.Context) string {
	return utils.ResolveReturnPath(c, r.URL(p))
}

func (r StaticRoute) ReturnOrURL(c *gin.Context) string {
	return r.Route.ReturnOrURL(NoParams{}, c)
}

func (r Route[P]) ReturnOrURLWithQuery(p P, c *gin.Context, additionalQueryParams url.Values) string {
	return utils.ResolveReturnPath(c, r.URLWithQuery(p, additionalQueryParams))
}

func (r StaticRoute) ReturnOrURLWithQuery(c *gin.Context, additionalQueryParams url.Values) string {
	return r.Route.ReturnOrURLWithQuery(NoParams{}, c, additionalQueryParams)
}

func (r Route[P]) TestURLWithQuery(queryParams url.Values) (string, error) {
	testURL, err := r.TestURL()
	if err != nil {
		return "", err
	}
	if len(queryParams) > 0 {
		testURL = fmt.Sprintf("%s?%s", testURL, queryParams.Encode())
	}
	return testURL, nil
}

func (r StaticRoute) TestURLWithQuery(queryParams url.Values) (string, error) {
	return r.Route.TestURLWithQuery(queryParams)
}

func (r Route[P]) TestURL() (string, error) {
	testParams := make(map[string]string, len(r.Params))
	for _, param := range r.Params {
		switch param.Type {
		case ParamUUID:
			testParams[param.Name] = "00000000-0000-0000-0000-000000000000"
		case ParamText:
			testParams[param.Name] = "test"
		}
	}
	url := replaceParams(r.Pattern(), testParams)
	if strings.Contains(url, ":") {
		return "", fmt.Errorf("unresolved params in route %s: %s", r.Path, url)
	}
	return url, nil
}

func (r StaticRoute) TestURL() (string, error) {
	return r.Route.TestURL()
}

func (r Route[P]) TestMustURL() (string, error) {
	testURL, err := r.TestURL()
	if err != nil {
		return "", err
	}
	return testURL, nil
}

func (r StaticRoute) TestMustURL() (string, error) {
	return r.Route.TestMustURL()
}

func (r Route[P]) ZeroParamValues() map[string]string {
	var zero P
	return zero.Values()
}

func (r StaticRoute) ZeroParamValues() map[string]string {
	return r.Route.ZeroParamValues()
}

func (r Route[P]) Register(router gin.IRoutes, handler gin.HandlerFunc) {
	var handlers []gin.HandlerFunc
	handlersCount := len(r.Middlewares) + 1
	if r.Group != nil {
		handlersCount += len(r.Group.FullMiddlewares())
	}
	handlers = make([]gin.HandlerFunc, 0, handlersCount)
	if r.Group != nil {
		handlers = append(handlers, r.Group.FullMiddlewares()...)
	}
	handlers = append(handlers, r.Middlewares...)
	handlers = append(handlers, handler)

	router.Handle(string(r.Method), r.Pattern(), handlers...)
}

func (r StaticRoute) Register(router gin.IRoutes, handler gin.HandlerFunc) {
	r.Route.Register(router, handler)
}

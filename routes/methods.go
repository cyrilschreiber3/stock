package routes

import (
	"fmt"
	"net/url"
	"strings"

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

func (r Route[P]) MethodType() Method {
	return r.Method
}

func (r Route[P]) HandlerFunc() gin.HandlerFunc {
	return r.Handler
}

func (r Route[P]) ParamSpecs() []ParamSpec {
	return r.Params
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

func (r Route[P]) PatternWithQuery(queryParams url.Values) string {
	pattern := r.Pattern()
	if len(queryParams) > 0 {
		pattern = fmt.Sprintf("%s?%s", pattern, queryParams.Encode())
	}
	return pattern
}

func (r Route[P]) URL(p P) (string, error) {
	raw := r.Pattern()
	url := replaceParams(raw, p.Values())

	if strings.Contains(url, ":") {
		return "", fmt.Errorf("unresolved params in route %s: %s", r.Path, url)
	}
	return url, nil
}

func (r Route[P]) URLWithQuery(p P, queryParams url.Values) (string, error) {
	url, err := r.URL(p)
	if err != nil {
		return "", err
	}

	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, queryParams.Encode())
	}

	return url, nil
}

func (r Route[P]) MustURL(p P) string {
	url, err := r.URL(p)
	if err != nil {
		panic(err)
	}
	return url
}

func (r Route[P]) MustURLWithQuery(p P, queryParams url.Values) string {
	url, err := r.URLWithQuery(p, queryParams)
	if err != nil {
		panic(err)
	}
	return url
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

func (r Route[P]) TestMustURL() (string, error) {
	testURL, err := r.TestURL()
	if err != nil {
		return "", err
	}
	return testURL, nil
}

func (r Route[P]) ZeroParamValues() map[string]string {
	var zero P
	return zero.Values()
}

func (r Route[P]) Register(router gin.IRoutes) {
	handlers := make([]gin.HandlerFunc, 0, len(r.Group.FullMiddlewares())+len(r.Middlewares)+1)
	handlers = append(handlers, r.Group.FullMiddlewares()...)
	handlers = append(handlers, r.Middlewares...)
	handlers = append(handlers, r.Handler)

	router.Handle(string(r.Method), r.Pattern(), handlers...)
}

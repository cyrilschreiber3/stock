package routes

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

type Method string

const (
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodPUT    Method = "PUT"
	MethodPATCH  Method = "PATCH"
	MethodDELETE Method = "DELETE"
)

type ParamType string

const (
	ParamUUID ParamType = "uuid"
	ParamText ParamType = "text"
)

type ParamSpec struct {
	Name string
	Type ParamType
}

type RouteGroup struct {
	Prefix      string
	Middlewares []gin.HandlerFunc
	Parent      *RouteGroup
}

type RouteParams interface {
	Values() map[string]string
}

type Route[P RouteParams] struct {
	Name        string
	Method      Method
	Group       *RouteGroup
	Path        string
	Params      []ParamSpec
	Middlewares []gin.HandlerFunc
	Handler     gin.HandlerFunc
}

type RouteHandler interface {
	Register(router gin.IRoutes)
	RouteName() string
	MethodType() Method
	HandlerFunc() gin.HandlerFunc
	ParamSpecs() []ParamSpec
	Pattern() string
	PatternWithQuery(queryParams url.Values) string
	TestURL() (string, error)
	TestMustURL() (string, error)
	TestURLWithQuery(queryParams url.Values) (string, error)
	ZeroParamValues() map[string]string
}

type ResourceRoutes struct {
	Group  *RouteGroup
	Routes []RouteHandler
}

type GroupOpt func(*RouteGroup)

type RouteOpt[P RouteParams] func(*Route[P])

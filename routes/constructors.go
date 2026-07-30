package routes

import "github.com/gin-gonic/gin"

func Group(prefix string, opts ...GroupOpt) *RouteGroup {
	group := &RouteGroup{Prefix: prefix}
	for _, opt := range opts {
		opt(group)
	}
	return group
}

func WithParent(parent *RouteGroup) GroupOpt {
	return func(g *RouteGroup) {
		g.Parent = parent
	}
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) GroupOpt {
	return func(g *RouteGroup) {
		g.Middlewares = append(g.Middlewares, middlewares...)
	}
}

func Spec[P RouteParams](name string, method Method, group *RouteGroup, path string, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:   name,
		Method: method,
		Group:  group,
		Path:   path,
	}
	for _, opt := range opts {
		opt(&route)
	}
	return route
}

func Spec0(name string, method Method, group *RouteGroup, path string, opts ...RouteOpt[NoParams]) StaticRoute {
	route := StaticRoute{
		Route: Route[NoParams]{
			Name:   name,
			Method: method,
			Group:  group,
			Path:   path,
		},
	}
	for _, opt := range opts {
		opt(&route.Route)
	}
	return route
}

func WithRouteMiddlewares[P RouteParams](middlewares ...gin.HandlerFunc) RouteOpt[P] {
	return func(r *Route[P]) {
		r.Middlewares = append(r.Middlewares, middlewares...)
	}
}

func WithRouteParams[P RouteParams](params ...ParamSpec) RouteOpt[P] {
	return func(r *Route[P]) {
		r.Params = params
	}
}

func registerAllSpecsForTests(router gin.IRoutes, routes ...RouteHandler) {
	for _, route := range routes {
		route.Register(router, func(c *gin.Context) {})
	}
}

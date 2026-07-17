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

func GET[P RouteParams](name string, group *RouteGroup, path string, handler gin.HandlerFunc, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:    name,
		Method:  MethodGET,
		Group:   group,
		Path:    path,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(&route)
	}
	return route
}

func POST[P RouteParams](name string, group *RouteGroup, path string, handler gin.HandlerFunc, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:    name,
		Method:  MethodPOST,
		Group:   group,
		Path:    path,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(&route)
	}
	return route
}

func PUT[P RouteParams](name string, group *RouteGroup, path string, handler gin.HandlerFunc, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:    name,
		Method:  MethodPUT,
		Group:   group,
		Path:    path,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(&route)
	}
	return route
}

func PATCH[P RouteParams](name string, group *RouteGroup, path string, handler gin.HandlerFunc, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:    name,
		Method:  MethodPATCH,
		Group:   group,
		Path:    path,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(&route)
	}
	return route
}

func DELETE[P RouteParams](name string, group *RouteGroup, path string, handler gin.HandlerFunc, opts ...RouteOpt[P]) Route[P] {
	route := Route[P]{
		Name:    name,
		Method:  MethodDELETE,
		Group:   group,
		Path:    path,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(&route)
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

func RegisterAll(router gin.IRoutes, routes ...RouteHandler) {
	for _, route := range routes {
		route.Register(router)
	}
}

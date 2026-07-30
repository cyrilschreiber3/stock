package routes

var IndexGroup = Group("/")

var Index = Spec0("Index", MethodGET, IndexGroup, "/")

var IndexRoutes = ResourceRoutes{
	Group: IndexGroup,
	Routes: []RouteHandler{
		Index,
	},
}

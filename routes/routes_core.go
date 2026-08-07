package routes

var IndexGroup = Group("/")

var Index = Spec0("Index", MethodGET, IndexGroup, "/")

var LangRouteParams = WithRouteParams[LangParam](ParamLang)
var LangRoute = Spec("Lang", MethodPOST, IndexGroup, "/lang/:lang", LangRouteParams)

var IndexRoutes = ResourceRoutes{
	Group: IndexGroup,
	Routes: []RouteHandler{
		Index,
		LangRoute,
	},
}

package router

import (
	"testing"

	"github.com/cyrilschreiber3/stock/routes"
	"github.com/gin-gonic/gin"
)

func TestRouteRegistration(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	testRouter := gin.New()
	routes.RegisterSpecsForTests(testRouter)

	prodRouter := gin.New()
	RegisterRoutes(prodRouter)

	prodRouteTable := map[string]gin.RouteInfo{}
	for _, routeInfo := range prodRouter.Routes() {
		prodRouteTable[routeInfo.Method+"_"+routeInfo.Path] = routeInfo
	}

	testRouteTable := map[string]gin.RouteInfo{}
	t.Run("test routes registered in production router", func(t *testing.T) {
		for _, routeInfo := range testRouter.Routes() {
			key := routeInfo.Method + "_" + routeInfo.Path
			testRouteTable[key] = routeInfo
			if _, exists := prodRouteTable[key]; !exists {
				t.Errorf("route %s not registered in production router", key)
			}
		}
	})

	t.Run("test routes registered in test router", func(t *testing.T) {
		for key := range prodRouteTable {
			if _, exists := testRouteTable[key]; !exists {
				t.Errorf("route %s not registered in test router", key)
			}
		}
	})
}

package routes

import (
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouteNameUnique(t *testing.T) {
	routeNames := make(map[string]bool)
	allRoutes := GetAllRoutes()
	for _, route := range allRoutes {
		if routeNames[route.RouteName()] {
			t.Errorf("duplicate route name: %s", route.RouteName())
		}
		routeNames[route.RouteName()] = true
	}
}

func TestRouteMethodAndPatternUnique(t *testing.T) {
	methodAndPattern := make(map[string]bool)
	allRoutes := GetAllRoutes()
	for _, route := range allRoutes {
		key := string(route.MethodType()) + " " + route.Pattern()
		if methodAndPattern[key] {
			t.Errorf("duplicate route method and pattern: %s", key)
		}
		methodAndPattern[key] = true
	}
}

func testQueryParameters(t *testing.T, encodedURL string, queryParams url.Values) {
	t.Helper()

	if !strings.Contains(encodedURL, "?") && len(queryParams) > 0 {
		t.Errorf("expected query parameters in URL, but none found: %s", encodedURL)
	}

	if strings.Contains(encodedURL, "?") && len(queryParams) == 0 {
		t.Errorf("expected no query parameters in URL, but found some: %s", encodedURL)
	}

	parsedURL, err := url.Parse(encodedURL)
	if err != nil {
		t.Errorf("error parsing generated URL: %v", err)
		return
	}
	parsedQuery := parsedURL.Query()
	for key, values := range queryParams {
		parsedValues, ok := parsedQuery[key]
		if !ok || len(parsedValues) != len(values) {
			t.Errorf("query parameter %q missing or has incorrect number of values in generated URL", key)
			continue
		}
		for i, value := range values {
			if parsedValues[i] != value {
				t.Errorf("query parameter %q has incorrect value in generated URL: got %q, want %q", key, parsedValues[i], value)
			}
		}
	}
}

func getTestQueryParams() []struct {
	name        string
	queryParams url.Values
} {
	tests := []struct {
		name        string
		queryParams url.Values
	}{
		{
			name: "simple query parameters",
			queryParams: url.Values{
				"simple": []string{"value"},
			},
		},
		{
			name: "query parameters with spaces",
			queryParams: url.Values{
				"with spaces": []string{"hello world"},
			},
		},
		{
			name: "query parameters with special characters",
			queryParams: url.Values{
				"special": []string{"!@#$%^&*()"},
			},
		},
		{
			name: "query parameters with empty value",
			queryParams: url.Values{
				"empty": []string{""},
			},
		},
	}
	return tests
}

func TestURLQueryParameters(t *testing.T) {
	tests := getTestQueryParams()

	for _, route := range GetAllRoutes() {
		for _, tt := range tests {
			t.Run(route.RouteName()+"_"+tt.name, func(t *testing.T) {
				encodedURL, err := route.TestURLWithQuery(tt.queryParams)
				if err != nil {
					t.Errorf("error generating URL with query: %v", err)
					return
				}
				testQueryParameters(t, encodedURL, tt.queryParams)
			})
		}
	}
}

func TestPatternQueryParameters(t *testing.T) {
	tests := getTestQueryParams()

	for _, route := range GetAllRoutes() {
		for _, tt := range tests {
			t.Run(route.RouteName()+"_"+tt.name, func(t *testing.T) {
				pattern := route.PatternWithQuery(tt.queryParams)
				testQueryParameters(t, pattern, tt.queryParams)
			})
		}
	}
}

func TestRouteFullPrefix(t *testing.T) {
	testCases := []struct {
		name     string
		group    *RouteGroup
		expected string
	}{
		{
			name:     "single group",
			group:    &RouteGroup{Prefix: "/api"},
			expected: "/api",
		},
		{
			name:     "nested groups",
			group:    &RouteGroup{Prefix: "/v1", Parent: &RouteGroup{Prefix: "/api"}},
			expected: "/api/v1",
		},
		{
			name:     "multiple nested groups",
			group:    &RouteGroup{Prefix: "/v1", Parent: &RouteGroup{Prefix: "/api", Parent: &RouteGroup{Prefix: "/root"}}},
			expected: "/root/api/v1",
		},
		{
			name:     "no prefix",
			group:    &RouteGroup{Prefix: "", Parent: &RouteGroup{Prefix: "/api"}},
			expected: "/api/",
		},
		{
			name:     "with trailing slash",
			group:    &RouteGroup{Prefix: "/v1/", Parent: &RouteGroup{Prefix: "/api/"}},
			expected: "/api/v1",
		},
		{
			name:     "with no slashes",
			group:    &RouteGroup{Prefix: "v1", Parent: &RouteGroup{Prefix: "api"}},
			expected: "/api/v1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.group.FullPrefix()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestRoutePattern(t *testing.T) {
	testCases := []struct {
		name     string
		route    Route[NoParams]
		expected string
	}{
		{
			name:     "simple route",
			route:    Route[NoParams]{Path: "test", Group: &RouteGroup{Prefix: "/api"}},
			expected: "/api/test",
		},
		{
			name:     "simple route with no group",
			route:    Route[NoParams]{Path: "test", Group: nil},
			expected: "/test",
		},
		{
			name:     "nested group route",
			route:    Route[NoParams]{Path: "test", Group: &RouteGroup{Prefix: "/v1", Parent: &RouteGroup{Prefix: "/api"}}},
			expected: "/api/v1/test",
		},
		{
			name:     "route with leading slash",
			route:    Route[NoParams]{Path: "/test", Group: &RouteGroup{Prefix: "/api"}},
			expected: "/api/test",
		},
		{
			name:     "route with trailing slash",
			route:    Route[NoParams]{Path: "test/", Group: &RouteGroup{Prefix: "/api"}},
			expected: "/api/test",
		},
		{
			name:     "route with empty group prefix",
			route:    Route[NoParams]{Path: "test", Group: &RouteGroup{Prefix: ""}},
			expected: "/test",
		},
		{
			name:     "route with empty path",
			route:    Route[NoParams]{Path: "", Group: &RouteGroup{Prefix: "/api"}},
			expected: "/api/",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.route.Pattern()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestRouteFullMiddlewares(t *testing.T) {
	testCases := []struct {
		name     string
		group    *RouteGroup
		expected []string // We'll use strings to identify middlewares for testing
	}{
		{
			name: "single group",
			group: &RouteGroup{
				Prefix:      "/api",
				Middlewares: []gin.HandlerFunc{func(c *gin.Context) { c.Set("mw1", true) }},
			},
			expected: []string{"mw1"},
		},
		{
			name: "nested groups",
			group: &RouteGroup{
				Prefix:      "/v1",
				Middlewares: []gin.HandlerFunc{func(c *gin.Context) { c.Set("mw2", true) }},
				Parent: &RouteGroup{
					Prefix:      "/api",
					Middlewares: []gin.HandlerFunc{func(c *gin.Context) { c.Set("mw1", true) }},
				},
			},
			expected: []string{"mw1", "mw2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middlewares := tc.group.FullMiddlewares()
			if len(middlewares) != len(tc.expected) {
				t.Errorf("expected %d middlewares, got %d", len(tc.expected), len(middlewares))
				return
			}
			for i, mw := range middlewares {
				c := &gin.Context{}
				mw(c)
				if _, exists := c.Get(tc.expected[i]); !exists {
					t.Errorf("expected middleware %q to be executed", tc.expected[i])
				}
			}
		})
	}
}

func TestRouteRegistration(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	RegisterSpecsForTests(router)
	ginRouteTable := map[string]gin.RouteInfo{}
	for _, routeInfo := range router.Routes() {
		ginRouteTable[routeInfo.Method+"_"+routeInfo.Path] = routeInfo
	}

	for _, route := range GetAllRoutes() {
		routeKey := string(route.MethodType()) + "_" + route.Pattern()

		if _, exists := ginRouteTable[routeKey]; !exists {
			t.Errorf("route %s (%s) not found in gin route table", route.RouteName(), routeKey)
		}
	}
}

package routes

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func countPathParams(t *testing.T, pattern string) int {
	t.Helper()

	count := 0
	for _, segment := range strings.Split(pattern, "/") {
		if strings.HasPrefix(segment, ":") {
			count++
		}
	}
	return count
}

func TestRouteParamSpecCount(t *testing.T) {
	for _, route := range GetAllRoutes() {
		t.Run(route.RouteName(), func(t *testing.T) {
			pattern := route.Pattern()
			inPattern := countPathParams(t, pattern)
			inSpecs := len(route.ParamSpecs())

			if inPattern != inSpecs {
				t.Errorf("pattern %q has %d path params but ParamSpecs declares %d",
					pattern, inPattern, inSpecs)
			}
		})
	}
}

func TestRouteParamSpecNames(t *testing.T) {
	for _, route := range GetAllRoutes() {
		t.Run(route.RouteName(), func(t *testing.T) {
			pattern := route.Pattern()
			paramNamesInPattern := make(map[string]bool)
			for _, segment := range strings.Split(pattern, "/") {
				if strings.HasPrefix(segment, ":") {
					paramNamesInPattern[segment[1:]] = true
				}
			}

			paramNamesInParamSpecs := make(map[string]bool)
			for _, paramSpec := range route.ParamSpecs() {
				paramNamesInParamSpecs[paramSpec.Name] = true

				if !paramNamesInPattern[paramSpec.Name] {
					t.Errorf("ParamSpec %q is not present in pattern %q", paramSpec.Name, pattern)
				}
			}

			for paramName := range paramNamesInPattern {
				if !paramNamesInParamSpecs[paramName] {
					t.Errorf("pattern %q contains param %q which is not declared in ParamSpecs", pattern, paramName)
				}
			}
		})
	}
}

func TestRouteURLSubstitution(t *testing.T) {
	for _, route := range GetAllRoutes() {
		t.Run(route.RouteName(), func(t *testing.T) {
			url, err := route.TestURL()
			if err != nil {
				t.Error(err)
				return
			}
			if strings.Contains(url, ":") {
				t.Errorf("URL %q still contains unresolved params", url)
			}
		})
	}
}

func TestRouteMustURL(t *testing.T) {
	validRoute := Spec("ValidRoute", MethodGET, nil, "/test/:id", WithRouteParams[IDParam](ParamID))
	invalidRoute := Spec("InvalidRoute", MethodGET, nil, "/test/:id", WithRouteParams[NoParams]())

	t.Run("MustURL panics on unresolved params", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustURL did not panic on unresolved params")
			}
		}()
		invalidRoute.URL(NoParams{})
	})

	t.Run("MustURL does not panic on valid params", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustURL panicked on valid params: %v", r)
			}
		}()
		validRoute.URL(IDParam{ID: uuid.New()})
	})
}

func assertParamSpecsMatchRouteParamsType(t *testing.T, route RouteHandler) []error {
	t.Helper()

	errors := []error{}

	values := route.ZeroParamValues()
	routeParamsValuesKeys := make(map[string]struct{})
	for k := range values {
		routeParamsValuesKeys[k] = struct{}{}
	}

	paramSpecKeys := make(map[string]struct{})
	for _, spec := range route.ParamSpecs() {
		paramSpecKeys[spec.Name] = struct{}{}
	}

	for k := range paramSpecKeys {
		if _, ok := routeParamsValuesKeys[k]; !ok {
			errors = append(errors, fmt.Errorf("route %q: ParamSpec %q missing from RouteParams.Values()", route.RouteName(), k))
		}
	}

	for k := range routeParamsValuesKeys {
		if _, ok := paramSpecKeys[k]; !ok {
			errors = append(errors, fmt.Errorf("route %q: RouteParams.Values() key %q missing from ParamSpecs", route.RouteName(), k))
		}
	}
	return errors
}

func TestRouteParamSpecsMatchRouteParamsType(t *testing.T) {
	for _, route := range GetAllRoutes() {
		t.Run(route.RouteName(), func(t *testing.T) {
			errors := assertParamSpecsMatchRouteParamsType(t, route)
			for _, err := range errors {
				t.Error(err)
			}
		})
	}
}

// TODO: Check if there is other edge cases to test for
func TestAssertParamSpecsMatchRouteParamsType_CatchesMismatch(t *testing.T) {
	t.Run("missing values key for declared paramspec", func(t *testing.T) {
		route := Spec(
			"BadMissingField",
			MethodGET,
			nil,
			"/test/:id/:field",
			WithRouteParams[IDParam](ParamID, ParamField),
		)

		errors := assertParamSpecsMatchRouteParamsType(t, route)
		if len(errors) == 0 {
			t.Errorf("expected errors, but got none")
		}
	})

	t.Run("values key missing from paramspecs", func(t *testing.T) {
		route := Spec(
			"BadExtraValueKey",
			MethodGET,
			nil,
			"/test/:id/:field",
			WithRouteParams[ObjectFieldParams](ParamID),
		)

		errors := assertParamSpecsMatchRouteParamsType(t, route)
		if len(errors) == 0 {
			t.Errorf("expected errors, but got none")
		}
	})
}

func assertParamSpecTypesMatchValues(t *testing.T, route RouteHandler) []error {
	t.Helper()

	errors := []error{}

	values := route.ZeroParamValues()

	for _, spec := range route.ParamSpecs() {
		value, ok := values[spec.Name]
		if !ok {
			continue // name mismatch already caught by assertParamSpecsMatchRouteParamsType
		}
		switch spec.Type {
		case ParamUUID:
			if _, err := uuid.Parse(value); err != nil {
				errors = append(errors, fmt.Errorf("route %q: param %q declared as ParamUUID but Values() returns non-UUID value %q",
					route.RouteName(), spec.Name, value))
			}
		case ParamText:
			// no type constraint on text values
		default:
			errors = append(errors, fmt.Errorf("route %q: param %q has unknown ParamType %q",
				route.RouteName(), spec.Name, spec.Type))
		}
	}
	return errors
}

func TestRouteParamSpecTypesMatchValues(t *testing.T) {
	for _, route := range GetAllRoutes() {
		t.Run(route.RouteName(), func(t *testing.T) {
			errors := assertParamSpecTypesMatchValues(t, route)
			for _, err := range errors {
				t.Error(err)
			}
		})
	}
}

// FIXME: This test will catch if a param is defined as complex in the RouteParams but is declared as a simple type in the ParamSpecs, but it won't catch the reverse case.
func TestAssertParamSpecTypesMatchValues_CatchesMismatch(t *testing.T) {
	t.Run("uuid declared as text does not fail with current helper", func(t *testing.T) {
		route := Spec(
			"TextSpecForUUIDValue",
			MethodGET,
			nil,
			"/test/:id",
			WithRouteParams[IDParam](ParamSpec{Name: "id", Type: ParamText}),
		)

		errors := assertParamSpecTypesMatchValues(t, route)
		if len(errors) > 0 {
			t.Errorf("expected no errors, but got: %v", errors)
		}
	})

	t.Run("text declared as uuid does fail", func(t *testing.T) {
		route := Spec(
			"UUIDSpecForTextValue",
			MethodGET,
			nil,
			"/test/:field",
			WithRouteParams[ObjectFieldParams](ParamSpec{Name: "field", Type: ParamUUID}),
		)

		errors := assertParamSpecTypesMatchValues(t, route)
		if len(errors) == 0 {
			t.Errorf("expected errors, but got none")
		}
	})
}

package routes

import "github.com/google/uuid"

var (
	ParamID        = ParamSpec{Name: "id", Type: ParamUUID}
	ParamField     = ParamSpec{Name: "field", Type: ParamText}
	ParamProductID = ParamSpec{Name: "product_id", Type: ParamUUID}
	ParamName      = ParamSpec{Name: "name", Type: ParamText}
)

// Deprecated: Use the StaticRoute type instead of Route[NoParams]
type NoParams struct{}

func (p NoParams) Values() map[string]string {
	return map[string]string{}
}

type IDParam struct {
	ID uuid.UUID
}

func (p IDParam) Values() map[string]string {
	return map[string]string{"id": p.ID.String()}
}

func ID(id uuid.UUID) IDParam {
	return IDParam{ID: id}
}

type ObjectFieldParams struct {
	ID    uuid.UUID
	Field string
}

func (p ObjectFieldParams) Values() map[string]string {
	return map[string]string{"id": p.ID.String(), "field": p.Field}
}

func ObjectField(id uuid.UUID, field string) ObjectFieldParams {
	return ObjectFieldParams{ID: id, Field: field}
}

type NameParam struct {
	Name string
}

func (p NameParam) Values() map[string]string {
	return map[string]string{"name": p.Name}
}

func Name(name string) NameParam {
	return NameParam{Name: name}
}

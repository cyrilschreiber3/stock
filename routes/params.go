package routes

import "github.com/google/uuid"

var (
	ParamID        = ParamSpec{Name: "id", Type: ParamUUID}
	ParamField     = ParamSpec{Name: "field", Type: ParamText}
	ParamProductID = ParamSpec{Name: "product_id", Type: ParamUUID}
)

type NoParams struct{}

func (p NoParams) Values() map[string]string {
	return map[string]string{}
}

type SimpleUUIDParam struct {
	ID uuid.UUID
}

func (p SimpleUUIDParam) Values() map[string]string {
	return map[string]string{"id": p.ID.String()}
}

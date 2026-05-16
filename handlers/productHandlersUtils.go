package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func parseProductForm(c *gin.Context) (*productFormPayload, int, error) {
	var form productFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	categoryID, err := parseUUIDField(form.CategoryID, "category")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	subcategoryID, err := parseUUIDField(form.SubcategoryID, "subcategory")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	defaultSupplierID, err := parseUUIDField(form.DefaultSupplierID, "default supplier")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var buyPrice pgtype.Numeric
	if err := buyPrice.Scan(form.DefaultBuyPrice); err != nil {
		return nil, http.StatusBadRequest, err
	}

	var sellPrice pgtype.Numeric
	if err := sellPrice.Scan(form.DefaultSellPrice); err != nil {
		return nil, http.StatusBadRequest, err
	}

	aliases := make([]string, 0, len(form.Aliases))
	for _, alias := range form.Aliases {
		cleanAlias := strings.TrimSpace(alias)
		if cleanAlias != "" {
			aliases = append(aliases, cleanAlias)
		}
	}

	payload := productFormPayload{
		Brand:             form.Brand,
		Name:              form.Name,
		CategoryID:        categoryID,
		SubcategoryID:     subcategoryID,
		DefaultSupplierID: defaultSupplierID,
		DefaultBuyPrice:   buyPrice,
		DefaultSellPrice:  sellPrice,
		Aliases:           aliases,
	}

	return &payload, http.StatusOK, nil
}

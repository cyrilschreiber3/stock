package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func parseUUIDField(raw string, field string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s ID: %w", field, err)
	}

	return parsed, nil
}

func parseUUIDParam(c *gin.Context, paramName string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(c.Param(paramName)))
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s ID: %w", paramName, err)
	}

	return parsed, nil
}

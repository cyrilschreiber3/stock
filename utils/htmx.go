package utils

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

const DefaultNotificationTTLMS = 3000

type NotificationTrigger struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	TTLMS    int    `json:"ttlMs"`
}

func SetJSONHeader(c *gin.Context, header string, payload any) bool {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("could not marshal %s header payload: %v", header, err)
		return false
	}

	c.Header(header, string(encodedPayload))
	return true
}

func SetHXTrigger(c *gin.Context, events map[string]any) bool {
	return SetJSONHeader(c, "HX-Trigger", events)
}

func NewNotificationTrigger(severity string, message string) NotificationTrigger {
	return NotificationTrigger{
		Severity: severity,
		Message:  message,
		TTLMS:    DefaultNotificationTTLMS,
	}
}

func HXNotify(c *gin.Context, status int, severity string, message string) {
	trigger := NewNotificationTrigger(severity, message)
	SetHXTrigger(c, map[string]any{
		"notify": trigger,
	})
	c.Status(status)
}

func HXNotifyWithEvents(c *gin.Context, status int, severity string, message string, events map[string]any) {
	trigger := NewNotificationTrigger(severity, message)
	events["notify"] = trigger
	SetHXTrigger(c, events)
	c.Status(status)
}

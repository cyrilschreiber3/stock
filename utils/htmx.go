package utils

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

const DefaultNotificationTTLMS = 3000
const DefaultErrorNotificationTTLMS = 10000

type NotificationTrigger struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	TTLMS    int    `json:"ttlMs"`
}

func SetJSONHeader(c *gin.Context, header string, payload any) {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("could not marshal %s header payload: %v", header, err)
		return
	}

	c.Header(header, string(encodedPayload))
}

func SetHXTrigger(c *gin.Context, events map[string]any) {
	SetJSONHeader(c, "HX-Trigger", events)
}

func SetHXLocation(c *gin.Context, url string, target string) {
	if target != "" {
		SetJSONHeader(c, "HX-Trigger", map[string]any{
			"redirect": map[string]any{
				"url":    url,
				"target": target,
			},
		})
		return
	} else {
		c.Header("HX-Location", url)
	}
}

func NewNotificationTrigger(severity string, message string) NotificationTrigger {
	ttl := DefaultNotificationTTLMS
	if severity == "error" {
		ttl = DefaultErrorNotificationTTLMS
	}
	return NotificationTrigger{
		Severity: severity,
		Message:  message,
		TTLMS:    ttl,
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

func HXNotifyWithRedirect(c *gin.Context, status int, severity string, message string, redirectURL string) {
	trigger := NewNotificationTrigger(severity, message)
	SetHXTrigger(c, map[string]any{
		"notify": trigger,
	})
	c.Header("HX-Location", redirectURL)

}

package utils

import (
	"encoding/json"
	"log/slog"

	"github.com/gin-gonic/gin"
)

const DefaultNotificationTTLMS = 3000
const DefaultErrorNotificationTTLMS = 10000

type NotificationTrigger struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	TTLMS    int    `json:"ttlMs"`
}

type MessageTrigger struct {
	Severity string `json:"severity"`
	Content  string `json:"content"`
}

func SetJSONHeader(c *gin.Context, header string, payload any) {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("could not marshal header payload", "header", header, "error", err)
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

func HXRedirectWithNotify(c *gin.Context, status int, severity string, message string, redirectURL string) {
	trigger := NewNotificationTrigger(severity, message)
	cookie_body, err := json.Marshal(trigger)
	if err != nil {
		slog.Error("could not marshal notification trigger", "error", err)
		return
	}

	c.SetCookie("notification-toast", string(cookie_body), 15, "/", "", false, false)
	c.Header("HX-Location", redirectURL)
	c.Status(status)
}

func NewMessageTrigger(severity string, content string) MessageTrigger {
	return MessageTrigger{
		Severity: severity,
		Content:  content,
	}
}

func HXRedirectWithMessage(c *gin.Context, status int, severity string, message string, redirectURL string) {
	trigger := NewMessageTrigger(severity, message)
	cookie_body, err := json.Marshal(trigger)
	if err != nil {
		slog.Error("could not marshal message trigger", "error", err)
		return
	}

	c.SetCookie("notification-message", string(cookie_body), 15, "/", "", false, false)
	c.Header("HX-Location", redirectURL)
	c.Status(status)
}

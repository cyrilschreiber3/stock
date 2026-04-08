package components

import (
	"fmt"

	"github.com/a-h/templ"
)

type badgeConfig struct {
	label     string
	icon      templ.Component
	iconAfter bool
	color     string // neutral, primary, secondary, accent, info, success, warning, error
	style     string // solid, outline, dash, soft, ghost, link
	size      string // xs, sm, md, lg, xl
	classes   []string
}

var defaultBadgeConfig = badgeConfig{
	color:   "neutral",
	style:   "solid",
	size:    "md",
	classes: []string{},
}

func addBadgeClass(config *badgeConfig, class string) {
	config.classes = append(config.classes, fmt.Sprintf("badge-%s", class))
}

func BadgeConfig(label string) *badgeConfig {
	return &badgeConfig{
		label:   label,
		color:   defaultBadgeConfig.color,
		style:   defaultBadgeConfig.style,
		size:    defaultBadgeConfig.size,
		classes: append([]string{}, defaultBadgeConfig.classes...),
	}
}

func (c *badgeConfig) Classes(classes ...string) *badgeConfig {
	c.classes = append(c.classes, classes...)
	return c
}

func (c *badgeConfig) Icon(icon templ.Component) *badgeConfig {
	c.icon = icon
	return c
}

func (c *badgeConfig) IconAfter(icon templ.Component) *badgeConfig {
	c.icon = icon
	c.iconAfter = true
	return c
}

func (c *badgeConfig) Color(color string) *badgeConfig {
	if color == "" {
		color = "neutral"
	}
	c.color = color
	addBadgeClass(c, color)
	return c
}

func (c *badgeConfig) Size(size string) *badgeConfig {
	c.size = size
	addBadgeClass(c, size)
	return c
}

func (c *badgeConfig) Style(style string) *badgeConfig {
	c.style = style
	addBadgeClass(c, style)
	return c
}

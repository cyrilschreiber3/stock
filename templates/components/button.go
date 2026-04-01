package components

import (
	"fmt"

	"github.com/a-h/templ"
)

type buttonConfig struct {
	label      string
	icon       templ.Component
	iconAfter  bool
	color      string // neutral, primary, secondary, accent, info, success, warning, error
	style      string // solid, outline, dash, soft, ghost, link
	size       string // xs, sm, md, lg, xl
	shape      string // circle, square, wide, block
	active     bool
	disabled   bool
	btnType    string
	onEvents   map[string]string
	classes    []string
	attributes templ.Attributes
}

var defaultButtonConfig = buttonConfig{
	color:      "neutral",
	style:      "solid",
	size:       "md",
	shape:      "rounded",
	active:     false,
	disabled:   false,
	btnType:    "button",
	onEvents:   map[string]string{},
	classes:    []string{},
	attributes: templ.Attributes{},
}

func addButtonClass(config *buttonConfig, class string) {
	config.classes = append(config.classes, fmt.Sprintf("btn-%s", class))
}

func ButtonConfig(label string) *buttonConfig {
	return &buttonConfig{
		label:      label,
		color:      defaultButtonConfig.color,
		style:      defaultButtonConfig.style,
		size:       defaultButtonConfig.size,
		shape:      defaultButtonConfig.shape,
		active:     defaultButtonConfig.active,
		disabled:   defaultButtonConfig.disabled,
		btnType:    defaultButtonConfig.btnType,
		onEvents:   map[string]string{},
		classes:    append([]string{}, defaultButtonConfig.classes...),
		attributes: templ.Attributes{},
	}
}

func (c *buttonConfig) Classes(classes ...string) *buttonConfig {
	c.classes = append(c.classes, classes...)
	return c
}

func (c *buttonConfig) Attributes(attributes templ.Attributes) *buttonConfig {
	for key, value := range attributes {
		c.attributes[key] = value
	}
	return c
}

func (c *buttonConfig) Icon(icon templ.Component) *buttonConfig {
	c.icon = icon
	return c
}

func (c *buttonConfig) IconAfter(icon templ.Component) *buttonConfig {
	c.icon = icon
	c.iconAfter = true
	return c
}

func (c *buttonConfig) Color(color string) *buttonConfig {
	if color == "" {
		color = "neutral"
	}
	c.color = color
	addButtonClass(c, color)
	return c
}

func (c *buttonConfig) Size(size string) *buttonConfig {
	c.size = size
	addButtonClass(c, size)
	return c
}

func (c *buttonConfig) Style(style string) *buttonConfig {
	c.style = style
	addButtonClass(c, style)
	return c
}

func (c *buttonConfig) Shape(shape string) *buttonConfig {
	if c.shape != "" {
		// remove old shape class
		oldShapeClass := fmt.Sprintf("btn-%s", c.shape)
		for i, class := range c.classes {
			if class == oldShapeClass {
				c.classes = append(c.classes[:i], c.classes[i+1:]...)
				break
			}
		}
	}
	c.shape = shape
	addButtonClass(c, shape)
	return c
}

func (c *buttonConfig) Active(active bool) *buttonConfig {
	c.active = active
	if active {
		addButtonClass(c, "active")
	}
	return c
}

func (c *buttonConfig) Disabled(disabled bool) *buttonConfig {
	c.disabled = disabled
	if disabled {
		addButtonClass(c, "disabled")
	}
	return c
}

func (c *buttonConfig) Type(btnType string) *buttonConfig {
	c.btnType = btnType
	c.attributes["type"] = btnType
	return c
}

func (c *buttonConfig) On(event string, handler string) *buttonConfig {
	c.onEvents[event] = handler
	c.attributes[fmt.Sprintf("on%s", event)] = handler
	return c
}

func (c *buttonConfig) OnClick(handler string) *buttonConfig {
	return c.On("click", handler)
}

func (c *buttonConfig) At(event string, handler string) *buttonConfig {
	c.attributes[fmt.Sprintf("@%s", event)] = handler
	return c
}

func (c *buttonConfig) AtClick(handler string) *buttonConfig {
	return c.At("click", handler)
}

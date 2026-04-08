package components

import (
	"fmt"

	"github.com/a-h/templ"
)

type formFieldConfig struct {
	name            string
	label           string
	inputType       string
	value           string
	options         []string
	placeholder     string
	color           string // neutral, primary, secondary, accent, info, success, warning, error
	size            string // xs, sm, md, lg, xl
	prefix          string
	suffix          string
	required        bool
	checked         bool
	validationHint  string
	inputClasses    []string
	inputAttributes templ.Attributes
	customInput     templ.Component
}

func FormFieldConfig(name string) *formFieldConfig {
	return &formFieldConfig{
		name:            name,
		inputType:       "text",
		color:           "neutral",
		size:            "md",
		prefix:          "",
		suffix:          "",
		required:        false,
		checked:         false,
		inputClasses:    []string{"input"},
		inputAttributes: templ.Attributes{},
		customInput:     nil,
	}
}

func (c *formFieldConfig) getInputComponent() templ.Component {
	var classes []string
	if c.prefix == "" && c.suffix == "" {
		classes = c.inputClasses
	}

	switch c.inputType {
	case "select":
		return selectFormInput(c.inputClasses, c.inputAttributes, c.options, c.placeholder)
	case "radio":
		return radioFormInput(c.name, c.inputClasses, c.options)
	case "textarea":
		return textareaFormInput(c.inputClasses, c.inputAttributes)
	case "custom":
		return c.customInput
	default:
		return formInput(classes, c.inputAttributes)
	}

}

func (c *formFieldConfig) Label(label string) *formFieldConfig {
	c.label = label
	return c
}

func (c *formFieldConfig) Type(inputType string) *formFieldConfig {
	c.inputType = inputType
	c.inputAttributes["type"] = inputType
	switch inputType {
	case "text", "email", "password", "number":
		c.Classes("input")
	case "checkbox":
		c.Classes("toggle")
	case "custom":
		// No default classes for custom input
	default:
		c.Classes(inputType)
	}
	return c
}

func (c *formFieldConfig) Color(color string) *formFieldConfig {
	return c.TypeClass(color)
}

func (c *formFieldConfig) Size(size string) *formFieldConfig {
	return c.TypeClass(size)
}

func (c *formFieldConfig) Value(value string) *formFieldConfig {
	c.value = value
	c.inputAttributes["value"] = value
	return c
}

func (c *formFieldConfig) Placeholder(placeholder string) *formFieldConfig {
	c.placeholder = placeholder
	c.inputAttributes["placeholder"] = placeholder
	return c
}

func (c *formFieldConfig) Prefix(prefix string) *formFieldConfig {
	c.prefix = prefix
	return c
}

func (c *formFieldConfig) Suffix(suffix string) *formFieldConfig {
	c.suffix = suffix
	return c
}

func (c *formFieldConfig) Required() *formFieldConfig {
	c.required = true
	c.inputAttributes["required"] = true
	return c
}

func (c *formFieldConfig) Checked() *formFieldConfig {
	c.checked = true
	c.inputAttributes["checked"] = true
	return c
}

func (c *formFieldConfig) Options(options []string) *formFieldConfig {
	c.options = options
	return c
}

func (c *formFieldConfig) ValidationAttributes(attrs map[string]string) *formFieldConfig {
	for key, value := range attrs {
		c.inputAttributes[key] = value
	}
	return c
}

func (c *formFieldConfig) ValidationPreset(preset string) *formFieldConfig {
	switch preset {
	case "money":
		c.ValidationAttributes(map[string]string{
			"step": "0.01",
			"min":  "0",
		})
	case "integer":
		c.ValidationAttributes(map[string]string{
			"step": "1",
		})
	case "positive":
		c.ValidationAttributes(map[string]string{
			"min": "0",
		})
	}
	return c
}

func (c *formFieldConfig) ValidationHint(hint string) *formFieldConfig {
	c.validationHint = hint
	return c
}

func (c *formFieldConfig) TypeClass(class string) *formFieldConfig {
	c.color = class
	switch c.inputType {
	case "text", "email", "password", "number":
		c.Classes(fmt.Sprintf("input-%s", class))
	case "checkbox":
		c.Classes(fmt.Sprintf("toggle-%s", class))
	default:
		c.Classes(fmt.Sprintf("%s-%s", c.inputType, class))
	}
	return c
}

func (c *formFieldConfig) Classes(classes ...string) *formFieldConfig {
	c.inputClasses = append(c.inputClasses, classes...)
	return c
}

func (c *formFieldConfig) Attributes(attributes templ.Attributes) *formFieldConfig {
	for key, value := range attributes {
		c.inputAttributes[key] = value
	}
	return c
}

func (c *formFieldConfig) CustomInput(component templ.Component) *formFieldConfig {
	c.customInput = component
	c.inputType = "custom"
	return c
}

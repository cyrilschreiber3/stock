package components

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/a-h/templ"
)

type formFieldConfig struct {
	name            string
	label           string
	inputType       string
	value           string
	valueList       []string
	options         map[string]string
	optionsEndpoint string
	dependsOn       string
	placeholder     string
	color           string // neutral, primary, secondary, accent, info, success, warning, error
	size            string // xs, sm, md, lg, xl
	prefix          string
	suffix          string
	required        bool
	checked         bool
	asList          bool
	validationHint  string
	inputClasses    []string
	inputAttributes templ.Attributes
	customInput     templ.Component
}

func FormFieldConfig(name string) *formFieldConfig {
	return &formFieldConfig{
		name:         name,
		inputType:    "text",
		color:        "neutral",
		size:         "md",
		prefix:       "",
		suffix:       "",
		required:     false,
		checked:      false,
		inputClasses: []string{"validator"},
		inputAttributes: templ.Attributes{
			"id":   name,
			"name": name,
		},
		customInput: nil,
	}
}

func (c *formFieldConfig) getInputComponent() templ.Component {
	var classes []string
	var component templ.Component
	if c.prefix == "" && c.suffix == "" {
		classes = c.inputClasses
	}

	if c.asList {
		classes = append(classes, "join-item")
	}

	switch c.inputType {
	case "select":
		component = selectFormInput(classes, c.inputAttributes, c.options, c.placeholder)
	case "radio":
		component = radioFormInput(c.name, classes, c.options)
	case "textarea":
		component = textareaFormInput(classes, c.inputAttributes)
	case "custom":
		component = c.customInput
	default:
		component = formInput(classes, c.inputAttributes)
	}

	if c.asList {
		c = c.Attributes(templ.Attributes{"x-model": fmt.Sprintf("%s[index]", c.name)})
		return listFormInput(c, component)
	} else {
		return component
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
	switch c.inputType {
	case "checkbox":
		c.inputAttributes["checked"] = value == "true"
	default:
		c.inputAttributes["value"] = value
	}
	return c
}

func (c *formFieldConfig) ValueList(values []string) *formFieldConfig {
	c.valueList = values
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

func (c *formFieldConfig) Options(options map[string]string) *formFieldConfig {
	c.options = options
	return c
}

func (c *formFieldConfig) OptionsFromEndpoint(endpoint string) *formFieldConfig {
	query := url.Values{}
	if c.placeholder != "" {
		query.Set("placeholder", c.placeholder)
	}
	if c.value != "" {
		query.Set("value", c.value)
	}

	if len(query) > 0 {
		c.optionsEndpoint = fmt.Sprintf("%s?%s", endpoint, query.Encode())
	} else {
		c.optionsEndpoint = endpoint
	}

	c.inputAttributes["hx-get"] = c.optionsEndpoint
	c.inputAttributes["hx-trigger"] = "load"
	c.inputAttributes["hx-target"] = "this"
	c.inputAttributes["hx-swap"] = "innerHTML"
	// c.inputAttributes["hx-on:htmx:afterSwap"] = `this.dispatchEvent(new CustomEvent("change", { bubbles: true }))`
	return c
}

func (c *formFieldConfig) DependsOn(fieldName string) *formFieldConfig {
	c.dependsOn = fieldName
	c.inputAttributes["hx-trigger"] = fmt.Sprintf("load, change from:#%s, htmx:afterSwap from:#%s", fieldName, fieldName)
	c.inputAttributes["hx-target"] = "this"
	c.inputAttributes["hx-swap"] = "innerHTML"

	// If endpoint contains a path token like :id, replace it at request time.
	// Example: /categories/:id/subcategories/options
	if strings.Contains(c.optionsEndpoint, ":id") {
		c.inputAttributes["hx-on:htmx:config-request"] = fmt.Sprintf(
			`event.detail.path = %q.replace(":id", encodeURIComponent(document.getElementById(%q)?.value || ""))`,
			c.optionsEndpoint,
			fieldName,
		)
	} else {
		// Otherwise pass parent value as query param (e.g. ?category_id=...)
		c.inputAttributes["hx-vals"] = fmt.Sprintf(
			`js:{ %q: document.getElementById(%q)?.value || "" }`,
			fieldName,
			fieldName,
		)
	}
	return c
}

func (c *formFieldConfig) AsList() *formFieldConfig {
	c.asList = true
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

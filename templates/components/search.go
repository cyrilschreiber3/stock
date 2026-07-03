package components

import (
	"net/url"

	"github.com/a-h/templ"
)

type searchConfig struct {
	name             string
	label            string // not implemented
	placeholder      string
	icon             templ.Component
	value            string
	endpoint         string
	filters          map[string]string
	filterDefault    string
	filtersEndpoint  string
	floatingResults  bool
	inputClasses     []string
	buttonClasses    []string
	inputAttributes  templ.Attributes
	buttonAttributes templ.Attributes
}

func SearchConfig(name string) *searchConfig {
	return &searchConfig{
		name:             name,
		label:            "",
		placeholder:      "Search",
		icon:             nil,
		value:            "",
		endpoint:         "",
		filters:          map[string]string{},
		filterDefault:    "Filter",
		filtersEndpoint:  "",
		floatingResults:  false,
		inputClasses:     []string{"input", "join-item", "w-full"},
		buttonClasses:    []string{"select", "join-item", "w-fit", "shrink-0"},
		inputAttributes:  templ.Attributes{},
		buttonAttributes: templ.Attributes{},
	}
}

func (c *searchConfig) Label(label string) *searchConfig {
	c.label = label
	return c
}

func (c *searchConfig) Placeholder(placeholder string) *searchConfig {
	c.placeholder = placeholder
	c.inputAttributes["placeholder"] = placeholder
	return c
}

func (c *searchConfig) Icon(icon templ.Component) *searchConfig {
	c.icon = icon
	return c
}

func (c *searchConfig) Value(value string) *searchConfig {
	c.value = value
	c.inputAttributes["value"] = value
	return c
}

func (c *searchConfig) Endpoint(endpoint string) *searchConfig {
	c.endpoint = endpoint
	return c
}

func (c *searchConfig) Filters(filters map[string]string) *searchConfig {
	c.filters = filters
	return c
}

func (c *searchConfig) FilterDefault(filterDefault string) *searchConfig {
	c.filterDefault = filterDefault
	return c
}

func (c *searchConfig) FiltersFromEndpoint(endpoint string) *searchConfig {
	c.filtersEndpoint = endpoint

	query := url.Values{}
	if c.filterDefault != "" {
		query.Set("placeholder", c.filterDefault)
	}

	if len(query) > 0 {
		c.filtersEndpoint += "?" + query.Encode()
	}

	c.buttonAttributes["hx-get"] = c.filtersEndpoint
	c.buttonAttributes["hx-trigger"] = "load, refreshFilters"
	c.buttonAttributes["hx-target"] = "this"
	c.buttonAttributes["hx-swap"] = "innerHTML"
	return c
}

func (c *searchConfig) FloatingResults(floatingResults bool) *searchConfig {
	c.floatingResults = floatingResults
	return c
}

func (c *searchConfig) InputClasses(classes ...string) *searchConfig {
	c.inputClasses = append(c.inputClasses, classes...)
	return c
}

func (c *searchConfig) ButtonClasses(classes ...string) *searchConfig {
	c.buttonClasses = append(c.buttonClasses, classes...)
	return c
}

func (c *searchConfig) InputAttributes(attributes templ.Attributes) *searchConfig {
	for key, value := range attributes {
		c.inputAttributes[key] = value
	}
	return c
}

func (c *searchConfig) ButtonAttributes(attributes templ.Attributes) *searchConfig {
	for key, value := range attributes {
		c.buttonAttributes[key] = value
	}
	return c
}

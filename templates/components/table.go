package components

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/cyrilschreiber3/stock/routes"
	"github.com/gin-gonic/gin"
)

type DataTableConfig struct {
	Label              string
	Name               string
	SearchPlaceholder  string
	SearchValue        string
	SearchDisabled     bool
	Columns            []DataTableColumn
	ItemCount          int
	ItemLabel          string
	ItemLabelSingular  string
	ItemLabelKey       string
	SortKey            string
	SortDirection      string
	CurrentPage        int
	TotalPages         int
	PageSize           int
	PageSizeOptions    []int
	URL                routes.StaticRoute
	URLString          string
	AdditionalTriggers []string
}

type DataTableColumn struct {
	Label    string
	Key      string
	Sortable bool
}

type pageNumberDisplay struct {
	ShowPrefix  bool
	ShowSuffix  bool
	DisablePrev bool
	DisableNext bool
	Pages       []int
}

func NewDataTableColumn(label, key string) DataTableColumn {
	return DataTableColumn{
		Label:    label,
		Key:      key,
		Sortable: true,
	}
}

func NewUnsortedDataTableColumn(label, key string) DataTableColumn {
	return DataTableColumn{
		Label:    label,
		Key:      key,
		Sortable: false,
	}
}

func NewActionDataTableColumn() DataTableColumn {
	return DataTableColumn{
		Label:    "Actions",
		Key:      "actions",
		Sortable: false,
	}
}

func NewDataTableConfig(name string) *DataTableConfig {
	return &DataTableConfig{
		Name:               name,
		Columns:            []DataTableColumn{},
		PageSize:           10,
		PageSizeOptions:    []int{10, 25, 50, 100},
		CurrentPage:        1,
		TotalPages:         1,
		SearchDisabled:     false,
		AdditionalTriggers: []string{},
	}
}

func (c *DataTableConfig) SetLabel(label string) *DataTableConfig {
	c.Label = label
	return c
}

func (c *DataTableConfig) SetSearchPlaceholder(placeholder string) *DataTableConfig {
	c.SearchPlaceholder = placeholder
	return c
}

func (c *DataTableConfig) SetSearchValue(value string) *DataTableConfig {
	c.SearchValue = value
	return c
}

func (c *DataTableConfig) DisableSearch() *DataTableConfig {
	c.SearchDisabled = true
	return c
}

func (c *DataTableConfig) AddColumns(column ...DataTableColumn) *DataTableConfig {
	c.Columns = append(c.Columns, column...)
	return c
}

func (c *DataTableConfig) SetItemCount(count int) *DataTableConfig {
	c.ItemCount = count
	return c
}

func (c *DataTableConfig) SetItemLabel(label string, singular string) *DataTableConfig {
	c.ItemLabel = label
	c.ItemLabelSingular = singular
	return c
}

func (c *DataTableConfig) SetItemLabelKey(key string) *DataTableConfig {
	c.ItemLabelKey = key
	return c
}

func (c *DataTableConfig) SetSortKey(key string) *DataTableConfig {
	c.SortKey = key
	return c
}

func (c *DataTableConfig) SetSortDirection(direction string) *DataTableConfig {
	c.SortDirection = direction
	return c
}

func (c *DataTableConfig) SetCurrentPage(page int) *DataTableConfig {
	c.CurrentPage = page
	return c
}

func (c *DataTableConfig) SetPageSize(size int) *DataTableConfig {
	c.PageSize = size
	return c
}

func (c *DataTableConfig) SetPageSizeOptions(options []int) *DataTableConfig {
	c.PageSizeOptions = options
	return c
}

func (c *DataTableConfig) SetURL(url routes.StaticRoute) *DataTableConfig {
	c.URL = url
	c.URLString = url.URL()
	return c
}

func (c *DataTableConfig) SetURLString(url string) *DataTableConfig {
	c.URLString = url
	return c
}

func (c *DataTableConfig) AddAdditionalTriggers(triggers ...string) *DataTableConfig {
	c.AdditionalTriggers = append(c.AdditionalTriggers, triggers...)
	return c
}

func (c *DataTableConfig) buildConfig() *DataTableConfig {
	c.TotalPages = c.ItemCount / c.PageSize
	if c.ItemCount%c.PageSize != 0 {
		c.TotalPages++
	}

	if c.CurrentPage > c.TotalPages {
		c.CurrentPage = c.TotalPages
	}

	if c.CurrentPage < 1 {
		c.CurrentPage = 1
	}
	return c
}

func (c *DataTableConfig) GetConfigFromURL(ctx *gin.Context) *DataTableConfig {
	if search := ctx.Query(c.getSearchName()); search != "" {
		c.SearchValue = search
	}

	if page := ctx.Query(c.getPageSelectName()); page != "" {
		_, err := fmt.Sscanf(page, "%d", &c.CurrentPage)
		if err != nil {
			c.CurrentPage = 1
		}
	}

	if pageSize := ctx.Query(c.getPageSizeName()); pageSize != "" {
		_, err := fmt.Sscanf(pageSize, "%d", &c.PageSize)
		if err != nil {
			c.PageSize = 10
		}
	}

	if sortKey := ctx.Query(c.getSortName()); sortKey != "" {
		c.SortKey = sortKey
	}

	if sortDirection := ctx.Query(c.getSortDirectionName()); sortDirection != "" {
		c.SortDirection = sortDirection
	}

	return c
}

func (c *DataTableConfig) SetConfigToURL(ctx *gin.Context) *DataTableConfig {
	query := ctx.Request.URL.Query()
	refererURL, err := url.Parse(ctx.GetHeader("Referer"))
	if err != nil {
		slog.Error("Error parsing referer URL", "error", err)
		refererURL = &url.URL{Path: ctx.Request.URL.Path, RawQuery: query.Encode()}
	}

	for key, value := range refererURL.Query() {
		query.Set(key, value[0])
	}

	if c.SearchValue != "" {
		query.Set(c.getSearchName(), c.SearchValue)
	} else {
		query.Del(c.getSearchName())
	}

	query.Set(c.getPageSelectName(), fmt.Sprintf("%d", c.CurrentPage))
	query.Set(c.getPageSizeName(), fmt.Sprintf("%d", c.PageSize))

	if c.SortKey != "" {
		query.Set(c.getSortName(), c.SortKey)
	} else {
		query.Del(c.getSortName())
	}

	if c.SortDirection != "" {
		query.Set(c.getSortDirectionName(), c.SortDirection)
	} else {
		query.Del(c.getSortDirectionName())
	}

	encodedQuery := query.Encode()
	pushURL := refererURL.Path
	if encodedQuery != "" {
		pushURL += "?" + encodedQuery
	}
	ctx.Header("HX-Push-Url", pushURL)
	return c
}

func (c *DataTableConfig) GetSortDirectionForQuery() string {
	switch c.SortDirection {
	case "asc", "ASC":
		return "ASC"
	case "desc", "DESC":
		return "DESC"
	default:
		return "ASC"
	}
}

func (c *DataTableConfig) GetPageOffset() int {
	return (c.CurrentPage - 1) * c.PageSize
}

func (c *DataTableConfig) GetPageLimit(elementCount int) int {
	if c.GetPageOffset()+c.PageSize > elementCount {
		return elementCount - c.GetPageOffset()
	}
	return c.PageSize
}

func (c *DataTableConfig) getFormTriggerString() string {
	trigger := fmt.Sprintf("submit, input from:[form=%s-form] delay:200ms", c.Name)
	for _, additionalTrigger := range c.AdditionalTriggers {
		trigger = fmt.Sprintf("%s, %s from:body", trigger, additionalTrigger)
	}
	return trigger
}

func (c *DataTableConfig) sortState(col string) string {
	if c.SortKey == col {
		switch c.SortDirection {
		case "asc":
			return "asc"
		case "desc":
			return "desc"
		}
	}
	return ""
}

func (c *DataTableConfig) nextSortDirection(col string) string {
	if c.SortKey == col {
		switch c.SortDirection {
		case "asc":
			return "desc"
		case "desc":
			return "asc"
		}
	}
	return "asc"
}

func (c *DataTableConfig) shouldShowArrow(col string) bool {
	return c.SortKey == col
}

func (c *DataTableConfig) countLabel() string {
	label := c.ItemLabel
	if c.ItemCount == 1 {
		label = c.ItemLabelSingular
	}
	return fmt.Sprintf("%d %s", c.ItemCount, label)
}

func (c *DataTableConfig) pageNumbers() *pageNumberDisplay {
	pc := &pageNumberDisplay{
		Pages: []int{0},
	}

	if c.CurrentPage == 1 {
		pc.DisablePrev = true
	}

	if c.CurrentPage == c.TotalPages {
		pc.DisableNext = true
	}

	if c.TotalPages <= 1 {
		pc.Pages = []int{1}
	} else if c.TotalPages <= 5 {
		pages := make([]int, c.TotalPages)
		for i := 1; i <= c.TotalPages; i++ {
			pages[i-1] = i
		}
		pc.Pages = pages
	} else if c.TotalPages > 5 {
		if c.CurrentPage <= 3 {
			pc.ShowSuffix = true
			pc.Pages = []int{1, 2, 3}
		} else if c.CurrentPage >= c.TotalPages-2 {
			pc.ShowPrefix = true
			pc.Pages = []int{c.TotalPages - 2, c.TotalPages - 1, c.TotalPages}
		} else {
			pc.ShowPrefix = true
			pc.ShowSuffix = true
			pc.Pages = []int{c.CurrentPage}
		}
	}
	return pc
}

func (c *DataTableConfig) getFormName() string {
	return c.Name + "-form"
}

func (c *DataTableConfig) getSearchName() string {
	return c.Name + "-search"
}

func (c *DataTableConfig) getSortName() string {
	return c.Name + "-sort"
}

func (c *DataTableConfig) getSortDirectionName() string {
	return c.Name + "-sort-direction"
}

func (c *DataTableConfig) getPageSizeName() string {
	return c.Name + "-page-size"
}

func (c *DataTableConfig) getPageSelectName() string {
	return c.Name + "-page-select"
}

func (c *DataTableConfig) getElementName(name string) string {
	return c.Name + "-" + name
}

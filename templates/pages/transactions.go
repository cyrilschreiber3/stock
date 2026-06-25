package pages

import (
	"github.com/a-h/templ"
	"github.com/cyrilschreiber3/stock/templates/components"
)

func transactionTypeBadge(transactionType string) templ.Component {
	switch transactionType {
	case "buy":
		return components.Badge(components.BadgeConfig("Purchase").Color("primary"))
	case "sell":
		return components.Badge(components.BadgeConfig("Sale").Color("secondary"))
	case "adjustment":
		return components.Badge(components.BadgeConfig("Adjustment").Color("accent"))
	case "correction":
		return components.Badge(components.BadgeConfig("Correction").Color("warning"))
	default:
		return components.Badge(components.BadgeConfig(transactionType).Color("neutral"))
	}
}

func transactionStateBadge(status string) templ.Component {
	switch status {
	case "draft":
		return components.Badge(components.BadgeConfig("Draft").Color("neutral"))
	case "pendingRefund":
		return components.Badge(components.BadgeConfig("Pending Refund").Color("warning"))
	case "completed":
		return components.Badge(components.BadgeConfig("Completed").Color("success"))
	default:
		return components.Badge(components.BadgeConfig(status).Color("neutral"))
	}
}

package locales

import (
	"context"

	"github.com/cyrilschreiber3/stock/locales/lang"
	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
)

var SupportedLanguages = []string{"en", "fr"}
var ShowMissing = false

func init() {
	if err := ctxi18n.LoadWithDefault(lang.LocalesFS, "en"); err != nil {
		panic(err)
	}
}

func New(ctx context.Context) *Locale {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	return &Locale{
		Context: ctx,
	}
}

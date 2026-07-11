package static

import "embed"

//go:embed *
var StaticAssets embed.FS

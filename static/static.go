package static

import "embed"

//go:embed *
var StaticAssets embed.FS

//go:embed favicon.ico
var Favicon []byte

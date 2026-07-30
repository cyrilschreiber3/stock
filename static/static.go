package static

import (
	"embed"
	"runtime/debug"
)

//go:embed *
var StaticAssets embed.FS

//go:embed favicon.ico
var Favicon []byte

var Version = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			return s.Value[:7]
		}
	}
	return "dev"
}()

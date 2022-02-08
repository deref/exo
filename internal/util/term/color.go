package term

import (
	"crypto/md5"
	"os"

	"github.com/lucasb-eyer/go-colorful"
)

const ResetCode = "\u001b[0m"

func colorIsBlack(c colorful.Color) bool {
	return c.R == 0 && c.G == 0 && c.B == 0
}

// See <https://bixense.com/clicolors/>.
func IsColorEnabled() bool {
	force, _ := os.LookupEnv("CLICOLOR_FORCE")
	if force != "" {
		return force != "0"
	}
	prefer, _ := os.LookupEnv("CLICOLOR")
	return IsInteractive() && prefer != "0"
}

// Maps a key to a stable, arbitrary color.
type ColorCache struct {
	palette []colorful.Color
	colors  map[string]colorful.Color
}

func NewColorCache() *ColorCache {
	pal, err := colorful.HappyPalette(256)
	if err != nil {
		// An error should only be possible if the number of colours requested is
		// too high. Since this is a fixed constant this panic should be impossible.
		panic(err)
	}
	return &ColorCache{
		palette: pal,
		colors:  make(map[string]colorful.Color),
	}
}

func (cache *ColorCache) Color(key string) colorful.Color {
	color := cache.colors[key]
	if colorIsBlack(color) {
		b := md5.Sum([]byte(key))[0]
		color = cache.palette[b]
		cache.colors[key] = color
	}
	return color
}

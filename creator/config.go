package creator

import "image/color"

type Config struct {
	ImageSize  int
	ImageRows  int
	ImageCols  int
	ImageApart int
	BgColor    color.Color
	Reverse    bool
	CacheDir   string
	MaxActive  int
}

var DefaultConfig = &Config{
	60,
	0,
	10,
	10,
	color.White,
	true,
	"./github-avatar-cache",
	10,
}

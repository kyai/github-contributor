package creator

type Config struct {
	ImageSize  int
	ImageRows  int
	ImageCols  int
	ImageApart int
	CacheDir   string
	MaxActive  int
}

var DefaultConfig = &Config{
	60,
	0,
	10,
	10,
	"./github-avatar-cache",
	10,
}

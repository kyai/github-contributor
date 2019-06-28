package creator

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"code.google.com/p/graphics-go/graphics"
)

type Creator struct {
	repo string
	conf *Config
}

func New(repo string) *Creator {
	return &Creator{
		repo: repo,
		conf: DefaultConfig,
	}
}

func (c *Creator) Set(conf *Config) *Creator {
	c.conf = conf
	return c
}

func (c *Creator) Create() (path string, err error) {
	images, err := c.downInRestricted()
	if err != nil {
		return
	}
	return c.write(c.merge(images...))
}

func (c *Creator) downInRestricted() (images []image.Image, err error) {
	active := make(chan int, c.conf.MaxActive)
	defer close(active)

	users, err := getContributorsByRepo(c.repo)
	if err != nil {
		return
	}

	if c.conf.Reverse {
		for i, j := 0, len(users)-1; i < j; i, j = i+1, j-1 {
			users[i], users[j] = users[j], users[i]
		}
	}

	images = make([]image.Image, len(users))

	var wg sync.WaitGroup

	for i, user := range users {
		wg.Add(1)
		active <- 1

		go func(i int, user *Contributor) {
			defer func() {
				<-active
				wg.Done()
			}()

			var errImg error
			images[i], errImg = httpGetImage(user.Author.Avatar)
			if errImg != nil {
				fmt.Println(fmt.Sprintf("Get <%s> has error: %s", user.Author.Avatar, errImg.Error()))
				return
			}
		}(i, user)
	}

	wg.Wait()

	return
}

func (c *Creator) merge(images ...image.Image) image.Image {

	cols := c.conf.ImageCols
	rows := int(math.Ceil(float64(len(images)) / float64(cols)))

	width := cols*c.conf.ImageSize + (cols-1)*c.conf.ImageApart
	height := rows*c.conf.ImageSize + (rows-1)*c.conf.ImageApart

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// background color
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{c.conf.BgColor}, image.ZP, draw.Src)

	for i, img := range images {
		img, _ = c.resize(img)

		var (
			row  = i / cols
			col  = i % cols
			left = col * (c.conf.ImageSize + c.conf.ImageApart)
			top  = row * (c.conf.ImageSize + c.conf.ImageApart)
		)

		offset := image.Pt(left, top)
		draw.Draw(canvas, img.Bounds().Add(offset), img, image.ZP, draw.Over)
	}

	return canvas
}

func (c *Creator) resize(img image.Image) (image.Image, error) {
	size := c.conf.ImageSize
	if img.Bounds().Size().X == size {
		return img, nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	err := graphics.Scale(dst, img)
	return dst, err
}

func (c *Creator) write(img image.Image) (path string, err error) {
	if err = os.MkdirAll(c.conf.CacheDir, os.ModePerm); err != nil {
		return
	}

	fileName := strings.Replace(c.repo, "/", "-", 1) + ".jpeg"
	file, err := os.Create(filepath.Join(c.conf.CacheDir, fileName))
	if err != nil {
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, img, &jpeg.Options{c.conf.Quality})
	path = file.Name()
	return
}

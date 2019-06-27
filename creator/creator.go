package creator

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math"
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

func (c *Creator) Create() {
	images, err := c.downInRestricted()
	if err != nil {
		panic(err)
	}
	c.merge(images...)
}

func (c *Creator) downInRestricted() (images []image.Image, err error) {
	active := make(chan int, c.conf.MaxActive)

	users, err := getContributorsByRepo(c.repo)
	if err != nil {
		return
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

func (c *Creator) merge(images ...image.Image) {

	cols := c.conf.ImageCols
	rows := int(math.Ceil(float64(len(images)) / float64(cols)))

	width := cols*c.conf.ImageSize + (cols-1)*c.conf.ImageApart
	height := rows*c.conf.ImageSize + (rows-1)*c.conf.ImageApart

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

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

	imgw, err := ioutil.TempFile("./", "image-")
	if err != nil {
		panic(err)
	}
	defer imgw.Close()
	err = jpeg.Encode(imgw, canvas, &jpeg.Options{jpeg.DefaultQuality})

	fmt.Println(imgw.Name(), err)
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

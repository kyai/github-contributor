package creator

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
)

func httpGet(url string) (res []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func httpGetWithHeader(url string, headers map[string]string) (res []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func httpGetImage(url string) (img image.Image, err error) {
	res, err := httpGet(url)
	if err != nil {
		return
	}

	if img, err = jpeg.Decode(bytes.NewReader(res)); err == nil {
		return
	}

	if img, err = png.Decode(bytes.NewReader(res)); err == nil {
		return
	}

	return
}

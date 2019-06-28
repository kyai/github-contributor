package creator

import (
	"fmt"
	"testing"
)

func TestCreator(t *testing.T) {
	path, err := New("golang/go").Create()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
}

func TestDownload(t *testing.T) {
	imgs, err := New("golang/go").downInRestricted()
	if err != nil {
		panic(err)
	}
	for k, v := range imgs {
		fmt.Println(k, v != nil)
	}
}

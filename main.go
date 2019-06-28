package main

import (
	"fmt"
	"os"

	"github.com/kyai/github-contributors/creator"
)

func main() {
	var repo string
	if len(os.Args) > 1 {
		repo = os.Args[1]
	} else {
		fmt.Println("Please specify the repository.")
		os.Exit(0)
	}
	path, err := creator.New(repo).Create()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("The image is saved in", path)
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"log"
	"github.com/docker/docker/client"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Printf("error == %s","asdfasdf")
	if len(os.Args) != 3 {
		fmt.Println("Usage: injecto <image> <container>")
		os.Exit(1)
	}

	image := processImage(os.Args[1])
	container := os.Args[2]

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("error == %s",err )
		panic(err)
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Printf("error == %s",err )
		panic(err)
	}
	if err := save(cli, dir, image); err != nil {
		log.Printf("copying [%s]: %s\n", dir, image)
		log.Printf("error == %s",err )
		panic(err)
	}
	if err := copy(cli, dir, container); err != nil {
		log.Printf("error == %s",err )
		panic(err)
	}
}

func processImage(s string) string {
	parts := strings.Split(s, ":")
	if len(parts) == 1 {
		return strings.Join([]string{s, "latest"}, ":")
	}
	return s
}

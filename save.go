package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"log"
	"github.com/docker/docker/client"
	"github.com/hightouchio/injecto/tar"
)

var (
	blacklist = []string{
		"dev",
		"etc/hostname",
		"etc/hosts",
		"etc/motd",
		"etc/modules-load.d",
		"etc/mtab",
		"etc/resolv.conf",
		"media",
		"mnt",
		"sys",
		"tmp",
	}
)

type manifestEntry struct {
	Layers []string
}

func save(cli *client.Client, dir, image string) error {
	log.Printf("saving %s\n", image)

	reader, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		log.Printf("error == %s",err )
		return err
	}

	saveDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Printf("error == %s",err )
		return err
	}

	if err := tar.Extract(reader, saveDir, blacklist); err != nil {
		log.Printf("error == %s",err )
		return err
	}

	manifestBytes, err := ioutil.ReadFile(path.Join(saveDir, "manifest.json"))
	if err != nil {
		log.Printf("error == %s",err )
		return err
	}

	var manifest []manifestEntry
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		log.Printf("error == %s",err )
		return err
	}

	for i, layer := range manifest[0].Layers {
		filename := path.Join(saveDir, layer)

		layerFile, err := os.Open(filename)
		if err != nil {
			log.Printf("error == %s",err )
			return err
		}

		log.Printf("extracting layer [%d/%d] %s\n", i+1, len(manifest[0].Layers),manifest[0].Layers)
		if err := tar.Extract(layerFile, dir, blacklist); err != nil {
			log.Printf("under extracting %s\n", err)
			log.Printf("error == %s",err )
			return err
		}
	}

	return nil
}

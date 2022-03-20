package main

import (
	"fmt"

	"github.com/heroku/docker-registry-client/registry"
)

func main() {
	ruleList := map[string]int{
		"ubuntu":  2,
		"alpine":  3,
		"grafana": 1,
		"redis":   2,
	}

	url := "https://registry.duy.io/"
	username := "duyuser"
	password := "password123456"
	hub, _ := registry.NewInsecure(url, username, password)
	repositories, _ := hub.Repositories()

	// Fetch list of all repo
	// Go 1 by 1, match with the hashList (whitelist)
	// If hit with whitelist,
	// Get all tags
	// ex: keep 2, then remove all, except the latest 2
	for _, repo := range repositories {
		fmt.Printf("Scanning repo: %s\n", repo)
		if imgCount, ok := ruleList[repo]; ok {
			tags, _ := hub.Tags(repo)
			fmt.Println(tags)

			for i, tag := range tags {
				// Remove all image, except the latest "imgCount" tags!
				if i >= imgCount {
					break
				}

				digest, _ := hub.ManifestDigest(repo, tag)
				fmt.Printf("Delete image %s:%s\n", repo, tag)
				fmt.Printf("Img Digest %s\n", digest)
				err := hub.DeleteManifest(repo, digest)

				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

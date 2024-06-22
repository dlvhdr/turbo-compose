package cmd

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func GetLocalImages() (map[string]image.Summary, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()
	images, err := apiClient.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return nil, err
	}

	imgDict := make(map[string]image.Summary)
	for _, img := range images {
		firstTag := ""
		tags := img.RepoTags
		if len(tags) > 0 {
			firstTag = tags[0]
		}
		nameParts := strings.Split(firstTag, ":")
		name := "none"
		if len(nameParts) > 1 {
			name = strings.Join(nameParts[:len(nameParts)-1], ":")
		}
		imgDict[name] = img
	}
	return imgDict, nil
}

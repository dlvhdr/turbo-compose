package docker

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"
)

type ServiceOption struct {
	Name       string
	LocalImage *image.Summary
	Image      string
}

type DockerCompose struct {
	Version  string
	Services map[string]ServiceDefinition
}

type ServiceDefinition struct {
	Name  string
	Image string
}

func GetLocalImages(repo string) (map[string]image.Summary, error) {
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
		name = strings.TrimPrefix(name, repo+"/")

		imgDict[name] = img
	}
	return imgDict, nil
}

func ListServicesFromComposeFile(repo string) []ServiceDefinition {
	yamlFile, err := os.ReadFile("/Users/dlvhdr/code/komodor/mono/docker-compose.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	var dockerCompose DockerCompose
	err = yaml.Unmarshal(yamlFile, &dockerCompose)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	res := make([]ServiceDefinition, 0)
	for serviceName, service := range dockerCompose.Services {
		if strings.HasPrefix(service.Image, repo) {
			res = append(res, ServiceDefinition{Name: serviceName, Image: service.Image})
		}
	}

	return res
}

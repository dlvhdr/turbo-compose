/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/dlvhdr/turbo-compose/cmd/internals/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "turbo-compose",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		defer apiClient.Close()

		images, err := apiClient.ImageList(context.Background(), image.ListOptions{})
		if err != nil {
			panic(err)
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

		composeServices := listFromComposeFile()

		for _, service := range composeServices.Services {
			if !strings.HasPrefix(service.Image, "634375685434.dkr.ecr.us-east-1.amazonaws.com") {
				continue
			}

			name := strings.Split(service.Image, ":")[0]
			shortName := strings.Split(name, "/")[1]
			if img, ok := imgDict[name]; ok {
				fmt.Printf("✅ name: %v, created: %v\n", shortName, utils.TimeElapsed(time.Unix(img.Created, 0)))
			} else {
				fmt.Printf("❌ name: %v\n", shortName)
			}

		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func listFromComposeFile() DockerCompose {
	yamlFile, err := os.ReadFile("/Users/dlvhdr/code/komodor/mono/docker-compose.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	// Parse the YAML content into a DockerCompose struct
	var dockerCompose DockerCompose
	err = yaml.Unmarshal(yamlFile, &dockerCompose)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	// Print the parsed DockerCompose struct
	// fmt.Printf("Version: %s\n", dockerCompose.Version)
	for serviceName, service := range dockerCompose.Services {
		if strings.HasPrefix(service.Image, "634375685434.dkr.ecr.us-east-1.amazonaws.com") {
			// fmt.Printf("Service: %s", serviceName)
			// fmt.Printf("  Image: %s\n", service.Image)
		} else {
			delete(dockerCompose.Services, serviceName)
		}
	}

	return dockerCompose
}

type DockerCompose struct {
	Version  string
	Services map[string]Services
}

type Services struct {
	Image string
	Ports []string
}

func init() {
	// Here you will define your flags and configuration settings.
}

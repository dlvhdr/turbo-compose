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

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "turbo-compose",
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
		services := listFromComposeFile()
		for _, service := range services {
			fmt.Printf("%v\n", service.Image)
		}
		opts := make([]ServiceOption, 0)
		for _, service := range services {
			if !strings.HasPrefix(service.Image, "634375685434.dkr.ecr.us-east-1.amazonaws.com") {
				continue
			}
			// name := strings.Split(service.Image, ":")[0]
			// shortName := strings.Split(name, "/")[1]
			name := service.Name
			if img, ok := imgDict[name]; ok {
				opts = append(opts, ServiceOption{
					Name:       name,
					Image:      service.Image,
					LocalImage: &img,
				})
			} else {
				opts = append(opts, ServiceOption{
					Name:       name,
					Image:      service.Image,
					LocalImage: nil,
				})
			}
		}

		selection := make([]ServiceOption, 0)
		form := GetForm(opts, &selection)
		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}
		for _, service := range selection {
			fmt.Printf("Selected: %s\n", service.Name)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func listFromComposeFile() []Service {
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
	res := make([]Service, 0)
	for serviceName, service := range dockerCompose.Services {
		if strings.HasPrefix(service.Image, "634375685434.dkr.ecr.us-east-1.amazonaws.com") {
			res = append(res, Service{Name: serviceName, Image: service.Image})
		}
	}

	return res
}

type ServiceOption struct {
	Name       string
	LocalImage *image.Summary
	Image      string
}

type DockerCompose struct {
	Version  string
	Services map[string]Service
}

type Service struct {
	Name  string
	Image string
}

func init() {
	// Here you will define your flags and configuration settings.
}

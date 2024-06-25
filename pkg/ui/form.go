package ui

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/docker/docker/api/types/image"

	"github.com/dlvhdr/turbo-compose/pkg/docker"
	"github.com/dlvhdr/turbo-compose/pkg/utils"
)

type model struct {
	repo string
}

func NewModel(repo string) model {
	return model{repo: repo}
}

func GetForm(services []docker.ServiceOption, selection *[]docker.ServiceOption) *huh.Form {
	opts := make([]huh.Option[docker.ServiceOption], 0)
	sort.Slice(services, func(i, j int) bool {
		if services[i].LocalImage != nil && services[j].LocalImage == nil {
			return true
		}
		if services[i].LocalImage == nil && services[j].LocalImage != nil {
			return false
		}

		return services[i].Name < services[j].Name
	})
	for _, opt := range services {
		hasImage := opt.LocalImage != nil
		createdAt := ""
		if hasImage {
			createdAt = fmt.Sprintf("✅ %s", utils.TimeElapsed(time.Unix(opt.LocalImage.Created, 0)))
		} else {
			createdAt = "❌"
		}
		opts = append(opts, huh.NewOption(fmt.Sprintf("%s (%s)", opt.Name, createdAt), opt))
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[docker.ServiceOption]().
				Title("Select Services").
				Options(
					opts...,
				).
				Value(selection),
		),
	)
	return form
}

func (m model) makeOpts(definitions []docker.ServiceDefinition, localImgs map[string]image.Summary) []docker.ServiceOption {
	opts := make([]docker.ServiceOption, 0)
	for _, serviceDef := range definitions {
		if !strings.HasPrefix(serviceDef.Image, m.repo) {
			continue
		}
		name := serviceDef.Name
		opt := docker.ServiceOption{
			Name:  name,
			Image: serviceDef.Image,
		}
		if img, ok := localImgs[name]; ok {
			opt.LocalImage = &img
		} else {
			opt.LocalImage = nil
		}

		opts = append(opts, opt)
	}
	return opts
}

func (m model) Run() error {
	images, err := docker.GetLocalImages(m.repo)
	if err != nil {
		return err
	}
	services := docker.ListServicesFromComposeFile(m.repo)

	opts := m.makeOpts(services, images)

	selection := make([]docker.ServiceOption, 0)
	form := GetForm(opts, &selection)
	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range selection {
		fmt.Printf("Selected: %s\n", service.Name)
	}

	return nil
}

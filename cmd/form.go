package cmd

import (
	"github.com/charmbracelet/huh"
)

func GetForm(services []ServiceOption, selection *[]ServiceOption) *huh.Form {
	opts := make([]huh.Option[ServiceOption], 0)
	for _, opt := range services {
		opts = append(opts, huh.NewOption(opt.Name, opt))
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[ServiceOption]().
				Title("Select Services").
				Options(
					opts...,
				).
				Value(selection),
		),
	)
	return form
}

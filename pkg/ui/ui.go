package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/image"

	"github.com/dlvhdr/turbo-compose/pkg/docker"
	"github.com/dlvhdr/turbo-compose/pkg/utils"
)

type model struct {
	repo            string
	composeFilePath string
	list            *list.Model
	images          map[string]image.Summary
	services        []docker.ServiceDefinition
	options         []Option
	selection       []Option
}

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type Option struct {
	Name       string
	LocalImage *image.Summary
	Image      string
}

func (o Option) Title() string {
	return o.Name
}

func (o Option) Description() string {
	return o.Image
}

func (o Option) FilterValue() string {
	return o.Name
}

func NewModel(composeFilePath string, repo string) model {
	return model{repo: repo, composeFilePath: composeFilePath, selection: make([]Option, 0)}
}

func (m model) Init() tea.Cmd {
	return initCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		err error
		cmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case initMsg:
		// m, err = m.init()
		return m.init()
		if err != nil {
			panic(err)
		}
		return m, nil
	}

	if m.list != nil {
		l, cmd := m.list.Update(msg)
		m.list = &l
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	if m.list == nil {
		return "Loading..."
	}
	return m.list.View()
}

type initMsg struct{}

func initCmd() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}

func init() tea.Cmd {
	images, err := docker.GetLocalImages(m.repo)
	if err != nil {
		return model{}, nil
	}
	services := docker.ListServicesFromComposeFile(m.composeFilePath, m.repo)

	options = makeOpts(services, images)
	items := make([]list.Item, len(m.options))
	for i, opt := range m.options {
		items[i] = opt
	}

	return m, func() tea.Msg {
		return optionsMsg{
			images:   images,
			services: services,
			options:  options,
		}
	}
}

type optionsMsg struct {
	images   map[string]image.Summary
	services []docker.ServiceDefinition
	options  []Option
}

func (m model) makeOpts(definitions []docker.ServiceDefinition, localImgs map[string]image.Summary) []Option {
	opts := make([]Option, 0)
	for _, serviceDef := range definitions {
		if !strings.HasPrefix(serviceDef.Image, m.repo) {
			continue
		}
		name := serviceDef.Name
		opt := Option{
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

func (m model) getForm() *huh.Form {
	opts := make([]huh.Option[Option], 0)
	sort.Slice(m.options, func(i, j int) bool {
		if m.options[i].LocalImage != nil && m.options[j].LocalImage == nil {
			return true
		}
		if m.options[i].LocalImage == nil && m.options[j].LocalImage != nil {
			return false
		}

		return m.options[i].Name < m.options[j].Name
	})
	for _, opt := range m.options {
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
			huh.NewMultiSelect[Option]().
				Title("Select Services").
				Options(
					opts...,
				).
				Value(&m.selection),
		),
	)
	return form
}

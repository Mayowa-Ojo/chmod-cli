package ui

import (
	"math"
	"strings"

	"github.com/Mayowa-Ojo/chmod-cli/internal/common"
	"github.com/Mayowa-Ojo/chmod-cli/internal/generate"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	OptionsSection     = Section("options")
	PermissionsSection = Section("permissions")
)

type Section string

type Model struct {
	cursor      int
	section     Section
	options     *Options
	permissions *Permissions
	state       *generate.State
}

// Options store the state for selected options
type Options struct {
	values   []string
	selected string
	cursor   int
}

// Permissions store the state for selected permissions
type Permissions struct {
	blocks []PermissionsBlock
	cursor int
	values []string
}

// PermissionsBlock store the state for each permissions block
type PermissionsBlock struct {
	cursor   int
	selected []string
}

type PWDPermissionMsg string

type UpdateCommandMsg struct {
	User   generate.User
	Access generate.Access
	Active bool
}

func InitScreen() error {
	model := createModel()
	p := tea.NewProgram(model)

	return p.Start()
}

func createModel() tea.Model {
	optionValues := []string{"Verbose", "Changes", "Silent", "Default"}
	options := &Options{
		values:   optionValues,
		selected: optionValues[0],
	}

	blocks := make([]PermissionsBlock, 3)

	blocks[0].cursor = -1

	permissionValues := []string{"Read", "Write", "Execute"}
	permissions := &Permissions{
		values: permissionValues,
		blocks: blocks,
		cursor: -1,
	}

	state := generate.NewState()

	return Model{
		cursor:      0,
		section:     OptionsSection,
		options:     options,
		permissions: permissions,
		state:       state,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(getPWDPermission)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "down", "left", "right", "enter":
			if m.section == OptionsSection {
				if m.options.cursor+1 >= len(m.options.values) {
					switchSection(&m, msg.String())
				}

				m.options.updateOptions(msg.String())
			}

			if m.section == PermissionsSection {
				if m.permissions.blocks[m.permissions.cursor].cursor <= 0 {
					switchSection(&m, msg.String())
				}

				return m, m.permissions.updatePermissions(msg.String())
			}

		}

	case PWDPermissionMsg:
		m.state.PWD = string(msg)

	case UpdateCommandMsg:
		m.state.Users[msg.User][msg.Access] = msg.Active

		command := m.state.BuildCommand()

		m.state.Command = command
	}

	return m, nil
}

func switchSection(m *Model, msg string) {

	switch msg {
	case "up":
		m.section = OptionsSection
		m.permissions.cursor = -1
		m.permissions.blocks[0].cursor = -1
		m.options.cursor = len(m.options.values) - 1

	case "down":
		if m.section == PermissionsSection {
			break
		}

		m.section = PermissionsSection
		m.options.cursor = -2
		m.permissions.cursor = 0
	}
}

func updateCommand(user generate.User, access generate.Access, active bool) tea.Cmd {
	return func() tea.Msg {
		return UpdateCommandMsg{
			User:   user,
			Access: access,
			Active: active,
		}
	}
}

func getPWDPermission() tea.Msg {
	mode, err := generate.GetPWDMode()
	if err != nil {
		panic(err)
	}

	return PWDPermissionMsg(mode.String())
}

func (m Model) View() string {
	s := strings.Builder{}

	header := m.renderHeader()
	lists := m.permissions.renderPermissions()
	footer := m.renderFooter()

	s.WriteString(header)
	s.WriteString("\n")
	s.WriteString(m.options.renderOptions())
	s.WriteString("\n")
	s.WriteString(lists)
	s.WriteString("\n")
	s.WriteString(footer)

	return s.String()
}

func (o *Options) updateOptions(key string) bool {
	switch key {
	case "up":
		if o.cursor <= 0 {
			return false
		}

		o.cursor--
		return true

	case "down":
		if o.cursor >= len(o.values)-1 {
			return false
		}

		o.cursor++
		return true

	case "enter":
		o.selected = o.values[o.cursor]
	}

	return false
}

func (p *Permissions) updatePermissions(key string) tea.Cmd {
	if p.cursor < 0 {
		return nil
	}

	switch key {
	case "up":
		if p.blocks[p.cursor].cursor <= 0 {
			break
		}

		p.blocks[p.cursor].cursor--

	case "down":
		if p.blocks[p.cursor].cursor >= len(p.values)-1 {
			break
		}

		p.blocks[p.cursor].cursor++

	case "right":
		if p.cursor >= len(p.values)-1 {
			break
		}

		p.cursor++

	case "left":
		if p.cursor <= 0 {
			break
		}

		p.cursor--

	case "enter":
		item := p.values[p.blocks[p.cursor].cursor]
		selected := p.blocks[p.cursor].selected
		user := getBlockName(p.cursor)
		access := getAccessSymbol(item)

		// remove item if already exists [selected]
		if common.IncludesString(selected, item) {
			index := common.FindIndexString(selected, item)

			if math.Signbit(float64(index)) {
				break
			}

			selected = append(selected[:index], selected[index+1:]...)

			p.blocks[p.cursor].selected = selected

			return updateCommand(generate.User(user), generate.Access(access), false)
			// break
		}

		p.blocks[p.cursor].selected = append(selected, item)

		return updateCommand(generate.User(user), generate.Access(access), true)
	}

	return nil
}

func getBlockName(blockIndex int) string {
	switch blockIndex {
	case 0:
		return "owner"

	case 1:
		return "group"

	case 2:
		return "other"
	}

	return ""
}

func getAccessSymbol(access string) string {
	switch access {
	case "Read":
		return "r"

	case "Write":
		return "w"

	case "Execute":
		return "x"
	}

	return ""
}

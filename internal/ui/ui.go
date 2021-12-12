package ui

import (
	"math"
	"strings"

	"github.com/Mayowa-Ojo/chmod-cli/internal/common"
	"github.com/Mayowa-Ojo/chmod-cli/internal/generate"
	tea "github.com/charmbracelet/bubbletea"
)

const NumSections = 3

type Section int

const (
	OptionsSection Section = iota
	CommandModeSection
	PermissionsSection
)

func (s Section) String() string {
	return [...]string{"options", "command-mode", "permissions"}[s]
}

type Model struct {
	cursor      int
	section     Section
	options     *Options
	mode        *CommandMode
	permissions *Permissions
	state       *generate.State
}

// Options store the state for selected options
type Options struct {
	values   []string
	selected string
	cursor   int
}

// CommandMode stores the state for selected command mode
type CommandMode struct {
	values   []string
	selected string
	cursor   int
}

// PathType stores the state for selected path type
type PathType struct {
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

	commandModeValues := []string{"octal", "symbolic"}
	commandMode := &CommandMode{
		values:   commandModeValues,
		selected: commandModeValues[1],
		cursor:   -1,
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
		mode:        commandMode,
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
				m.options.updateOptions(msg.String())
			}

			if m.section == CommandModeSection {
				m.mode.updateCommandMode(msg.String())
			}

			if m.section == PermissionsSection {
				return m, m.permissions.updatePermissions(msg.String())
			}

		case "tab", " ", "shift+tab":
			switchSection(&m, msg.String())
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
	case "tab", " ":
		m.setSectionCursor(false)
		if m.cursor >= NumSections {
			m.cursor = -1
		}
		m.cursor++
		m.section = getSection(m.cursor)
		m.setSectionCursor(true)

	case "shift+tab":
		if m.cursor <= 0 {
			break
		}

		m.setSectionCursor(false)
		m.cursor--
		m.section = getSection(m.cursor)
		m.setSectionCursor(true)
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
	s.WriteString(m.mode.renderCommandMode())
	s.WriteString("\n")
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

func (c *CommandMode) updateCommandMode(key string) {
	switch key {
	case "left":
		if c.cursor <= 0 {
			break
		}
		c.cursor--

	case "right":
		if c.cursor >= len(c.values)-1 {
			break
		}
		c.cursor++

	case "enter":
		c.selected = c.values[c.cursor]
	}
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

func getSection(cursor int) Section {
	switch cursor {
	case 0:
		return OptionsSection

	case 1:
		return CommandModeSection

	case 2:
		return PermissionsSection
	}

	return 0
}

func (m Model) setSectionCursor(active bool) {
	switch m.section {
	case OptionsSection:
		if active {
			m.options.cursor = 0
			break
		}
		m.options.cursor = -1

	case CommandModeSection:
		if active {
			m.mode.cursor = 0
			break
		}
		m.mode.cursor = -1

	case PermissionsSection:
		if active {
			m.permissions.cursor = 0
			m.permissions.blocks[0].cursor = 0
			break
		}
		m.permissions.cursor = -1
	}
}

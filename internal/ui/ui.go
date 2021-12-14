package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Mayowa-Ojo/chmod-cli/internal/common"
	"github.com/Mayowa-Ojo/chmod-cli/internal/generate"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

const NumSections = 4
const ResetCommandDuration = time.Second * 3

type Section int

const (
	OptionsSection Section = iota
	CommandModeSection
	PathTypeSection
	PermissionsSection
)

func (s Section) String() string {
	return [...]string{"options", "command-mode", "path-type", "permissions"}[s]
}

type Model struct {
	cursor      int
	section     Section
	options     *Options
	mode        *CommandMode
	path        *PathType
	permissions *Permissions
	state       *generate.State
	keys        *KeyMap
	help        help.Model
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

type CopyCommandMsg struct{}

type ResetCommandMsg string

func InitScreen() error {
	model := createModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	return p.Start()
}

func createModel() tea.Model {
	optionValues := []string{"Verbose", "Changes", "Silent", "Default"}
	options := &Options{
		values:   optionValues,
		selected: optionValues[3],
	}

	commandModeValues := []string{"Octal", "Symbolic"}
	commandMode := &CommandMode{
		values:   commandModeValues,
		selected: commandModeValues[1],
		cursor:   -1,
	}

	pathTypeValues := []string{"File", "Directory"}
	pathType := &PathType{
		values:   pathTypeValues,
		selected: pathTypeValues[0],
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

	keyMap := NewKeyMap()

	help := help.NewModel()
	help.Width = 55

	return Model{
		cursor:      0,
		section:     OptionsSection,
		options:     options,
		mode:        commandMode,
		path:        pathType,
		permissions: permissions,
		state:       state,
		keys:        keyMap,
		help:        help,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(getPWDPermission)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit

		case "up", "down", "left", "right", "enter":
			if m.section == OptionsSection {
				return m, m.options.updateOptions(msg.String())
			}

			if m.section == CommandModeSection {
				return m, m.mode.updateCommandMode(msg.String())
			}

			if m.section == PathTypeSection {
				return m, m.path.updatePathType(msg.String())
			}

			if m.section == PermissionsSection {
				return m, m.permissions.updatePermissions(msg.String())
			}

		case "tab", " ", "shift+tab":
			switchSection(&m, msg.String())

		case "?":
			m.help.ShowAll = !m.help.ShowAll

		case "ctrl+c":
			if !strings.EqualFold(m.state.Command, "") {
				return m, copyCommand()
			}
		}

	case PWDPermissionMsg:
		m.state.PWD = string(msg)

	case UpdateCommandMsg:
		if !strings.EqualFold(string(msg.User), "") {
			m.state.Users[msg.User][msg.Access] = msg.Active
		}

		command := strings.Builder{}

		command.WriteString("chmod ")

		command.WriteString(fmt.Sprintf("%s ", getOptionFlag(&m)))

		if m.mode.selected == "Octal" {
			command.WriteString(m.state.BuildCommand(m.mode.selected))
		} else if m.path.selected == "Directory" {
			command.WriteString(fmt.Sprintf("d%s", m.state.BuildCommand(m.mode.selected)))
		} else {
			command.WriteString(fmt.Sprintf("-%s", m.state.BuildCommand(m.mode.selected)))
		}

		m.state.Command = command.String()

	case CopyCommandMsg:
		common.CopyToClipboard(m.state.Command)

		command := m.state.Command
		m.state.Command = "copied!"

		return m, resetCommand(command)

	case ResetCommandMsg:
		m.state.Command = string(msg)
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

func copyCommand() tea.Cmd {
	return func() tea.Msg {
		return CopyCommandMsg{}
	}
}

func resetCommand(cmd string) tea.Cmd {
	return tea.Tick(ResetCommandDuration, func(t time.Time) tea.Msg {
		return ResetCommandMsg(cmd)
	})
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
	help := m.help.View(m.keys)

	s.WriteString(header)
	s.WriteString("\n")
	s.WriteString(m.options.renderOptions())
	s.WriteString("\n")
	s.WriteString(m.mode.renderCommandMode())
	s.WriteString("\n\n")
	s.WriteString(m.path.renderPathType())
	s.WriteString("\n\n")
	s.WriteString(lists)
	s.WriteString("\n")
	s.WriteString(footer)
	s.WriteString("\n\n")
	s.WriteString(help)

	return s.String()
}

func (o *Options) updateOptions(key string) tea.Cmd {
	switch key {
	case "up":
		if o.cursor <= 0 {
			break
		}
		o.cursor--

	case "down":
		if o.cursor >= len(o.values)-1 {
			break
		}
		o.cursor++

	case "enter":
		o.selected = o.values[o.cursor]
		return updateCommand(generate.User(""), generate.Access(""), false)
	}

	return nil
}

func (c *CommandMode) updateCommandMode(key string) tea.Cmd {
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
		return updateCommand(generate.User(""), generate.Access(""), false)
	}

	return nil
}

func (p *PathType) updatePathType(key string) tea.Cmd {
	switch key {
	case "left":
		if p.cursor <= 0 {
			break
		}
		p.cursor--

	case "right":
		if p.cursor >= len(p.values)-1 {
			break
		}
		p.cursor++

	case "enter":
		p.selected = p.values[p.cursor]
		return updateCommand(generate.User(""), generate.Access(""), false)
	}

	return nil
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
		return PathTypeSection

	case 3:
		return PermissionsSection
	}

	return 0
}

func getOptionFlag(m *Model) string {
	switch m.options.selected {
	case "Verbose":
		return "--verbose"

	case "Changes":
		return "--changes"

	case "Silent":
		return "--silent"

	case "Default":
		return ""
	}

	return ""
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

	case PathTypeSection:
		if active {
			m.path.cursor = 0
			break
		}
		m.path.cursor = -1

	case PermissionsSection:
		if active {
			m.permissions.cursor = 0
			m.permissions.blocks[0].cursor = 0
			break
		}
		m.permissions.cursor = -1
	}
}

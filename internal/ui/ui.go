package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/Mayowa-Ojo/chmod-cli/internal/common"
	"github.com/Mayowa-Ojo/chmod-cli/internal/generate"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	green  = "#73F59F"
	red    = "#F25D94"
	purple = "#874BFD"
	yellow = "#FDE68A"
)

const (
	// radioActive   = "(o)"
	radioActive   = "(●)"
	radioInactive = "(o)"
	checkActive   = "[✓]"
	checkInactive = "[ ]"
	snowflake     = "❄ "
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

func (m Model) renderHeader() string {
	h := strings.Builder{}

	bannerContent := lipgloss.NewStyle().Foreground(lipgloss.Color("#FDE68A")).Render("chmod-cli v.0.1.0")
	bannerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true)
	banner := lipgloss.Place(
		54,
		5,
		lipgloss.Center,
		lipgloss.Center,
		bannerStyle.Render(bannerContent),
		lipgloss.WithWhitespaceChars(snowflake),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)

	h.WriteString("\n")
	h.WriteString(banner)
	h.WriteString("\n")
	h.WriteString(strings.Repeat("-", 54))
	h.WriteString("\n")
	h.WriteString(fmt.Sprintf("[PWD: %s]", m.state.PWD))
	h.WriteString("\n")
	h.WriteString(strings.Repeat("-", 54))
	h.WriteString("\n")

	return h.String()
}

func (m Model) renderFooter() string {
	footer := lipgloss.NewStyle().
		Width(55).
		Foreground(lipgloss.Color("#FDE68A")).
		Background(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		Padding(0, 1)
	footerContent := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Command: %s", m.state.Command))

	return footer.Render(footerContent)
}

func (o *Options) renderOptions() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#874BFD")).
		Padding(0, 3).Bold(true)

	option := lipgloss.NewStyle().Padding(0)
	activeOption := option.Copy().Foreground(lipgloss.Color("#F25D94"))
	options := strings.Builder{}

	for _, v := range o.values {
		focused := !math.Signbit(float64(o.cursor)) && o.values[o.cursor] == v
		active := o.selected == v

		if focused && active {
			options.WriteString(activeOption.Render(fmt.Sprintf("%s %s", radioActive, v)))
		} else if focused {
			options.WriteString(activeOption.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		} else if active {
			options.WriteString(fmt.Sprintf("%s %s", activeOption.Render(radioActive), v))
		} else {
			options.WriteString(option.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		}

		options.WriteString("\n")
	}

	optionsStyle := lipgloss.JoinVertical(
		lipgloss.Left,
		header.Render("Options"),
		options.String(),
	)

	return optionsStyle
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

func (p *Permissions) renderPermissions() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#874BFD")).
		Padding(0, 3).Bold(true)

	block := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		MarginRight(2).
		Height(5).
		Width(15)

	activeBlock := block.Copy().BorderForeground(lipgloss.Color("#874BFD"))

	blockItem := lipgloss.NewStyle().PaddingLeft(2)
	activeBlockItem := blockItem.Copy().Foreground(lipgloss.Color("#F25D94"))
	// focusedBlockItem := blockItem.Copy().Foreground(lipgloss.Color("#FDE68A"))

	var ownerBlock, groupBlock, otherBlock []string

	for _, v := range p.values {
		var currBlock PermissionsBlock
		if p.cursor >= 0 {
			currBlock = p.blocks[p.cursor]
		}

		{
			// owner block
			if len(ownerBlock) < 1 {
				ownerBlock = append(ownerBlock, blockItem.Copy().Foreground(lipgloss.Color(yellow)).Render("[Owner]"))
				ownerBlock = append(ownerBlock, blockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 0 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[0].selected, v)

			if focused && active {
				ownerBlock = append(ownerBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				ownerBlock = append(ownerBlock, fmt.Sprintf("%s %s", activeBlockItem.Render(checkActive), v))
			} else if focused {
				ownerBlock = append(ownerBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				ownerBlock = append(ownerBlock, blockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}

		{
			// group block
			if len(groupBlock) < 1 {
				groupBlock = append(groupBlock, blockItem.Copy().Foreground(lipgloss.Color(yellow)).Render("[Group]"))
				groupBlock = append(groupBlock, blockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 1 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[1].selected, v)

			if focused && active {
				groupBlock = append(groupBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				groupBlock = append(groupBlock, fmt.Sprintf("%s %s", activeBlockItem.Render(checkActive), v))
			} else if focused {
				groupBlock = append(groupBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				groupBlock = append(groupBlock, blockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}

		{
			// other block
			if len(otherBlock) < 1 {
				otherBlock = append(otherBlock, blockItem.Copy().Foreground(lipgloss.Color(yellow)).Render("[Other]"))
				otherBlock = append(otherBlock, blockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 2 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[2].selected, v)

			if focused && active {
				otherBlock = append(otherBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				otherBlock = append(otherBlock, fmt.Sprintf("%s %s", activeBlockItem.Render(checkActive), v))
			} else if focused {
				otherBlock = append(otherBlock, activeBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				otherBlock = append(otherBlock, blockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}
	}

	var (
		ownerBlockStyle = block
		groupBlockStyle = block
		otherBlockStyle = block
	)

	if p.cursor == 0 {
		ownerBlockStyle = activeBlock
	}

	if p.cursor == 1 {
		groupBlockStyle = activeBlock
	}

	if p.cursor == 2 {
		otherBlockStyle = activeBlock
	}

	blocks := lipgloss.JoinHorizontal(
		lipgloss.Top,
		ownerBlockStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			ownerBlock...,
		)),
		groupBlockStyle.Copy().Render(lipgloss.JoinVertical(
			lipgloss.Left,
			groupBlock...,
		)),
		otherBlockStyle.Copy().Render(lipgloss.JoinVertical(
			lipgloss.Left,
			otherBlock...,
		)),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header.Render("Permissions"), blocks)
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

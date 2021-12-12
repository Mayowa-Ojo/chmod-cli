package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/Mayowa-Ojo/chmod-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderHeader() string {
	styles := GetStyles()
	h := strings.Builder{}

	banner := styles.Banner

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
	styles := GetStyles()

	footer := styles.Footer
	footerContent := styles.FooterContent.Render(fmt.Sprintf("Command: %s", m.state.Command))

	return footer.Render(footerContent)
}

func (o *Options) renderOptions() string {
	styles := GetStyles()

	options := strings.Builder{}

	for _, v := range o.values {
		focused := !math.Signbit(float64(o.cursor)) && o.values[o.cursor] == v
		active := o.selected == v

		if focused && active {
			options.WriteString(styles.OptionsActiveItem.Render(fmt.Sprintf("%s %s", radioActive, v)))
		} else if focused {
			options.WriteString(styles.OptionsActiveItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		} else if active {
			options.WriteString(fmt.Sprintf("%s %s", styles.OptionsActiveItem.Render(radioActive), v))
		} else {
			options.WriteString(styles.OptionsItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		}

		options.WriteString("\n")
	}

	return styles.OptionsContainer(options)
}

func (c *CommandMode) renderCommandMode() string {
	styles := GetStyles()

	modes := []string{}

	for _, v := range c.values {
		focused := !math.Signbit(float64(c.cursor)) && c.values[c.cursor] == v
		active := c.selected == v

		if focused && active {
			modes = append(modes, styles.CommandModeActiveItem.Render(fmt.Sprintf("%s %s", radioActive, v)))
		} else if focused {
			modes = append(modes, styles.CommandModeActiveItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		} else if active {
			modes = append(modes, fmt.Sprintf("%s %s", styles.CommandModeActiveItem.Render(radioActive), v))
		} else {
			modes = append(modes, styles.CommandModeItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		}
	}

	return styles.CommandModeContainer(modes...)
}

func (p *PathType) renderPathType() string {
	styles := GetStyles()

	paths := []string{}

	for _, v := range p.values {
		focused := !math.Signbit(float64(p.cursor)) && p.values[p.cursor] == v
		active := p.selected == v

		if focused && active {
			paths = append(paths, styles.PathTypeActiveItem.Render(fmt.Sprintf("%s %s", radioActive, v)))
		} else if focused {
			paths = append(paths, styles.PathTypeActiveItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		} else if active {
			paths = append(paths, fmt.Sprintf("%s %s", styles.PathTypeActiveItem.Render(radioActive), v))
		} else {
			paths = append(paths, styles.PathTypeItem.Render(fmt.Sprintf("%s %s", radioInactive, v)))
		}
	}

	return styles.PathTypeContainer(paths...)
}

func (p *Permissions) renderPermissions() string {
	styles := GetStyles()

	var ownerBlock, groupBlock, otherBlock []string

	for _, v := range p.values {
		var currBlock PermissionsBlock
		if p.cursor >= 0 {
			currBlock = p.blocks[p.cursor]
		}

		{
			if len(ownerBlock) < 1 {
				ownerBlock = append(ownerBlock, styles.PermissionsBlockItem.Copy().Foreground(lipgloss.Color(ColorYellow)).Render("[Owner]"))
				ownerBlock = append(ownerBlock, styles.PermissionsBlockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 0 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[0].selected, v)

			if focused && active {
				ownerBlock = append(ownerBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				ownerBlock = append(ownerBlock, fmt.Sprintf("%s %s", styles.PermissionsActiveBlockItem.Render(checkActive), v))
			} else if focused {
				ownerBlock = append(ownerBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				ownerBlock = append(ownerBlock, styles.PermissionsBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}

		{
			if len(groupBlock) < 1 {
				groupBlock = append(groupBlock, styles.PermissionsBlockItem.Copy().Foreground(lipgloss.Color(ColorYellow)).Render("[Group]"))
				groupBlock = append(groupBlock, styles.PermissionsBlockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 1 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[1].selected, v)

			if focused && active {
				groupBlock = append(groupBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				groupBlock = append(groupBlock, fmt.Sprintf("%s %s", styles.PermissionsActiveBlockItem.Render(checkActive), v))
			} else if focused {
				groupBlock = append(groupBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				groupBlock = append(groupBlock, styles.PermissionsBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}

		{
			if len(otherBlock) < 1 {
				otherBlock = append(otherBlock, styles.PermissionsBlockItem.Copy().Foreground(lipgloss.Color(ColorYellow)).Render("[Other]"))
				otherBlock = append(otherBlock, styles.PermissionsBlockItem.Render(strings.Repeat("-", 7)))
			}

			focused := p.cursor == 2 && p.values[currBlock.cursor] == v
			active := common.IncludesString(p.blocks[2].selected, v)

			if focused && active {
				otherBlock = append(otherBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkActive, v)))
			} else if active {
				otherBlock = append(otherBlock, fmt.Sprintf("%s %s", styles.PermissionsActiveBlockItem.Render(checkActive), v))
			} else if focused {
				otherBlock = append(otherBlock, styles.PermissionsActiveBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			} else {
				otherBlock = append(otherBlock, styles.PermissionsBlockItem.Render(fmt.Sprintf("%s %s", checkInactive, v)))
			}
		}
	}

	var (
		ownerBlockStyle = styles.PermissionsBlock
		groupBlockStyle = styles.PermissionsBlock
		otherBlockStyle = styles.PermissionsBlock
	)

	if p.cursor == 0 {
		ownerBlockStyle = styles.PermissionsActiveBlock
	}

	if p.cursor == 1 {
		groupBlockStyle = styles.PermissionsActiveBlock
	}

	if p.cursor == 2 {
		otherBlockStyle = styles.PermissionsActiveBlock
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.PermissionsHeader.Render("Permissions"),
		blocks,
	)
}

package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// colors
type Color = string

const (
	ColorGreen       = Color("#73F59F")
	ColorRed         = Color("#F25D94")
	ColorPurple      = Color("#874BFD")
	ColorYellow      = Color("#FDE68A")
	ColorGray50      = Color("#FFF7DB")
	ColorSubtleLight = Color("#D9DCCF")
	ColorSubtleDark  = Color("#383838")
)

// symbols
const (
	// radioActive   = "(o)"
	radioActive   = "(●)"
	radioInactive = "(o)"
	checkActive   = "[✓]"
	checkInactive = "[ ]"
	snowflake     = "❄ "
)

type Styles struct {
	Banner        string
	BannerContent lipgloss.Style
	BannerText    lipgloss.Style

	Footer        lipgloss.Style
	FooterContent lipgloss.Style

	OptionsContainer  func(opts strings.Builder) string
	OptionsHeader     lipgloss.Style
	OptionsItem       lipgloss.Style
	OptionsActiveItem lipgloss.Style

	CommandModeContainer  func(modes ...string) string
	CommandModeHeader     lipgloss.Style
	CommandModeItem       lipgloss.Style
	CommandModeActiveItem lipgloss.Style

	PathTypeContainer  func(paths ...string) string
	PathTypeHeader     lipgloss.Style
	PathTypeItem       lipgloss.Style
	PathTypeActiveItem lipgloss.Style

	PermissionsHeader          lipgloss.Style
	PermissionsBlock           lipgloss.Style
	PermissionsActiveBlock     lipgloss.Style
	PermissionsBlockItem       lipgloss.Style
	PermissionsActiveBlockItem lipgloss.Style
}

func GetStyles() *Styles {
	s := new(Styles)

	s.BannerText = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorYellow))

	s.BannerContent = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorPurple)).
		Padding(0, 1).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true)

	s.Banner = lipgloss.Place(
		54,
		5,
		lipgloss.Center,
		lipgloss.Center,
		s.BannerContent.Render(s.BannerText.Render("chmod-cli v.0.1.0")),
		lipgloss.WithWhitespaceChars(snowflake),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: ColorSubtleLight, Dark: ColorSubtleDark}),
	)

	s.Footer = lipgloss.NewStyle().
		Width(55).
		Foreground(lipgloss.Color(ColorYellow)).
		Background(lipgloss.AdaptiveColor{Light: ColorSubtleLight, Dark: ColorSubtleDark}).
		Padding(0, 1)

	s.FooterContent = lipgloss.NewStyle().Bold(true)

	s.OptionsHeader = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray50)).
		Background(lipgloss.Color(ColorPurple)).
		Padding(0, 3).Bold(true)

	s.OptionsItem = lipgloss.NewStyle().Padding(0)
	s.OptionsActiveItem = s.OptionsItem.Copy().Foreground(lipgloss.Color(ColorRed))

	s.OptionsContainer = func(opts strings.Builder) string {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			s.OptionsHeader.Render("Options"),
			opts.String(),
		)
	}

	s.CommandModeHeader = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray50)).
		Background(lipgloss.Color(ColorPurple)).
		Padding(0, 3).Bold(true)

	s.CommandModeItem = lipgloss.NewStyle().Padding(0)

	s.CommandModeActiveItem = s.CommandModeItem.Copy().Foreground(lipgloss.Color(ColorRed))

	s.CommandModeContainer = func(modes ...string) string {
		// hacky way to add spacing between horizontal elements
		// because margin's not giving expected behavior
		modes = append(modes, lipgloss.NewStyle().Render("  "))

		modes[1], modes[2] = modes[2], modes[1]

		return lipgloss.JoinVertical(
			lipgloss.Left,
			s.CommandModeHeader.Render("Command Mode"),
			lipgloss.JoinHorizontal(lipgloss.Top, modes...),
		)
	}

	s.PathTypeHeader = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray50)).
		Background(lipgloss.Color(ColorPurple)).
		Padding(0, 3).Bold(true)

	s.PathTypeItem = lipgloss.NewStyle().Padding(0)

	s.PathTypeActiveItem = s.PathTypeItem.Copy().Foreground(lipgloss.Color(ColorRed))

	s.PathTypeContainer = func(paths ...string) string {
		paths = append(paths, lipgloss.NewStyle().Render("  "))

		paths[1], paths[2] = paths[2], paths[1]

		return lipgloss.JoinVertical(
			lipgloss.Left,
			s.PathTypeHeader.Render("Path Type"),
			lipgloss.JoinHorizontal(lipgloss.Top, paths...),
		)
	}

	s.PermissionsHeader = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray50)).
		Background(lipgloss.Color(ColorPurple)).
		Padding(0, 3).Bold(true)

	s.PermissionsBlock = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		MarginRight(2).
		Height(5).
		Width(15)

	s.PermissionsActiveBlock = s.PermissionsBlock.Copy().BorderForeground(lipgloss.Color(ColorPurple))

	s.PermissionsBlockItem = lipgloss.NewStyle().PaddingLeft(2)

	s.PermissionsActiveBlockItem = s.PermissionsBlockItem.Copy().Foreground(lipgloss.Color(ColorRed))

	return s
}

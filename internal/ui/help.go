package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	TabUp   key.Binding
	TabDown key.Binding
	Select  key.Binding
	Copy    key.Binding
	Quit    key.Binding
	Help    key.Binding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
		),
		TabUp: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("⇧ + tab", "prev section"),
		),
		TabDown: key.NewBinding(
			key.WithKeys("tab/space"),
			key.WithHelp("tab/space", "next section"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Copy: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "copy command"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.TabDown, k.TabUp, k.Select, k.Copy},
		{k.Help, k.Quit},
	}
}

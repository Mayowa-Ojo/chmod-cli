package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	var msg tea.KeyMsg
	is := require.New(t)

	t.Run("test update options", func(t *testing.T) {
		t.Skip()
		model := createModel()
		msg = tea.KeyMsg{
			Type:  tea.KeyDown,
			Runes: nil,
			Alt:   false,
		}

		model, cmd := model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).options.cursor; cursor != 1 {
			t.Errorf("Expected cursor to be '1', instead got '%d'", cursor)
		}

		msg = tea.KeyMsg{
			Type:  tea.KeyUp,
			Runes: nil,
			Alt:   false,
		}

		model, cmd = model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).options.cursor; cursor != 0 {
			t.Errorf("Expected cursor to be '0', instead got '%d'", cursor)
		}

		msg = tea.KeyMsg{
			Type:  tea.KeyEnter,
			Runes: nil,
			Alt:   false,
		}

		model, cmd = model.Update(msg)
		is.NotNil(cmd)
		is.NotNil(model)

		if selected := model.(Model).options.selected; selected != "Verbose" {
			t.Errorf("Expected selected to be 'Verbose', instead got '%s'", selected)
		}
	})

	t.Run("test update command-mode", func(t *testing.T) {
		t.Skip()
		model := createModel()
		msg = tea.KeyMsg{
			Type:  tea.KeyTab,
			Runes: nil,
			Alt:   false,
		}

		model, cmd := model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).cursor; cursor != 1 {
			t.Errorf("Expected cursor to be '1', instead got '%d'", cursor)
		}

		if section := model.(Model).section; section != CommandModeSection {
			t.Errorf("Expected section to be 'command-mode', instead got '%d'", section)
		}

		msg = tea.KeyMsg{
			Type:  tea.KeyRight,
			Runes: nil,
			Alt:   false,
		}

		model, cmd = model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).mode.cursor; cursor != 1 {
			t.Errorf("Expected cursor to be '1', instead got '%d'", cursor)
		}

		msg = tea.KeyMsg{
			Type:  tea.KeyEnter,
			Runes: nil,
			Alt:   false,
		}

		model, cmd = model.Update(msg)
		is.NotNil(cmd)
		is.NotNil(model)

		if selected := model.(Model).mode.selected; selected != "Symbolic" {
			t.Errorf("Expected selected to be 'Symbolic', instead got '%s'", selected)
		}
	})

	t.Run("test update path-type", func(t *testing.T) {
		// t.Skip()
		model := createModel()
		msg := tea.KeyMsg{
			Type:  tea.KeyTab,
			Runes: nil,
			Alt:   false,
		}

		model, _ = model.Update(msg)
		model, cmd := model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).cursor; cursor != 2 {
			t.Errorf("Expected cursor to be '2', instead got '%d'", cursor)
		}

		if section := model.(Model).section; section != PathTypeSection {
			t.Errorf("Expected section to be 'path-type', instead got '%s'", section)
		}

		msg = tea.KeyMsg{
			Type:  tea.KeyRight,
			Runes: nil,
			Alt:   false,
		}

		model, cmd = model.Update(msg)
		is.Nil(cmd)
		is.NotNil(model)

		if cursor := model.(Model).path.cursor; cursor != 1 {
			t.Errorf("Expected cursor to be '1', instead got '%d'", cursor)
		}
	})
}

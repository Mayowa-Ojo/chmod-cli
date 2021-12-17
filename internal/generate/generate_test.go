package generate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPWDMode(t *testing.T) {
	is := require.New(t)

	mode, err := GetPWDMode()

	is.NoError(err)
	is.NotNil(mode)

	expected := "drwxrwxr-x"

	if mode.String() != expected {
		t.Errorf("Expected mode to be equal to '%s', instead got '%s'", expected, mode.String())
	}
}

func TestSortKeys(t *testing.T) {
	is := require.New(t)
	s := NewState()

	m := map[Access]bool{
		ExecuteAccess: false,
		WriteAccess:   false,
		ReadAccess:    true,
	}

	keys := s.SortKeys(m)
	is.NotEmpty(keys)

	if strings.Join([]string{string(keys[0]), string(keys[1]), string(keys[2])}, "") != "rwx" {
		t.Errorf("Expected keys to be sorted, instead got %v", keys)
	}
}

func TestToOctal(t *testing.T) {
	// is := require.New(t)

	cmd := "rw-rwxr-x"
	expected := "675"

	if got := toOctal(cmd); got != expected {
		t.Errorf("Expected octal value to be '%s', instead got '%s'", expected, got)
	}
}

func TestBuildCommand(t *testing.T) {
	is := require.New(t)

	s := NewState()

	s.Users[Owner][ReadAccess] = true
	s.Users[Owner][WriteAccess] = true

	s.Users[Group][ReadAccess] = true
	s.Users[Group][WriteAccess] = true
	s.Users[Group][ExecuteAccess] = true

	s.Users[Other][ReadAccess] = true
	s.Users[Other][ExecuteAccess] = true

	cmd := s.BuildCommand("Symbolic")
	is.NotEmpty(cmd)

	expected := "rw-rwxr-x"
	if cmd != expected {
		t.Errorf("Expected command to be '%s', instead got '%s'", expected, cmd)
	}
}

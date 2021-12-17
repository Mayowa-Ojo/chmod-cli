package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncludesString(t *testing.T) {

	slice := []string{"red", "blue", "green"}
	el := "blue"

	if exists := IncludesString(slice, el); !exists {
		t.Errorf("Expected 'exists' to be true, instead got '%v'", exists)
	}

	el = "orange"

	if exists := IncludesString(slice, el); exists {
		t.Errorf("Expected 'exists' to be false, instead got, %v", exists)
	}
}

func TestFindIndexString(t *testing.T) {
	is := require.New(t)

	slice := []string{"red", "blue", "green"}
	el := "blue"

	index := FindIndexString(slice, el)
	is.NotEmpty(index)

	if index != 1 {
		t.Errorf("Expected index to be '1', instead got '%d'", index)
	}
}

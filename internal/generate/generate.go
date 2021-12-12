package generate

import (
	"io/fs"
	"os"
	"sort"
	"strings"
)

type Access string

const (
	ReadAccess    = Access("r")
	WriteAccess   = Access("w")
	ExecuteAccess = Access("x")
)

type User string

const (
	Owner = User("owner")
	Group = User("group")
	Other = User("other")
)

type State struct {
	Users   map[User]map[Access]bool
	Command string
	PWD     string
}

func NewState() *State {
	return &State{
		Users: map[User]map[Access]bool{
			Owner: {
				ReadAccess:    false,
				WriteAccess:   false,
				ExecuteAccess: false,
			},
			Group: {
				ReadAccess:    false,
				WriteAccess:   false,
				ExecuteAccess: false,
			},
			Other: {
				ReadAccess:    false,
				WriteAccess:   false,
				ExecuteAccess: false,
			},
		},
	}
}

func GetPWDMode() (fs.FileMode, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	stat, err := os.Stat(pwd)
	if err != nil {
		return 0, err
	}

	return stat.Mode(), nil
}

func (s *State) BuildCommand() string {
	command := strings.Builder{}
	owner := strings.Builder{}
	group := strings.Builder{}
	other := strings.Builder{}
	// fmt.Printf("state: %+v\n", s.Users)

	sortedOwnerKeys := s.sortKeys(s.Users[Owner])
	sortedGroupKeys := s.sortKeys(s.Users[Group])
	sortedOtherKeys := s.sortKeys(s.Users[Other])

	for _, v := range sortedOwnerKeys {
		if s.Users[Owner][v] {
			owner.WriteString(string(v))
		} else {
			owner.WriteString("-")
		}
	}

	for _, v := range sortedGroupKeys {
		if s.Users[Group][v] {
			group.WriteString(string(v))
		} else {
			group.WriteString("-")
		}
	}

	for _, v := range sortedOtherKeys {
		if s.Users[Other][v] {
			other.WriteString(string(v))
		} else {
			other.WriteString("-")
		}
	}

	command.WriteString(owner.String())
	command.WriteString(group.String())
	command.WriteString(other.String())

	return command.String()
}

func (s *State) sortKeys(m map[Access]bool) []Access {
	keys := make([]Access, len(m))
	i := 0

	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	return keys
}

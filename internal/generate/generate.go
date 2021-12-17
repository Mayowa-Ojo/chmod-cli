package generate

import (
	"errors"
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

func (s *State) BuildCommand(mode string) string {
	command := strings.Builder{}
	owner := strings.Builder{}
	group := strings.Builder{}
	other := strings.Builder{}

	sortedOwnerKeys := s.SortKeys(s.Users[Owner])
	sortedGroupKeys := s.SortKeys(s.Users[Group])
	sortedOtherKeys := s.SortKeys(s.Users[Other])

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

	if mode == "Octal" {
		return toOctal(command.String())
	}

	return command.String()
}

func toOctal(cmd string) string {
	if len(cmd) != 9 {
		panic(errors.New("invalid chmod string"))
	}

	mapBinaryToDecimal := map[string]string{
		"000": "0",
		"001": "1",
		"010": "2",
		"011": "3",
		"100": "4",
		"101": "5",
		"110": "6",
		"111": "7",
	}

	binary := strings.Builder{}

	for _, s := range cmd {
		if string(s) == "-" {
			binary.WriteString("0")
		} else {
			binary.WriteString("1")
		}
	}

	var owner, group, other string

	owner = binary.String()[:3]
	group = binary.String()[3:6]
	other = binary.String()[6:]

	octal := strings.Builder{}

	octal.WriteString(mapBinaryToDecimal[owner])
	octal.WriteString(mapBinaryToDecimal[group])
	octal.WriteString(mapBinaryToDecimal[other])

	return octal.String()
}

func (s *State) SortKeys(m map[Access]bool) []Access {
	keys := make([]Access, len(m))
	i := 0

	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	return keys
}

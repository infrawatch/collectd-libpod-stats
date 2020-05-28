package cgroups

import (
	"github.com/pkg/errors"
)

//Memory memory control group controller
type Memory struct {
	path    string
	cgroup2 bool
}

//NewMemory create new memory cgroup control type
func NewMemory(path string) (*Memory, error) {
	memory := &Memory{
		path: path,
	}

	var err error
	memory.cgroup2, err = IsCgroup2UnifiedMode()
	if err != nil {
		return nil, errors.Wrap(err, "determining cgroup version")
	}
	return memory, nil
}

//Stats get memory stats
func (m *Memory) Stats() (uint64, error) {
	return 0, nil
}

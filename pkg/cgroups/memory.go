package cgroups

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

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
	if m.cgroup2 {
		return m.statsV2()
	}
	return m.statsV1()
}

func (m *Memory) statsV1() (uint64, error) {
	res, err := readFileAsUint64(filepath.Join(m.path, "memory.usage_in_bytes"))
	if err != nil {
		return 0, errors.Wrapf(err, "retrieving memory stats cgroup v1")
	}
	return res, nil
}

func (m *Memory) statsV2() (uint64, error) {
	p := filepath.Join(m.path, "memory.current")

	stat, err := ioutil.ReadFile(p)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.Atoi(strings.TrimSpace(string(stat)))

	if err != nil {
		return 0, errors.Wrapf(err, "retrieving memory stats cgroup v2")
	}

	return uint64(ret), nil
}

package cgroups

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"collectd.org/api"
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
func (m *Memory) Stats() (api.Value, error) {
	if m.cgroup2 {
		return m.statsV2()
	}
	return m.statsV1()
}

func (m *Memory) statsV1() (api.Value, error) {
	res, err := readFileAsUint64(filepath.Join(m.path, "memory.usage_in_bytes"))
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving memory stats cgroup v1")
	}
	return api.Gauge(res), nil
}

func (m *Memory) statsV2() (api.Value, error) {
	p := filepath.Join(m.path, "memory.current")

	stat, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	ret, err := strconv.Atoi(strings.TrimSpace(string(stat)))

	if err != nil {
		return nil, errors.Wrapf(err, "retrieving memory stats cgroup v2")
	}

	return api.Gauge(ret), nil
}

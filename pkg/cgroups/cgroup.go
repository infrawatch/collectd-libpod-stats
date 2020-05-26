package cgroups

import (
	"bufio"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	cgroupRoot = "/sys/fs/cgroup"
)

//ControlType supported cgroup controller types
type ControlType int

const (
	CPUAcctT ControlType = iota
	MemoryT
)

func (ct ControlType) String() string {
	return []string{"cpuacct", "memory"}[ct]
}

//CgroupControl represents a cgroup controller
type CgroupControl interface {
	Stats() ([]byte, error)
}

func readFileAsUint64(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, errors.Wrapf(err, "open %s", path)
	}
	v := cleanString(string(data))
	if v == "max" {
		return math.MaxUint64, nil
	}
	ret, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return ret, errors.Wrapf(err, "parse %s from %s", v, path)
	}
	return ret, nil
}

func readCgroup2MapPath(path string) (map[string][]string, error) {
	ret := map[string][]string{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ret, nil
		}
		return nil, errors.Wrapf(err, "open file %s", path)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		ret[parts[0]] = parts[1:]
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "parsing file %s", path)
	}
	return ret, nil
}

func cleanString(s string) string {
	return strings.Trim(s, "\n")
}

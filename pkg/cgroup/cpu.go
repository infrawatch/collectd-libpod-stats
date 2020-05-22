package cgroup

import (
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

// GetSystemCPUUsage returns the system usage for all the cgroups
func GetSystemCPUUsage() (uint64, error) {
	cgroupv2, err := IsCgroup2UnifiedMode()
	if err != nil {
		return 0, err
	}
	if !cgroupv2 {
		p := filepath.Join(cgroupRoot, CPUAcct, "cpuacct.usage")
		return readFileAsUint64(p)
	}

	files, err := ioutil.ReadDir(cgroupRoot)
	if err != nil {
		return 0, errors.Wrapf(err, "read directory %q", cgroupRoot)
	}
	var total uint64
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		p := filepath.Join(cgroupRoot, file.Name(), "cpu.stat")

		values, err := readCgroup2MapPath(p)
		if err != nil {
			return 0, err
		}

		if val, found := values["usage_usec"]; found {
			v, err := strconv.ParseUint(cleanString(val[0]), 10, 0)
			if err != nil {
				return 0, err
			}
			total += v * 1000
		}

	}
	return total, nil
}

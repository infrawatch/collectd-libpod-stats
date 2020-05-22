/*Package cgroup includes objects and functions for interfacing with cgroups on the host OS

As of 22 May 2020, there exist only two
1. v1 (<= RHEL7/Fedora30/CentOS7)
2. v2 (>= RHEL8/Fedora31/CentOS8)
*/
package cgroup

import (
	"bufio"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	cgroupRoot = "/sys/fs/cgroup"
	// CPU is the cpu controller
	CPU = "cpu"
	// CPUAcct is the cpuacct controller
	CPUAcct = "cpuacct"
	// Memory is the memory controller
	Memory = "memory"
)

//Cgroup represents a cgroup version
type Cgroup interface {
}

var (
	isUnifiedOnce sync.Once
	isUnified     bool
	isUnifiedErr  error
)

// IsCgroup2UnifiedMode returns whether we are running in cgroup 2 cgroup2 mode.
func IsCgroup2UnifiedMode() (bool, error) {
	isUnifiedOnce.Do(func() {
		var st syscall.Statfs_t
		if err := syscall.Statfs("/sys/fs/cgroup", &st); err != nil {
			isUnified, isUnifiedErr = false, err
		} else {
			isUnified, isUnifiedErr = st.Type == unix.CGROUP2_SUPER_MAGIC, nil
		}
	})
	return isUnified, isUnifiedErr
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

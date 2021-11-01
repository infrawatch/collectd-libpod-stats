package virt

/* Package virt contains objects and functions for handling cgroups and containers
on host OS

As of 22 May 2020:
1. cgroup v1 (<= RHEL7/Fedora30/CentOS7)
2. cgroup v2 (>= RHEL8/Fedora31/CentOS8)
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/infrawatch/collectd-libpod-stats/pkg/cgroups"
	"github.com/infrawatch/collectd-libpod-stats/pkg/containers"
)

const (
	//container path relative to user
	relativeContainersPath string = "containers/storage/overlay-containers/containers.json"
)

//MetricMatrix holds stats for each container according to
//control type. Usage is: map[container label]map[control type]data
type MetricMatrix map[string]map[cgroups.ControlType]uint64

//ContainersStats retrieves stats in specified cgroup controllers for all containers on host
func ContainersStats(cgroupControls ...cgroups.ControlType) (MetricMatrix, error) {
	cgroup2, err := cgroups.IsCgroup2UnifiedMode()
	if err != nil {
		return nil, errors.Wrapf(err, "determing cgroup version")
	}
	uid := os.Geteuid()

	retMatrix := MetricMatrix{}

	cMap, err := getContainers()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving host containers")
	}

	for cLabel, c := range cMap {
		retMatrix[cLabel] = map[cgroups.ControlType]uint64{}
		for _, control := range cgroupControls {
			ctrlPath, err := genContainerCgroupPath(control, c.ID, cgroup2, uid)
			if err != nil {
				return nil, err
			}

			cgCtrl, err := cgroups.CgroupControlFactory(control, ctrlPath)
			if err != nil {
				return nil, err
			}

			stat, err := cgCtrl.Stats()
			if err != nil && err != cgroups.ErrDoesNotExist {
				return nil, err
			}

			if err != cgroups.ErrDoesNotExist {
				retMatrix[cLabel][control] = stat
			}
		}
	}
	return retMatrix, nil
}

//getContainers returns map with containers created on host indexed by name
func getContainers() (map[string]*containers.Container, error) {
	/*
		libpod stores container related information in one of two places:
		1. /var/lib/containers/storage/overlay-containers (root)
		2. $HOME/.local/share/containers/storage/overlay-containers (rootless)
	*/

	uid := os.Geteuid()
	containersPath := filepath.Join("/var/lib", relativeContainersPath)
	if uid != 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving user home directory")
		}

		containersPath = filepath.Join(home, ".local/share", relativeContainersPath)
	}

	containersJSON, err := ioutil.ReadFile(containersPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading libpod container file")
	}

	containerList, err := containers.NewListFromJSON(containersJSON)
	if err != nil {
		return nil, errors.Wrap(err, "loading container json")
	}

	containerMap := make(map[string]*containers.Container)
	for _, c := range containerList {
		for _, name := range c.Names {
			containerMap[name] = c
		}
	}

	return containerMap, nil
}

func genContainerCgroupPath(ctype cgroups.ControlType, id string, cgroup2 bool, uid int) (string, error) {

	path, err := filepath.Abs("/sys/fs/cgroup")
	if err != nil {
		return "", errors.Wrapf(err, "retrieving cgroup root path")
	}

	if cgroup2 {
		if uid != 0 {
			path = filepath.Join(path, fmt.Sprintf("user.slice/user-%d.slice/user@%d.service/user.slice/libpod-%s.scope", uid, uid, id))
		} else {
			path = filepath.Join(path, fmt.Sprintf("machine.slice/libpod-%s.scope", id))
		}
	} else {
		if uid != 0 {
			return "", fmt.Errorf("rootless cgroups require Cgroups V2")
		}
		path = filepath.Join(path, fmt.Sprintf("%s/machine.slice/libpod-%s.scope", ctype.String(), id))
	}
	return path, nil
}

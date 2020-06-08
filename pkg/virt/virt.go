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
	"github.com/pleimer/collectd-libpod-stats/pkg/cgroups"
	"github.com/pleimer/collectd-libpod-stats/pkg/containers"
)

const (
	//container path relative to user
	relativeContainersPath string = "containers/storage/overlay-containers/containers.json"

	//libpod cgroup name template
	containerCgroupTemplate string = "libpod-%s.scope"
)

//MetricMatrix holds stats for each container according to
//control type
type MetricMatrix map[string]map[cgroups.ControlType]uint64

//ContainersStats retrieves stats in specified cgroup controllers for all containers on host
func ContainersStats(cgroupControls ...cgroups.ControlType) (MetricMatrix, error) {
	retMatrix := MetricMatrix{}

	cMap, err := getContainers()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving host containers")
	}

	for cLabel, c := range cMap {
		retMatrix[cLabel] = map[cgroups.ControlType]uint64{}
		for _, control := range cgroupControls {
			ctrlPath, err := genContainerCgroupPath(control, c.ID)
			if err != nil {
				return nil, err
			}

			cgCtrl, err := cgroups.CgroupControlFactory(control, ctrlPath)
			if err != nil {
				return nil, err
			}

			stat, err := cgCtrl.Stats()
			if err != nil {
				return nil, err
			}

			retMatrix[cLabel][control] = stat
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

func genContainerCgroupPath(ctype cgroups.ControlType, id string) (string, error) {
	cgroup2, err := cgroups.IsCgroup2UnifiedMode()
	if err != nil {
		return "", errors.Wrapf(err, "determing cgroup version")
	}

	if err != nil {
		return "", errors.Wrapf(err, "retrieving cgroup root path")
	}
	path, err := filepath.Abs("/sys/fs/cgroup")

	if !cgroup2 {
		path = filepath.Join(path, ctype.String())
	}

	uid := os.Geteuid()
	if uid != 0 {
		if !cgroup2 {
			return "", fmt.Errorf("rootless cgroups require Cgroups V2")
		}
		path = filepath.Join(path, fmt.Sprintf("user.slice/user-%d.slice/user@%d.service/user.slice", uid, uid))
	}

	path = filepath.Join(path, fmt.Sprintf(containerCgroupTemplate, id))
	return path, nil
}

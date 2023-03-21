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
	"strings"

	"github.com/pkg/errors"
	"github.com/infrawatch/collectd-libpod-stats/pkg/cgroups"
	"github.com/infrawatch/collectd-libpod-stats/pkg/containers"
)

const (
	//container path relative to user
	relativeContainersPath string = "containers/storage/overlay-containers/containers.json"
	relativeVolatileContainersPath string = "containers/storage/overlay-containers/volatile-containers.json"
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
			ctrlPath, err := genContainerCgroupPath(control, c.ID, cgroup2, uid, c.Names[0])
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
	volatileContainersPath := filepath.Join("/var/lib", relativeVolatileContainersPath)

	if uid != 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving user home directory")
		}

		containersPath = filepath.Join(home, ".local/share", relativeContainersPath)
		volatileContainersPath = filepath.Join(home, ".local/share", relativeVolatileContainersPath)
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

	/*
		volatile-containers.json holds the records for the containers which have been created
		with "--rm" flag which will destroy the container and the overlay mount point when
		the container completes. Existance of this file must be checked before reading it
		because it wouldn't exist if no containers use "--rm" flag.
		More info here https://www.redhat.com/sysadmin/container-volatile-overlay-mounts
	*/
	if _, err := os.Stat(volatileContainersPath); err == nil {
		volatileContainersJSON, err := ioutil.ReadFile(volatileContainersPath)
		if err != nil {
			return nil, errors.Wrap(err, "reading libpod volatile container file")
		}

		volatileContainerList, err := containers.NewListFromJSON(volatileContainersJSON)
		if err != nil {
			return nil, errors.Wrap(err, "loading volatile container json")
		}

		for _, c := range volatileContainerList {
			for _, name := range c.Names {
				containerMap[name] = c
			}
		}
	}

	return containerMap, nil
}

func genContainerCgroupPath(ctype cgroups.ControlType, id string, cgroup2 bool, uid int, container_name string) (string, error) {

	path, err := filepath.Abs("/sys/fs/cgroup")
	if err != nil {
		return "", errors.Wrapf(err, "retrieving cgroup root path")
	}

	if cgroup2 {
		if uid != 0 {
			path = filepath.Join(path, fmt.Sprintf("user.slice/user-%d.slice/user@%d.service/user.slice/libpod-%s.scope", uid, uid, id))
		} else {
			if strings.HasPrefix(container_name, "ceph-") {
				if _, err := os.Stat("/run/ceph"); err == nil {
					/*
						A ceph cgroup path "/sys/fs/cgroup/system.slice/system-ceph\\x2d332cfe0d\\x2dcd42\\x2d5667\\x2daa09\\x2dca57970f68cc.slice/"
						is equivalent to "/sys/fs/cgroup/system.slice/system-ceph<FSID>.slice/".
						To reach this cgroup path, ceph fsid (inside /run/ceph) must be known.
					*/
					contents, _ := os.ReadDir("/run/ceph")
					if len(contents) < 1 {
						return "", errors.New("fsid directory missing")
					}
					ceph_fsid := contents[0].Name()

					// systemd escaping algorithm replaces "-" with C-style "\x2d"
					ceph_fsid_escape := strings.ReplaceAll(ceph_fsid, "-", "\\x2d")

					// create absolute path to the root directory cgroup
					path = filepath.Join(path, fmt.Sprintf("/system.slice/system-ceph\\x2d%s.slice/", ceph_fsid_escape))

					// get hostname without the domain name
					node_hostname, _ := os.Hostname()
					node_hostname = strings.Split(node_hostname, ".")[0]

					// from `container_name` get the name of the ceph service for which cgroup path is to be found
					ceph_service_name := strings.Split(container_name, node_hostname)[0]
					ceph_service_name = strings.Trim(ceph_service_name, "-")
					ceph_service_name = strings.SplitAfter(ceph_service_name, fmt.Sprintf("ceph-%s-", ceph_fsid))[1]

					/*
						Find out the directory starting with the prefix "ceph-<FSID>-" corresponding to `ceph_service_name`
						which contains the stats for that service.
					*/
					ceph_service_cgroups, _ := os.ReadDir(path)
					for _, file := range ceph_service_cgroups {
						if file.IsDir() && strings.HasPrefix(file.Name(), fmt.Sprintf("ceph-%s", ceph_fsid)) && strings.Contains(file.Name(), ceph_service_name) {
							path = filepath.Join(path, file.Name())
							break
						}
					}
				}
			} else {
				path = filepath.Join(path, fmt.Sprintf("machine.slice/libpod-%s.scope", id))
			}
		}
	} else {
		if uid != 0 {
			return "", fmt.Errorf("rootless cgroups require Cgroups V2")
		}
		path = filepath.Join(path, fmt.Sprintf("%s/machine.slice/libpod-%s.scope", ctype.String(), id))
	}
	return path, nil
}

package cgroup

//V1 represents cgroup V1 on the host OS. V1 is non-heiarchal
//and requires root to access process attributes
type V1 struct {
	basePath string
}

func (v *V1) GetCPU() (uint64, error) {
	usage, err := GetSystemCPUUsage()
	if err != nil {
		return 0, err
	}
	return usage, nil
}

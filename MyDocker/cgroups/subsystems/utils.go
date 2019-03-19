package subsystems

import (
	"fmt"
	"strings"
	"os"
	"path"
	"bufio"
)

//通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {	//按行遍历
		// 40 29 0:34 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:21 - cgroup cgroup rw,memory
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			//rw,memory
			if opt == subsystem {
				return fields[4]	//  /sys/fs/cgroup/memory
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

//找到对应subsystem挂载的hierarchy相对路径对应的   cgroup    在虚拟文件系统中的绝对路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)	//比如 通过memory得到/sys/fs/cgroup/memory
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		//cgroupPath 相当于docker容器的id
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
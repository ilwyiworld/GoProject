package cgroups

import (
	"./subsystems"
	"github.com/sirupsen/logrus"
)

//把不同subsystem中的cgroup管理起来，并与容器建立关系
type CgroupManager struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path     string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程pid加入到这个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)	//以进程号为id
	}
	return nil
}

// 设置各个subsystem挂载中cgroup资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

//释放各个subsystem挂载中的cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		if err := subSysIns.Remove(c.Path); err != nil {	//删除cgroup文件夹
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}

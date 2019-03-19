package subsystems

//用于传递资源限制配额的结构体
type ResourceConfig struct {
	MemoryLimit string	//内存
	CpuShare    string	//CPU时间片权重
	CpuSet      string	//CPU核心数
}

//将cgroup抽象成path，原因是cgroup在hierarchy的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	//返回subsystem的名字，如cpu memory
	Name() string
	//设置某个cgroup在这个subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	//将进程添加到某个cgroup中
	Apply(path string, pid int) error
	//移除某个cgroup
	Remove(path string) error
}

//通过不同的subsystem初始化实例创建资源限制处理链数组
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)

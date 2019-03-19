package container

import (
	"os"
	"os/exec"
	"syscall"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "container.log"
	RootUrl				string = "/root"
	MntUrl				string = "/root/mnt/%s"
	WriteLayerUrl 		string = "/root/writeLayer/%s"
)

type ContainerInfo struct {
	Pid         string `json:"pid"`        //容器的init进程在宿主机上的 PID
	Id          string `json:"id"`         //容器Id
	Name        string `json:"name"`       //容器名
	Command     string `json:"command"`    //容器内init运行命令
	CreatedTime string `json:"createTime"` //创建时间
	Status      string `json:"status"`     //容器的状态
	Volume      string `json:"volume"`     //容器的数据卷
	PortMapping []string `json:"portmapping"` //端口映射
}

/*
这里是父进程，也就是当前进程执行的内容
1./proc/self/exe调用中，/proc/self值的是当前运行进程自己的环境，
	exec其实就是自己调用了自己，使用这种方式对创建出来的线程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境
4.如果用户指定了-ti参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, containerName, volume, imageName string, envSlice []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorf("get init process error %v", err)
		return nil, nil
	}

	cmd := exec.Command(initCmd, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {	//前台运行
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {	//后台运行
		//将容器进程的标准输出挂载到/var/run/mydocker/容器名/container.log中
		//生成容器对应目录的container.log文件
		dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
			return nil, nil
		}
		stdLogFilePath := dirURL + ContainerLogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		//把生成好的文件赋值给stdout，这样能把容器内的标准输出重定向到这个文件中
		cmd.Stdout = stdLogFile
	}

	//传入管道文件读取端的句柄  会外带这这个文件句柄去创建子进程
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), envSlice...)
	NewWorkSpace(volume, imageName, containerName)
	cmd.Dir = fmt.Sprintf(MntUrl, containerName)//指定容器初始化后的工作目录 全部是/root/mnt 里面有writerLayer和容器目录
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()	//生成一个匿名管道 读，写 都是文件类型
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

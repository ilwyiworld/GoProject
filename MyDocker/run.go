package main

import (
	"./container"
	log "github.com/sirupsen/logrus"
	"os"
	"./cgroups/subsystems"
	"./cgroups"
	"strings"
	"strconv"
	"fmt"
	"time"
	"encoding/json"
	"math/rand"
	"./network"
)
/*
这里的Start方法是真正开始调用前面创建好的command的调用，它首先会clone出来一个namespace隔离的进程，
然后在子进程中，调用/proc/self/exe，也就是调用自己，发送init参数，调用我们写的init方法，去初始化容器的一些资源
 */
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, containerName, volume, imageName string,
			envSlice []string, nw string, portmapping []string) {
	containerID := randStringBytes(10)		//生成10位数容器id
	//如果不指定容器名，那么就以容器id作为容器名
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, containerName, volume, imageName, envSlice)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName, containerID, volume)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}
	//use mydocker-cgroup as cgroup name
	//创建cgroup manager，并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager :=cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	//设置资源限制
	cgroupManager.Set(res)
	//将容器进程加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)
	//对容器设置完限制之后，初始化容器
	if nw != "" {
		// config container network
		network.Init()
		containerInfo := &container.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portmapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
	}

	//发送用户命令
	sendInitCommand(comArray, writePipe)

	if tty {
		parent.Wait()		//用于父进程等待子进程结束
		//下面是容器退出后的操作
		deleteContainerInfo(containerName)
		container.DeleteWorkSpace(volume, containerName)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

//记录容器信息 将容器信息持久化到磁盘的/var/run/mydocker/容器名/config.json中
func recordContainerInfo(containerPID int, commandArray []string, containerName, id, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		Volume:      volume,
	}

	//将容器信息的对象json序列化为字符串
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	//存储容器的路径
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	//创建配置文件----config.json
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	//将json化之后的数据写入到文件中
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerName string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

//ID生成器
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


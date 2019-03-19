package main

import (
	log "github.com/sirupsen/logrus"
	"syscall"
	"strconv"
	"./container"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
)

func stopContainer(containerName string) {
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Get contaienr pid by name %s error %v", containerName, err)
		return
	}
	//将string类型的pid转化为int类型
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	//系统调用kill可以发送信号给进程，通过传递syscall.SIGTERM信号杀掉容器主进程
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerName, err)
		return
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	//修改容器状态，PID置为空
	containerInfo.Status = container.STOP
	containerInfo.Pid = " "
	//将修改后的信息序列化成json字符串
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerName, err)
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	//重新写入新的数据覆盖原来的信息
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Errorf("Write file %s error", configFilePath, err)
	}
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Errorf("GetContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}

func removeContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	//只删除处于停止状态的容器
	if containerInfo.Status != container.STOP {
		log.Errorf("Couldn't remove running container")
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove file %s error %v", dirURL, err)
		return
	}
	container.DeleteWorkSpace(containerInfo.Volume, containerName)
}

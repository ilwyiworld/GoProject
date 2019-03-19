package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"./container"
	"io/ioutil"
	"encoding/json"
	"strings"
	"os/exec"
	"os"
	_ "./nsenter"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func ExecContainer(containerName string, comArray []string) {
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Exec container getContainerPidByName %s error %v", containerName, err)
		return
	}

	//把命令以空格为分隔符拼接成一个字符串，便于传递
	cmdStr := strings.Join(comArray, " ")
	log.Infof("container pid %s", pid)
	log.Infof("command %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)
	//获取对应PID的环境变量，其实也是容器的环境变量
	containerEnvs := getEnvsByPid(pid)
	//将宿主机的环境变量的话容器的环境变量都防止在exec进程内
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerName, err)
	}
}

func GetContainerPidByName(containerName string) (string, error) {
	//从config.json中获得
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}

func getEnvsByPid(pid string) []string {
	//进程存放环境变量的位置 /proc/<PID>/environ
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}

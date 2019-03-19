package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"./container"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

func ListContainers() {
	// /var/run/mydocker//
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1] // /var/run/mydocker/
	//读取文件夹下的所有文件
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		log.Errorf("Read dir %s error %v", dirURL, err)
		return
	}

	var containers []*container.ContainerInfo
	//遍历该文件夹下的所有文件
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		//根据容器配置文件获取对应的信息，然后转换成容器信息的对象
		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			log.Errorf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	//在控制台打印出容器信息
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	//控制台输出的信息列
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	//刷新标准输出流缓存区，将容器列表打印出来
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
	containerName := file.Name()
	//根据文件名生成文件绝对路径
	configFileDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFileDir = configFileDir + container.ConfigName
	//读取config.json文件中的容器信息
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("Read file %s error %v", configFileDir, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	//将json文件信息反序列化成容器信息对象
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}

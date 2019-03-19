package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"./container"
	"os"
	"io/ioutil"
)

func logContainer(containerName string) {
	//找到对应文件夹的位置
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	logFileLocation := dirURL + container.ContainerLogFile
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	//将读出来的内容输入到标准输出，也就是控制台上
	fmt.Fprint(os.Stdout, string(content))
}

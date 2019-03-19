package main

import (
	"github.com/urfave/cli"
	"os"
	log "github.com/sirupsen/logrus"
)

/*
1.cgroup hierarchy中的节点，用于管理进程和subsystem的控制关系
2.subsystem作用于hierarchy上的cgroup节点，并控制节点进程的资源占用
3.hierarchy将cgroup通过树状结构串起来，并通过虚拟文件系统的方式暴露给用户
 */

const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {

	app:=cli.NewApp()
	app.Name="myDocker"
	app.Usage=usage

	app.Commands=[]cli.Command{
		initCommand,
		runCommand,
		listCommand,
		logCommand,
		execCommand,
		stopCommand,
		removeCommand,
		commitCommand,
		networkCommand,
	}

	//初始化logrus日志配置
	app.Before=func(context *cli.Context) error{
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})

		//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
		log.SetOutput(os.Stdout)
		return nil
	}

	if err:=app.Run(os.Args);err!=nil{
		log.Fatal(err)
	}

}

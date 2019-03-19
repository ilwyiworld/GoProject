package main

import(
	"os/exec"
	"syscall"
	"os"
	"log"
)

func main(){
	cmd:=exec.Command("sh")			//用来指定被fork出来的新进程内的初始命令
	cmd.SysProcAttr=&syscall.SysProcAttr{
		Cloneflags : syscall.CLONE_NEWUTS,		//创建一个UTS Namespace
	}
	cmd.Stdin=os.Stdin
	cmd.Stdout=os.Stdout
	cmd.Stderr=os.Stderr
	
	if err:=cmd.Run(); err!=nil{
		log.Fatal(err)
	}
}
package main

import (
	"os/exec"
	"os"
	"syscall"
)

func main() {
	// Go requires an absolute path to the binary we want to execute,
	// so we’ll use exec.LookPath to find it (probably /bin/ls).
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	// Exec requires arguments in slice form (as apposed to one big string).
	// We’ll give ls a few common arguments.
	// Note that the first argument should be the program name.
	args := []string{"ls", "-a", "-l", "-h"}

	// Exec also needs a set of environment variables to use.
	// Here we just provide our current environment.
	env := os.Environ()

	// the execution of our process will end here and be replaced by the /bin/ls -a -l -h process
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}


/*
$ go run ExecingProcesses.go
total 16
drwxr-xr-x  4 mark 136B Oct 3 16:29 .
drwxr-xr-x 91 mark 3.0K Oct 3 12:50 ..
-rw-r--r--  1 mark 1.3K Oct 3 16:28 execing-processes.go*/

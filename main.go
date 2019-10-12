package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main(){
	dir, _ := os.Getwd()
	if checkFileExists(dir + "/package-lock.json") && checkFileExists(dir + "/package.json"){
		npmPath, err := exec.LookPath("npm")
		if err != nil {
			log.Fatal(err)
		}
		completed := false
		for i:= 0; i < 5; i++ {
			fmt.Printf("Trying to npm ci in %s... This cmd will timeout in 3 minutes\n", dir)
			cmd := exec.Command(npmPath, "ci")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			done := make(chan error, 1)

			go func() {
				done <- cmd.Wait()
			}()

			select {
				case <- time.After(3 * time.Minute):
					if err := cmd.Process.Kill(); err != nil {
						log.Fatal("Failed to kill process: ", err)
					}
				case err := <- done:
					if err != nil{
						waitErr := err.Error()
						log.Println("Failed to execute npm ci. Err:", waitErr)
					} else {
						completed = true
					}

			}
			if completed {
				log.Println("Successfully executed npm ci")
				return
			}
		}
		log.Fatal("Failed to execute npm ci for 5 times")
	} else {
		log.Fatal("package.json or package-lock.json not found in directory: ", dir)
	}
}

func checkFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
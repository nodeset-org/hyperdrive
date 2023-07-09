package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Hyperdrive initializing...");
	rpServiceStatus := exec.Command("bash", "rocketpool service status");
	out,err := rpServiceStatus.Output();
	if err != nil {
		fmt.Println(err);
	}
	fmt.Println(string(out))
}
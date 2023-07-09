package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Hyperdrive initializing...");
	rpServiceStatus := exec.Command("rocketpool", "service status");
	out,err := rpServiceStatus.Output();
	if err != nil {
		fmt.Println("Error checking rp service status", err);
	}
	fmt.Println(string(out))
}
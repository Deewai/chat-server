package projectx

import (
	"flag"
	"fmt"
)

func main() {
	isCLI := flag.Bool("cli", false, "Run it in CLI")
	flag.Parse()

	if *isCLI {
		fmt.Println("it works")
		return
	}

	helpInfo()
}

func helpInfo() {
	fmt.Println("Usage: projectx <flag>")
	fmt.Println("\nAvailable flags are:")
	fmt.Println("cli        Start the Cli")
}

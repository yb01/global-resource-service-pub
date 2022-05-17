package main

import (
	"fmt"
	"os"

	"global-resource-service/resource-management/cmds/service-api/app"
)

func main() {

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

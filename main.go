package main

import (
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/cmd/eleutherios"
)

var version = "dev"

func main() {
	eleutherios.SetVersion(version)
	if err := eleutherios.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

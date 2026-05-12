package main

import (
	"fmt"
	"os"

	"github.com/BataevDaniil/eleutherios/cmd/eleutherios"
)

func main() {
	if err := eleutherios.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

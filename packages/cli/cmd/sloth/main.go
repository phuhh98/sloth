package main

import (
	"os"

	"github.com/phuhh98/sloth/packages/cli/internal/app"
)

func main() {
	root := app.NewRootCommand()
	if err := root.Execute(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

package main

import (
	"github.com/crispyarty/novelparser/cmd"
	"github.com/crispyarty/novelparser/internal/config"
)

func main() {
	defer config.Save()

	cmd.Execute()
}

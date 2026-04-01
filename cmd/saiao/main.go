package main

import (
	"flag"
	"fmt"
	"os"

	"saiao/internal/app"
)

func main() {
	configPath := flag.String("config", "/etc/saiao/config.yaml", "path to SAIAO config file")
	flag.Parse()

	if err := app.Run(app.Options{
		ConfigPath: *configPath,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

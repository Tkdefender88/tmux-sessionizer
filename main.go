package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"gitlab.com/Tkdefender88/tmux-sessionizer/cmd"
	"gitlab.com/Tkdefender88/tmux-sessionizer/config"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "error setting up config: %v\n", err)
		os.Exit(2)
	}

	if err := fang.Execute(context.Background(), cmd.RootCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

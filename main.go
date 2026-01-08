package main

import (
	"fmt"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/spf13/viper"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app"
	"gitlab.com/Tkdefender88/tmux-sessionizer/config"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "error setting up config: %v\n", err)
		os.Exit(2)
	}

	inputChan := make(chan string)
	outputChan := make(chan string)

	cfg := viper.GetViper()
	paths, err := app.FindDirs(cfg.GetStringSlice(config.TS_SEARCH_PATHS), cfg.GetInt(config.TS_MAX_SEARCH_DEPTH))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding dirs: %v\n", err)
	}

	go func() {
		for _, p := range paths {
			inputChan <- p
		}
	}()

	go func() {
		for s := range outputChan {
			fmt.Println("Got: " + s)
		}
	}()

	exit := func(code int, err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(code)
	}

	options, err := fzf.ParseOptions(
		true,
		[]string{"--multi", "--margin=20%", "--border", "--header=Tmux Sessionizer"},
	)
	if err != nil {
		exit(fzf.ExitError, err)
	}

	options.Input = inputChan
	options.Output = outputChan

	code, err := fzf.Run(options)
	exit(code, err)
}

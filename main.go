package main

import (
	"fmt"
	"os"
	"sync"

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
	paths, err := app.FindSessionTargets(
		cfg.GetStringSlice(config.TS_SEARCH_PATHS),
		cfg.GetInt(config.TS_MAX_SEARCH_DEPTH),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding session targets: %v\n", err)
	}

	go func() {
		for _, p := range paths {
			inputChan <- p
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Go(func() {
		for s := range outputChan {
			err := openTmuxSession(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening tmux session: %v\n", err)
			}
		}
	})

	exit := func(code int, err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(code)
	}

	code, err := launchFzf(inputChan, outputChan)
	wg.Wait()
	exit(code, err)
}

func launchFzf(input chan string, output chan string) (int, error) {
	options, err := fzf.ParseOptions(
		true,
		[]string{"--multi", "--margin=20%", "--border", "--header=Tmux Sessionizer"},
	)
	if err != nil {
		return fzf.ExitError, err
	}

	options.Input = input
	options.Output = output

	return fzf.Run(options)
}

func openTmuxSession(target string) error {
	return app.NewTmux().OpenTmuxSession(target)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	fzf "github.com/junegunn/fzf/src"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app"
	"gitlab.com/Tkdefender88/tmux-sessionizer/config"
	"gitlab.com/Tkdefender88/tmux-sessionizer/logging"
)

func init() {
	RootCmd.PersistentFlags().Bool("debug", false, "print debug logs to the terminal")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
}

var RootCmd = &cobra.Command{
	Use:   "tmux-sessionizer",
	Short: "tmux-sessionizer is a session manager for tmux, the terminal multiplexer",
	Long: `tmux-sessionizer is a session manager for tmux
	it allows you to select a directory and have a tmux session be opened to that directory.
	if a tmux session already exists for that directory you get switched to that existing tmux session
	ideal if you bind it to a keyboard shortcut in your bash/zsh/fish/whatever shell is cool now

	* configuration can be found in ~/.config/tmux-sessionizer/config.yaml
	`,
	Example: `  # Launch the interactive fuzzy finder
  tmux-sessionizer

  # Bind to a keyboard shortcut (bash/zsh)
  bind -x '"\C-f":"tmux-sessionizer"'`,
	RunE: rootCmd,
	Args: cobra.NoArgs,
}

func rootCmd(cmd *cobra.Command, args []string) error {
	inputChan := make(chan string)

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logFilePath := filepath.Join(homedir, ".local", "state", "tmux-sessionizer.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	logging.SetupLogging(logFile)

	sessionizer := app.NewSessionManager()

	go func() {
		cfg := viper.GetViper()
		search_paths := cfg.GetStringSlice(config.TS_SEARCH_PATHS)
		extra_paths := cfg.GetStringSlice(config.TS_EXTRA_SEARCH_PATHS)
		search_paths = append(search_paths, extra_paths...)
		paths, err := sessionizer.FindSessionTargets(search_paths, cfg.GetInt(config.TS_MAX_SEARCH_DEPTH))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error finding session targets: %v\n", err)
		}
		for _, p := range paths {
			inputChan <- p
		}
		close(inputChan)
	}()

	wg := new(sync.WaitGroup)
	outputChan := make(chan string)
	wg.Go(func() {
		for s := range outputChan {
			err := sessionizer.OpenSession(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening tmux session: %v\n", err)
			}
		}
	})

	_, err = launchFzf(inputChan, outputChan)
	close(outputChan)
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func launchFzf(input chan string, output chan string) (int, error) {
	options, err := fzf.ParseOptions(
		true,
		[]string{
			"--multi",
			"--margin=20%",
			"--border",
			"--header=Tmux Sessionizer",
		},
	)
	if err != nil {
		return fzf.ExitError, err
	}

	options.Input = input
	options.Output = output

	return fzf.Run(options)
}

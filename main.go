package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app"
)

type model int

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/tmux-sessionizer")

	viper.SetDefault("ts_search_paths", []string{"~/workspace"})
	viper.SetDefault("ts_extra_search_paths", []string{"~/dotfiles"})
	viper.SetDefault("ts_max_search_depth", 2)

	err := viper.ReadInConfig()
	if err != nil {
		var fileLookupError viper.ConfigFileNotFoundError
		if !errors.As(err, &fileLookupError) {
			fmt.Fprintf(os.Stderr, "encountered an error with configuration: %v\n", err)
			os.Exit(2)
		}

		fmt.Fprintf(os.Stderr, "No config found, writing file to config directory\n")
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "encountered an error writing configuration: %v\n", err)
			os.Exit(2)
		}
	}

	p := tea.NewProgram(app.Model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "encountered an error: %v\n", err)
		os.Exit(1)
	}
}

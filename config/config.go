package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	TS_SEARCH_PATHS       = "ts_search_paths"
	TS_EXTRA_SEARCH_PATHS = "ts_extra_search_paths"
	TS_MAX_SEARCH_DEPTH   = "ts_max_search_depth"
)

func SetupConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/tmux-sessionizer")

	viper.SetDefault(TS_SEARCH_PATHS, []string{"~/workspace"})
	viper.SetDefault(TS_EXTRA_SEARCH_PATHS, []string{"~/dotfiles"})
	viper.SetDefault(TS_MAX_SEARCH_DEPTH, 2)

	err := viper.ReadInConfig()
	if err != nil {
		var fileLookupError viper.ConfigFileNotFoundError
		if !errors.As(err, &fileLookupError) {
			return fmt.Errorf("error reading configuration: %v", err)
		}

		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("error writing configuration: %v", err)
		}
	}

	return nil
}

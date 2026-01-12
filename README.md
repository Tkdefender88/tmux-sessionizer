# Tmux Sessionizer

This is a session manager for [tmux](https://github.com/tmux/tmux/wiki) based on the project of the same name by 
[ThePrimeagen](https://www.youtube.com/channel/UC8ENHE5xdFSwx71u3fDH5Xw). His was a bash script that he built mostly
just for his workflow. This one is built in Go and is tailored for my workflows. I wanted to build this as a learning
project to get more familiar with tmux and understanding how to use it better. This isn't the first rewrite of this
type of project there is another pretty popular one written in rust but this one is mine and I had fun making it. :)

## Quick start

To install it you can run
```sh
go install gitlab.com/Tkdefender88/tmux-sessionizer

```

## Usage

Run `tmux-sessionizer` from your terminal. If it doesn't find a configuration file in ~/.config/tmux-sessionizer it will
create a config.yaml file there with some defaults. Select the directories that you want to have it search through with
the `ts_search_paths` parameter. I have it search my ~/workspace but you can have it search whatever list of directories
you want. You can customize the search depth (how many subdirectories it will traverse from the search paths given) I
find 2 to be a sensible default to find projects without going too deep into the projects folders. You can use
`ts_extra_search_paths` to search extra directories and you can specify custom search depths for these. For example you
can specify `~/dotfiles:1` to add your dotfiles directory to the search but have it only traverse 1 directory deep.

In your home directory if you have a ~/.tmux-sessionizer script this script will be used when starting new sessions. For
example I use it to open a second tmux window and autmoatically open my editor, neovim in the first tmux window for any
new tmux session I start. You can also include a .tmux-sessionizer script in your project directories and tmux-sessionizer
will run them when it finds them. For example, if there is a python project with a .tmux-sessionizer file in the top directory
that starts the python virtual env for the project. When tmux-sessionizer is started on that project it will autmoatically
run that script for that project and start you in your python env. There is a higher precedence placed on a local directory's
.tmux-sessionizer script rather than the one in the home directory. So if the one in the local directory is executed the
.tmux-sessionizer script in the home directory will not be executed.

I recommend setting up bash/zsh aliases or keyboard shortcuts to execute it faster. I have it bound to alt+f on my system.

## Contributing

Open to contribution, however this isn't my primary focus so I may be slow to respond. You're welcome to fork it and
customize it for yourself if you want.

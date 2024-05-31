package repl

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/chzyer/readline"
)

func getHistoryFile() (string, error) {
	var historyFile string
	user, err := user.Current()
	if err != nil {
		return historyFile, err
	}
	historyFile = filepath.Join(user.HomeDir, ".yusra_history")
	return historyFile, nil
}

func GetReader() *readline.Instance {

	historyFile, err := getHistoryFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	completer := readline.NewPrefixCompleter(
		readline.PcItem(HELP),
		readline.PcItem(HELP_SHORT),
		readline.PcItem(EXIT),
		readline.PcItem(EXIT_SHORT),
		readline.PcItem(VERSION),
		readline.PcItem(VERSION_SHORT),
		readline.PcItem(LOAD),
		readline.PcItem(LOAD_SHORT),
		readline.PcItem(LIST_LOADED_FILES),
		readline.PcItem(LIST_LOADED_FILES_SHORT),
		readline.PcItem(SHOW_FILE_DATA_INFO),
	)
	rl, err := readline.NewEx(&readline.Config{
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		Prompt:            "yusra>",
		HistoryFile:       historyFile,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("can't write to history file: %s", err))
	}

	return rl
}

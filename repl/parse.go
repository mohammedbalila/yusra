package repl

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mohammedbalila/yusra/storage"
	"github.com/xwb1989/sqlparser"
)

func isSQL(input string) bool {
	_, err := sqlparser.Parse(input)
	return err == nil
}

func ParseInput(input string) (*string, error) {
	sqlStmt := isSQL(input)
	if sqlStmt {
		result := processSQLStmt(input)
		return &result, nil
	}

	cmd, err := processCommand(input)
	return cmd, err
}

// Maybe build an actual AST to parse the incoming input instead of just stings
func processCommand(cmdLine string) (*string, error) {
	text := strings.TrimSpace(strings.Replace(cmdLine, "\n", "", -1))
	cmdWithPossibleArgs := strings.Split(text, " ")
	cmd := cmdWithPossibleArgs[0]
	switch cmd {
	case HELP, HELP_SHORT:
		return help()
	case VERSION, VERSION_SHORT:
		return version()
	case LIST_LOADED_FILES, LIST_LOADED_FILES_SHORT:
		return listLoadedFiles()
	case LOAD, LOAD_SHORT:
		return loadJsonFile(cmdWithPossibleArgs)
	case SHOW_FILE_DATA_INFO:
		return showDatasetInfo(cmdWithPossibleArgs)
	case EXIT, EXIT_SHORT:
		return exit()
	case "":
		return nil, nil
	default:
		err := fmt.Errorf("unrecognized command: %s", cmdLine)
		return nil, err

	}
}

func loadJsonFile(args []string) (*string, error) {
	if len(args) != 2 {
		return nil, errors.New("invalid syntax: use load filename.json")
	}

	filename := args[1]
	fmt.Printf("loading %s...\n", filename)
	err := storage.LoadNewJsonFile(filename)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func help() (*string, error) {
	help_msg := "yusra version 0.0.1\n"
	help_msg += "type \"help\" for usage hints.\n"
	help_msg += "type \"exit\" to exit\n"
	help_msg += "type \"version\" to show version\n"
	help_msg += "type \"load\" file.json to load a new json file\n"
	help_msg += "type \"files\" to list loaded sets\n"
	help_msg += "type \"info\" dataset_name to get info about a loaded dataset\n"
	return &help_msg, nil
}

func listLoadedFiles() (*string, error) {
	err := storage.GetLoadedDatasets()
	if err != nil {
		return nil, err
	}
	return nil, err
}

func showDatasetInfo(args []string) (*string, error) {
	if len(args) != 2 {
		return nil, errors.New("invalid syntax: use info dataset_name")
	}
	err := storage.GetTableStats(args[1])
	if err != nil {
		return nil, err
	}
	return nil, err
}

func exit() (*string, error) {
	os.Exit(0)
	return nil, nil
}

func version() (*string, error) {
	version := "yusra version 0.0.1"
	return &version, nil
}

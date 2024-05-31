package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mohammedbalila/yusra/repl"
	"github.com/mohammedbalila/yusra/storage"
	"github.com/spf13/cobra"
)

type Result []interface{}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads a new json file and creates a sqlite table for it",
	Long: `
	Loads a new json file and creates a sqlite table for it
	usage: yusra load file.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "invalid use of the command")
			os.Exit(1)
		}

		// read file
		filename := args[0]
		err := storage.LoadNewJsonFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("failed to load json file: %s", err))
			os.Exit(1)
		}

		rl := repl.GetReader()
		defer rl.Close()

		welcome_msg := fmt.Sprintf("yusra version 0.0.1 %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(welcome_msg)
		fmt.Println("Enter \"\\help\" for usage hints.")
		for {

			// Read the input from the user
			input, err := rl.Readline()
			if err != nil { // Handle EOF or Ctrl-D
				os.Exit(0)
			}

			result, err := repl.ParseInput(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				if result != nil {
					fmt.Fprintln(os.Stdout, *result)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

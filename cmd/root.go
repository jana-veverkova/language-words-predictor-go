package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "language-words-predictor",
	Short: "Description",
}

func Execute() {
	rootCmd.AddCommand(trainTestSplitCmd)
	rootCmd.AddCommand(populationSampleCmd)
	rootCmd.AddCommand(frequencyDictionaryCmd)
	rootCmd.AddCommand(trainCmd)
	rootCmd.AddCommand(runTestCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
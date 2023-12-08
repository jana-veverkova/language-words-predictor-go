package cmd

import (
	"github.com/jana-veverkova/language-words-predictor-go/pkg/model"
	"github.com/spf13/cobra"
)

var runTestCmd = &cobra.Command{
	Use:   "run-test",
	Short: "Runs new test in terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		model.RunTest()
		return nil
	},
}
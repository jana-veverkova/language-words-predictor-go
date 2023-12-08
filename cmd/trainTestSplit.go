package cmd

import (
	"github.com/jana-veverkova/language-words-predictor-go/pkg/traintestsplit"
	"github.com/spf13/cobra"
)

var trainTestSplitCmd = &cobra.Command{
	Use:   "train-test-split",
	Short: "Splits population sample data into train and test set.",
	RunE: func(cmd *cobra.Command, args []string) error {
		traintestsplit.Split("data/processed/populationSamples", "data/processed")
		return nil
	},
}
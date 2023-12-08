package cmd

import (
	"github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
	"github.com/spf13/cobra"
)

var populationSampleCmd = &cobra.Command{
	Use:   "create-population-sample",
	Short: "Creates population samples files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		populationsample.CreateSourceFiles("data/processed/populationSamples")
		return nil
	},
}
package cmd

import (
	"github.com/jana-veverkova/language-words-predictor-go/pkg/model"
	"github.com/spf13/cobra"
)

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Runs predictions on test set using different setting parameters to compare accuracy of result.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := model.Train()
		if err != nil {
			return err
		}
		return nil
	},
}
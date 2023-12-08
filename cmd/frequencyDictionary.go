package cmd

import (
	"github.com/jana-veverkova/language-words-predictor-go/pkg/frequencydictionary"
	"github.com/spf13/cobra"
)

var frequencyDictionaryCmd = &cobra.Command{
	Use:   "create-frequency-dictionary",
	Short: "Creates frequency dictionary from text files stored in data/original/'language' directory.",
	ValidArgs: []string{"english"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		frequencydictionary.CreateFile("data/original/"+args[0], "data/frequencyDictionaries/"+args[0]+".csv")
		return nil
	},
}
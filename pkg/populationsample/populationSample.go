package populationsample

import (
	"sync"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/filehandler/csvx"
	"github.com/pkg/errors"
)

type PopulationSample struct {
	Table   [][]bool
	RowSums []int
	ColSums []int
}

func GetAllSamples() (*PopulationSample, error) {
	// creates population samples struct

	targetDir := "data/processed/populationSamples/"
	files := []string{"populationSample0.csv", "populationSample1.csv", "populationSample2.csv", "populationSample3.csv"}

	table := make([][]bool, 0)

	for _, file := range files {
		err := appendToTable(&table, targetDir+file)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	rowSums, colSums := computeRowColSums(&table)

	return &PopulationSample{Table: table, RowSums: rowSums, ColSums: colSums}, nil
}

func GetTrainSample() (*PopulationSample, error) {
	sourceFile := "data/processed/train.csv"

	table := make([][]bool, 0)

	err := appendToTable(&table, sourceFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rowSums, colSums := computeRowColSums(&table)

	return &PopulationSample{Table: table, RowSums: rowSums, ColSums: colSums}, nil
}

func GetTestSample() (*PopulationSample, error) {
	sourceFile := "data/processed/test.csv"

	table := make([][]bool, 0)

	err := appendToTable(&table, sourceFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rowSums, colSums := computeRowColSums(&table)

	return &PopulationSample{Table: table, RowSums: rowSums, ColSums: colSums}, nil
}

func appendToTable(table *[][]bool, fileName string) error {
	var bools []bool

	var wgReader sync.WaitGroup
	chRecords, chError := csvx.Read(fileName, &wgReader)

	for {
		select {
		case record, ok := <-chRecords:
			if !ok {
				return nil
			}

			bools = []bool{}

			for _, s := range record {
				var b bool
				if s == "1" {
					b = true
				} else {
					b = false
				}
				bools = append(bools, b)
			}

			*table = append(*table, bools)

		case err := <-chError:
			return errors.WithStack(err)
		}
	}
}

func computeRowColSums(table *[][]bool) ([]int, []int) {
	rowSums := make([]int, len(*table))
	colSums := make([]int, len((*table)[0])+1)

	for rowIx, row := range *table {
		for colIx, isWordKnown := range row {
			if isWordKnown {
				rowSums[rowIx]++
				colSums[colIx]++
			}
		}
	}

	return rowSums, colSums
}

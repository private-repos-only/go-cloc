package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_report_GetTotalFromCSV(t *testing.T) {

	total := GetTotalFromCSV("test-files/test.csv")

	// Assert
	assert.Equal(t, 100, total)

}

func Test_report_ParseTotalsFromCSVs(t *testing.T) {
	csvFilePaths := []string{
		"test-files/lodash.csv",
		"test-files/test.csv",
	}

	repoTotals := ParseTotalsFromCSVs(csvFilePaths)
	sum := 0
	for _, total := range repoTotals {
		sum += total.Total
	}

	// Assert
	assert.Equal(t, 23105, sum)
}

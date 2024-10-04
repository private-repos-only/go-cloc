package report

import (
	"encoding/csv"
	"go-cloc/scanner"
	"log"
	"os"
	"sort"
	"strconv"
)

type RepoTotal struct {
	RepositoryName string
	Total          int
}

// OutputCSV writes the results of the scan to a CSV file
// Returns the total number of lines of code for all files scanned
func OutputCSV(inputArr []scanner.FileScanResults, outputFilePath string) int {

	// Sort by CodeLineCount desc
	sort.Slice(inputArr, func(a, b int) bool {
		return inputArr[a].CodeLineCount > inputArr[b].CodeLineCount
	})

	// Create CSV information
	records := [][]string{
		{"filePath", "blank", "comment", "code"},
	}
	sumBlankLinesCount := 0
	sumCommentLinesCount := 0
	sumCodeLinesCount := 0
	for _, results := range inputArr {
		row := []string{results.FilePath, strconv.Itoa(results.BlankLineCount), strconv.Itoa(results.CommentsLineCount), strconv.Itoa(results.CodeLineCount)}
		records = append(records, row)
		sumBlankLinesCount += results.BlankLineCount
		sumCommentLinesCount += results.CommentsLineCount
		sumCodeLinesCount += results.CodeLineCount
	}
	// Append Total Row
	totalRow := []string{"total", strconv.Itoa(sumBlankLinesCount), strconv.Itoa(sumCommentLinesCount), strconv.Itoa(sumCodeLinesCount)}
	records = append(records, totalRow)

	// Write to csv
	writeCsv(outputFilePath, records)

	return sumCodeLinesCount
}

// TODO return true or false
func writeCsv(outputFilePath string, records [][]string) {
	// Write to csv
	f, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, row := range records {
		err = w.Write(row)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func OutputCombinedCSV(repoResults []RepoTotal, outputFilePath string) int {
	// Create CSV information
	records := [][]string{
		{"repository", "lineOfCodeCount"},
	}
	sum := 0
	for _, repoResult := range repoResults {
		row := []string{repoResult.RepositoryName, strconv.Itoa(repoResult.Total)}
		records = append(records, row)
		// keep running total
		sum += repoResult.Total
	}
	// Create total row
	totalRow := []string{"total", strconv.Itoa(sum)}
	records = append(records, totalRow)

	// Write to csv
	writeCsv(outputFilePath, records)
	return sum
}

package report

import (
	"encoding/csv"
	"go-cloc/logger"
	"go-cloc/scanner"
	"log"
	"os"
	"sort"
	"strconv"
)

func OutputCSV(inputArr []scanner.FileScanResults, outputFilePath string) {

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

type RepoTotal struct {
	RepositoryName string
	Total          int
}

func ParseTotalsFromCSVs(csvFilePaths []string) []RepoTotal {
	repoResults := []RepoTotal{}
	for _, csvFilePath := range csvFilePaths {
		// read total line
		total := GetTotalFromCSV(csvFilePath)
		// convert it to an int
		// TODO make this the repository name instead of the CSV
		repoResult := RepoTotal{
			RepositoryName: csvFilePath,
			Total:          total,
		}
		repoResults = append(repoResults, repoResult)
	}
	// Sort by total
	sort.Slice(repoResults, func(a, b int) bool {
		return repoResults[a].Total > repoResults[b].Total
	})

	return repoResults
}

func OutputCombinedCSV(repoResults []RepoTotal, outputFilePath string) int {
	// Create CSV information
	records := [][]string{
		{"filePath", "code"},
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

func GetTotalFromCSV(csvFilePath string) int {
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	var lastLine []string

	// Iterate through all the lines in the CSV
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" { // End of file, stop reading
				break
			}
			log.Fatal(err)
		}
		lastLine = record // Keep track of the last read line
	}

	// Print the last line
	if len(lastLine) > 3 {
		// Convert the string at index 3 to an integer
		value, err := strconv.Atoi(lastLine[3])
		if err != nil {
			logger.Error("Error converting index 3 to an integer:", err)
			return 0
		}
		return value
	}
	// Error
	return 0
}

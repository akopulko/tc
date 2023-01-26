package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/navossoc/bayesian"
	"golang.org/x/exp/slices"
)

// finds string in slice of strings
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// checks if string is numeric
func isNumeric(word string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(word)
}

// reads content of csv file
func readCSV(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// checks if feature is not single character and not numeric
func checkFeature(feature string) bool {
	if (len(feature) > 1) && (!isNumeric(feature)) {
		return true
	} else {
		return false
	}

}

// makes training data from CSV line entry
// returns category and slice of unique features
func parseTrainingEntry(data []string) (string, []string) {
	category := data[1]
	words := strings.Split(data[0], " ")
	var features []string
	for _, word := range words {
		if (checkFeature(word)) && (!contains(features, word)) {
			features = append(features, word)
		}
	}
	return category, features

}

// get slice of categories from training map
func getCategories(training map[string][]string) []bayesian.Class {
	var result []bayesian.Class
	for key := range training {
		result = append(result, bayesian.Class(key))
	}
	return result
}

// build training data map
// map[Category] = [feature1, feature2, ...]
func buildTrainingMap(data [][]string) map[string][]string {
	resultMap := make(map[string][]string)
	var features []string
	var category string
	for _, line := range data {
		category, features = parseTrainingEntry(line)
		_, exist := resultMap[category]
		if exist {
			resultMap[category] = append(resultMap[category], features...)
		} else {
			resultMap[category] = features
		}
	}
	return resultMap
}

// saves content to csv file
func saveCSV(content [][]string, file string) error {

	var header = []string{"DATE", "AMT", "TRN", "CATEG"}

	f, err := os.Create(file)

	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(f)

	_ = csvWriter.Write(header)

	for _, row := range content {
		_ = csvWriter.Write(row)
	}
	csvWriter.Flush()
	f.Close()

	return nil
}

// extracts features from transaction description and returns slice
// removes the following: single character words, numeric words, duplicates
func extractTransactionFeatures(transaction string) []string {

	var transFeatures []string
	features := strings.Split(transaction, " ")
	for _, feature := range features {
		if (len(feature) > 1) && (!contains(transFeatures, feature)) && (!isNumeric(feature)) {
			transFeatures = append(transFeatures, feature)
		}
	}
	return transFeatures
}

// print statement on the screen
func printStatement(statement [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "DATE\tAMT\tTRANSACTION\tCATEGORY")
	for _, line := range statement {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", line[0], line[1], line[2], line[3])
	}
	w.Flush()
}

// print classifier internal information
func printClassifierInfo(classifier bayesian.Classifier) {
	var words map[string]float64
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "CATRGORY\tFEATURE\tPROBABILITY")
	for _, class := range classifier.Classes {
		words = classifier.WordsByClass(class)
		for key, value := range words {
			fmt.Fprintf(w, "%s\t%s\t%f\n", class, key, value)
		}
	}
	w.Flush()

}

// helper function to compare two string slices of slices
func sliceOfSlicesEqual(slice1, slice2 [][]string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	var count int = 0
	for _, itemSlice1 := range slice1 {
		for _, itemSlice2 := range slice2 {
			if slices.Equal(itemSlice1, itemSlice2) {
				count++
			}
		}
	}

	return count == len(slice1)
}

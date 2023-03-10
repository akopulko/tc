package main

import (
	"reflect"
	"testing"

	"github.com/navossoc/bayesian"
	"golang.org/x/exp/slices"
)

func Test_sliceOfSlicesEqual(t *testing.T) {

	var tests = []struct {
		name     string
		slice1   [][]string
		slice2   [][]string
		expected bool
	}{
		{
			"Test two string slice of slices equal",
			[][]string{
				{"Food Market 233 Some Street", "Groceries"},
				{"Corner Shop SOMEWHERE", "Groceries"},
			},
			[][]string{
				{"Food Market 233 Some Street", "Groceries"},
				{"Corner Shop SOMEWHERE", "Groceries"},
			},
			true,
		},
		{
			"Test two string slice of slices equal un-ordered",
			[][]string{
				{"Food Market 233 Some Street", "Groceries"},
				{"Corner Shop SOMEWHERE", "Groceries"},
			},
			[][]string{
				{"Corner Shop SOMEWHERE", "Groceries"},
				{"Food Market 233 Some Street", "Groceries"},
			},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := sliceOfSlicesEqual(tc.slice1, tc.slice2)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected (%v) is not same as actual (%v)", tc.expected, actual)
			}
		})
	}

}

func Test_contains(t *testing.T) {
	var slice = []string{"one", "two", "three"}

	var search string = "two"
	actual := contains(slice, search)
	var expected bool = true
	if actual != expected {
		t.Errorf("Expected (%t) is not same as actual (%t)", expected, actual)
	}

	search = "five"
	actual = contains(slice, search)
	expected = false
	if actual != expected {
		t.Errorf("Expected (%t) is not same as actual (%t)", expected, actual)
	}

}

func Test_isNumeric(t *testing.T) {

	var tests = []struct {
		str      string
		expected bool
	}{
		{"12345", true},
		{"hello", false},
		{"hello123", false},
	}

	for _, tc := range tests {
		tname := tc.str
		t.Run(tname, func(t *testing.T) {
			actual := isNumeric(tc.str)
			if actual != tc.expected {
				t.Errorf("Expected (%t) is not same as actual (%t)", tc.expected, actual)
			}
		})
	}
}

func Test_ReadCSV(t *testing.T) {

	var f string = "../testdata/invalid.csv"
	_, err := readCSV(f)
	if err == nil {
		t.Errorf("Expected error but got (%e)", err)
	}

	expected := [][]string{
		{"Food Market 233 Some Street", "Groceries"},
		{"Corner Shop SOMEWHERE", "Groceries"},
		{"Beer & Wine World", "Alcohol"},
		{"Amazon Shopping", "Shopping"},
		{"Bottle Shop 555 Street", "Alcohol"},
		{"Amazon Marketplace", "Shopping"},
		{"eBay", "Shopping"},
	}

	f = "../sample/training.csv"
	actual, err := readCSV(f)
	if err != nil {
		t.Errorf("Expected no error but got (%e)", err)
	}

	if !sliceOfSlicesEqual(actual, expected) {
		t.Errorf("Expected (%v) is not same as actual (%v)", expected, actual)
	}
}

func Test_checkFeature(t *testing.T) {

	var tests = []struct {
		feature  string
		expected bool
	}{
		{"Q", false},
		{"12345", false},
		{"hello", true},
		{"hello123", true},
	}

	for _, tc := range tests {
		tname := tc.feature
		t.Run(tname, func(t *testing.T) {
			actual := checkFeature(tc.feature)
			if actual != tc.expected {
				t.Errorf("Expected (%t) is not same as actual (%t)", tc.expected, actual)
			}
		})
	}
}

func Test_ParseTrainingEntry(t *testing.T) {

	var tests = []struct {
		training         []string
		expectedCat      string
		expectedFeatures []string
	}{
		{[]string{"one two three", "cat1"}, "cat1", []string{"one", "two", "three"}},
		{[]string{"one two three one", "cat1"}, "cat1", []string{"one", "two", "three"}},
		{[]string{"one two three S", "cat1"}, "cat1", []string{"one", "two", "three"}},
		{[]string{"one two three 12345", "cat1"}, "cat1", []string{"one", "two", "three"}},
	}

	for _, tc := range tests {
		tname := tc.training[0]
		t.Run(tname, func(t *testing.T) {
			actual1, actual2 := parseTrainingEntry(tc.training)
			if (actual1 != tc.expectedCat) || (!reflect.DeepEqual(actual2, tc.expectedFeatures)) {
				t.Errorf(
					"Expected (%s | %v) is not same as actual (%s | %v)",
					tc.expectedCat,
					tc.expectedFeatures,
					actual1,
					actual2,
				)
			}
		})
	}
}

func Test_getCategories(t *testing.T) {
	var tests = []struct {
		name     string
		data     map[string][]string
		expected []bayesian.Class
	}{
		{
			"Test to extract correct categories",
			map[string][]string{
				"key1": {"value1", "value2"},
				"key2": {"value1", "value2"},
			},
			[]bayesian.Class{"key1", "key2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getCategories(tc.data)
			slices.Sort(actual)
			slices.Sort(tc.expected)
			if !slices.Equal(actual, tc.expected) {
				t.Errorf("Expected (%v) is not same as actual (%v)", tc.expected, actual)
			}
		})
	}
}

func Test_buildTrainingMap(t *testing.T) {
	var tests = []struct {
		name     string
		data     [][]string
		expected map[string][]string
	}{
		{
			"Test to build training map from CSV data",
			[][]string{
				{"Food Market 233", "Groceries"},
				{"Corner Shop 233", "Groceries"},
				{"Beer & Wine World", "Alcohol"},
			},
			map[string][]string{
				"Groceries": {"Food", "Market", "Corner", "Shop"},
				"Alcohol":   {"Beer", "Wine", "World"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := buildTrainingMap(tc.data)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected (%v) is not same as actual (%v)", tc.expected, actual)
			}
		})
	}

}

func Test_extractTransactionFeatures(t *testing.T) {
	var tests = []struct {
		name     string
		testData string
		expected []string
	}{
		{
			"Test to extract features from transaction",
			"food shop number 23 X",
			[]string{"food", "shop", "number"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractTransactionFeatures(tc.testData)
			if !slices.Equal(actual, tc.expected) {
				t.Errorf("Expected (%v) is not same as actual (%v)", tc.expected, actual)
			}
		})
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/navossoc/bayesian"
)

const modelFile = "model.gob"

var usage = `
Usage:
	tc train <training_file>
	tc classify <statement_file> [--out <output_file>] [--learn]
	tc show model

Options:
	-h --help     			Show this screen.
	---out <output_file> 	CSV file to save categorised transactions. By default they are printed on screen
	--learn      			LEt the model learn during transactions classification
`

type CmdLineArgs struct {
	TrainingFile  string `docopt:"<training_file>"`
	Train         bool   `docopt:"train"`
	Classify      bool   `docopt:"classify"`
	StatementFile string `docopt:"<statement_file>"`
	Output        bool   `docopt:"--out"`
	OutputFile    string `docopt:"<output_file>"`
	Learn         bool   `docopt:"--learn"`
	Show          bool   `docopt:"show"`
	Model         bool   `docopt:"model"`
}

func main() {

	// parse docopt
	args, err := docopt.ParseDoc(usage)

	if err != nil {
		log.Panic(err)
	}

	// bind cmd line arguments
	var cmdLineArgs CmdLineArgs
	err = args.Bind(&cmdLineArgs)
	if err != nil {
		panic(err)
	}

	if cmdLineArgs.Show && cmdLineArgs.Model {
		classifier, err := bayesian.NewClassifierFromFile(modelFile)
		if err != nil {
			panic(err)
		}
		printClassifierInfo(*classifier)
	}

	if cmdLineArgs.Train && len(cmdLineArgs.TrainingFile) > 0 {

		// let's train the thing by using some training data

		// read CSV with training data
		training, err := readCSV(cmdLineArgs.TrainingFile)
		if err != nil {
			panic(err)
		}

		// build Training map "categ" = ["feature1", "feature2"]
		trainingMap := buildTrainingMap(training)

		//get categories from training map
		categories := getCategories(trainingMap)

		// make classifier and train it with data
		classifier := bayesian.NewClassifier(categories...)
		for _, categ := range categories {
			classifier.Learn(trainingMap[string(categ)], categ)
		}
		err = classifier.WriteToFile(modelFile)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Trained model successfully saved to: %s\n", modelFile)
	}

	if cmdLineArgs.Classify && len(cmdLineArgs.StatementFile) > 0 {

		classifier, err := bayesian.NewClassifierFromFile(modelFile)
		if err != nil {
			panic(err)
		}

		// now let's read that statement and add categories for transactions
		// statement CSV with the following fields
		// Date, Amount, Transaction
		statement, err := readCSV(cmdLineArgs.StatementFile)
		if err != nil {
			panic(err)
		}

		var category bayesian.Class
		var features []string

		// for the output CSV we will add new field, Category
		for i, trn := range statement {
			features = extractTransactionFeatures(trn[2])
			_, likely, _ := classifier.LogScores(features)
			category = classifier.Classes[likely]
			trn = append(trn, string(category))
			statement[i] = trn
			// update model with additional learning if --learn was provided
			if cmdLineArgs.Learn {
				classifier.Learn(features, category)
			}
		}

		// save updated model if --learn was provided
		if cmdLineArgs.Learn {
			err = classifier.WriteToFile(modelFile)
			if err != nil {
				panic(err)
			}
		}

		// check if user want to save to file or print on screen (default)
		if cmdLineArgs.Output && len(cmdLineArgs.OutputFile) > 0 {
			err = saveCSV(statement, cmdLineArgs.OutputFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Classified transactions saved in: %s\n", cmdLineArgs.OutputFile)
		} else {
			printStatement(statement)
		}

	}

}

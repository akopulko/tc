# Little tool to categorise transactions in bank statements. 

## Usage
```
	tc train <training_file>
	tc classify <statement_file> [--out <output_file>] [--learn]
	tc show model

```
### Train
`tc train <training_file>`

You need to train model with some training data before it can properly categorise your transactions. `<training_file>` - see `sample` folder for file format. In general it's just a CSV file with the following format:`<transaction>,<category>`

### Show Model
`tc show model`

This will print you list of categories model learned, words associated with the category and their probabilities. The model will be stored as `model.gob` file in the same folder as the tool. 

### Classify
`tc classify <statement_file> [--out <output_file>] [--learn]`

This will classify provided transactions in statement. By default, tool will print output in `stdout` but you can specify output to CSV file with `--out` option. If you provide `--learn` option, this will add more training data to model from transactions that are classified.
`<statement_file>` - currently tool supports CSV statements int he following format: `<date>,<amount>,<transaction>`. Feel free add other statements support if you interested. 
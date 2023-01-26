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


### Sample Run
```shell

foo@test [tc]$ ./tc train training.csv
Trained model successfully saved to: model.gob

foo@test [tc]$ ls
model.gob      statement.csv  tc*            training.csv

foo@test [tc]$ ./tc classify statement.csv
DATE       |AMT       |TRANSACTION          |CATEGORY
22/10/2022 | -10.00   | Purchase Food Stall  |Groceries
12/12/2022 | -50.00   | Amazon MKT          |Shopping
31/12/2022 | -30.00   | Wine Bottle Shop    |Alcohol
01/01/2023 | -20.99   | eBay Purchase 12345 |Shopping

foo@test [tc]$ ./tc show model
CATRGORY  |FEATURE     |PROBABILITY
Groceries |Shop        |0.142857
Groceries |SOMEWHERE   |0.142857
Groceries |Food        |0.142857
Groceries |Market      |0.142857
Groceries |Some        |0.142857
Groceries |Street      |0.142857
Groceries |Corner      |0.142857
Alcohol   |Bottle      |0.166667
Alcohol   |Shop        |0.166667
Alcohol   |Street      |0.166667
Alcohol   |Beer        |0.166667
Alcohol   |Wine        |0.166667
Alcohol   |World       |0.166667
Shopping  |Amazon      |0.400000
Shopping  |Shopping    |0.200000
Shopping  |Marketplace |0.200000
Shopping  |eBay        |0.200000

```
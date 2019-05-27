# e164-csv-normalizer
Application formats various phone number strings input by users to E.164 common format. The application generates a CSV file with normalized phone numbers based on an input CSV file.

## Usage

```bash
$ go get github.com/4rakcheev/e164-csv-normalizer
$ go build
$ ./e164-csv-normalizer -i testNumbersDB.csv
```
Launch parameters

param | default | description
--- | --- | ---
`-i` | - | Input `csv file`
`-d` | y | Set to `n` for Don't `Remove duplicates` after format
`-h` | n | Set to `y` for Remove `first row as header` in the IN file
`-o` | - | Path for output normalized `csv file`
`-n` | - |  Set a National Prefix for non e164 numbers. Choose the scenario parameter `sn` for use this feature
`-sn` | - |  Set one of Scenarios for the National prefix replacement (you can use multiple scenarios like `za`):<br>`z` replace first zero to the prefix<br>`a` add the prefix to all numbers except National Prefix itself


## Example
Normalize the test database with numbers in varied formats (like user inputs).
This example skip first row as header, replace first 0 to prefix 358, append 358 prefix to other numbers without this prefix and removes duplicates
```bash
$ ./e164-csv-normalizer -i testNumbersDB.csv -n 358 -sn za -d y -h y
Processed [16] rows from file `testNumbersDB.csv`
Normalized numbers [12] (removed [2] duplicates) with wrong number [1] saved in `normalized_testNumbersDB.csv`
(!)skipped line 8 with error: "number 3.58504E+11 in large exponent format"
```

Input example:
```csv
some csv header
+358 40 727 9689
+358 400903691
+358 44 0308202
+358040 8545115
+358400290288
00250
3.58504E+11
0034600222090
040 1597474
40 1597474
040 252 7629
040 861 41 81
35840 861 41 81
0400109848
358407411963

```

Output:
```csv
358400903691
358440308202
3580408545115
358400290288
3580250
358034600222090
358401597474
358402527629
358408614181
358400109848
358407411963

```
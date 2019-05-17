# e164-csv-normalizer
Application formats various phone number strings input by users to E.164 common format. The application generates a CSV file with normalized phone numbers based on an input CSV file.

## Usage

```bash
$ go build
```
Launch parameters

param | default | description
--- | --- | ---
`-i` | - | Input `csv file`
`-d` | y | Set to "n" for Don't `Remove duplicates` after format
`-n` | - | Replace first 0 to this `National Prefix`
`-h` | n | Set to "y" for Remove `first row as header` in the IN file
`-o` | - | Path for output normalized `csv file`

## Example
Normalize the test database with numbers in varied formats (like user inputs):
```bash
$ ./e164-csv-normalizer -i testNumbersDB.csv -n 358
Processed [14] rows from file `testNumbersDB.csv`
Normalized numbers [13] (removed [1] duplicates) saved in `normalized_testNumbersDB.csv`
```

Input example:
```$xslt
+358 40 727 9689
+358 400903691
+358 44 0308202
+358040 8545115
+358400290288
00250
0034600222090
040 1597474
040 252 7629
040 861 41 81
35840 861 41 81
0400109848
358407411963
```

Output:
```$xslt
358407279689
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
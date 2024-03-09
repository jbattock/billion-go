generate measurements.txt
`git clone https://github.com/gunnarmorling/1brc.git`

`./mvnw clean verify`
`./create_measurements.sh 1000000000`
`mv ../1brc/measurements.txt measurements.txt`
expected output
`./calculate_average_baseline.sh > expected_output.txt`
`mv ../1brc/expected_output.txt expected_output.txt`

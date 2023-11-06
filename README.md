# ocm-csv-parser
Parser CLI to format csv data to be consumed by OCM services

## Build
Run `go mod tidy`

Then, run `go build`

## Usage
`ocm-csv-parser download` downloads the CSV file from Google Sheets and places in `/tmp/`. Not yet finished

`ocm-csv-parser parse` parses the CSV file in `/tmp/`, or, a specific file using a full path with the `--file / -f` flag, then deletes the temporary file in `/tmp/` if used

You can also specify the output directory / filename using `--output / -o`. Example: `ocm-csv-parser parse --output ~/Documents/my_file.configmap.yaml`

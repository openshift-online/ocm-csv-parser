# ocm-csv-parser
Parser CLI to format csv data to be consumed by OCM services

## Build
Run `go mod tidy`

Then, run `go build`

## Usage
`ocm-csv-parser download` downloads the CSV file from Google Sheets and places in `/tmp/`. Not yet finished

`ocm-csv-parser parse [command]` parses the CSV file in `/tmp/`, or, a specific file using a full path with the `--file / -f` flag, then deletes the temporary file in `/tmp/` if used

Options for `[command]` are the following:
* `ocm-csv-parser parse machinetypes`

You can also specify the output directory / filename using `--output / -o`. Example: `ocm-csv-parser parse machinetypes --output ~/Documents/my_file.configmap.yaml`

## Testing
Go into the `pkg` directory (`cd pkg`) and run `ginkgo` to run unit tests
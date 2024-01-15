# ocm-csv-parser
Parser CLI to format csv data to be consumed by OCM services

## Build
Run `go mod tidy`

Then, run `go build`

## Usage
`ocm-csv-parser download` downloads the CSV file from Google Sheets and places in `/tmp/`. Not yet finished.

(If you download the csv yourself, rename it to `cloudresources.csv` for the next step)

`ocm-csv-parser parse [command]` parses the CSV file as `/tmp/cloudresources.csv`, or, a specific file using a full path with the `--file / -f` flag, then deletes the temporary file in `/tmp/` if used

Options for `[command]` are the following:
* `ocm-csv-parser parse machinetypes`

You can also specify the output directory / filename using `--output / -o`. Example: `ocm-csv-parser parse machinetypes --output ~/Documents/my_file.configmap.yaml`

### Update to usage

If you would like, you can now specify the `--app-interface` flag with the *root directory* of your local app-interface. Doing so will automatically update the configmap within app-interface itself to save time.

Example usage:

`ocm-csv-parser parse regions --file ~/Downloads/new-configmap.csv --app-interface ../app-interface/`

**You can now parse regions**

## Testing
Go into the `pkg` directory (`cd pkg`) and run `ginkgo` to run unit tests

## Supported cloud resources

* Cloud Regions
* Machine/Instance Types

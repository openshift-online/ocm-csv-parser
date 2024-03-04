# ocm-csv-parser
Parser CLI to format csv data to be consumed by OCM services

## Build

---

Run `go mod tidy`

Then, run `go build`

## Usage

---

### Parse

`ocm-csv-parser parse [command]` parses the CSV file as `/tmp/cloudresources.csv`, or, a specific file using a full path with the `--file / -f` flag, then deletes the temporary file in `/tmp/` if used

Options for `[command]` are the following:
* `ocm-csv-parser parse machinetypes`

You can also specify the output directory / filename using `--output / -o`. Example: `ocm-csv-parser parse machinetypes --output ~/Documents/my_file.configmap.yaml`

For regions, look at the section below for an example command (including the new app-interface flag)

### Special `app-interface` flag for parse usage

---

If you would like, you can specify the `--app-interface` flag with the *root directory* of your local app-interface. Doing so will automatically update the configmap within app-interface itself to save time.

Example usage:

`ocm-csv-parser parse regions --file ~/Downloads/new-configmap.csv --app-interface ../app-interface/`

### Edit constraints

---

### Region

If you would like to edit a specific constraint in a specific env directly, please use the following command:

`ocm-csv-parser edit constraint region --app-interface ../app-interface --id me-central-1 --environment integration --ccs-only=true --enabled=true --govcloud=false`

#### Required flags:
* environment
* id
* app-interface

You may omit any of the other flags such as `--govcloud` if you do not wish to change the value

### Machine type

If you would like to edit a specific constraint in a specific env directly, please use the following command:

`ocm-csv-parser edit constraint machinetype --app-interface ../app-interface --id me-central-1 --environment integration --ccs-only=true --enabled=true`

#### Required flags:
* environment
* id
* app-interface

You may omit any of the other flags such as `--ccs-only` if you do not wish to change the value

## Testing

---

Go into the `pkg` directory (`cd pkg`) and run `ginkgo run` to run unit tests

## Supported cloud resources

---

* Cloud Regions
* Machine/Instance Types

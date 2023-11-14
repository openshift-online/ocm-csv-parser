package convert

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	. "github.com/openshift-online/ocm-csv-parser/pkg/fieldnames"
	"github.com/openshift-online/ocm-csv-parser/pkg/helper"
	"github.com/openshift-online/ocm-csv-parser/pkg/validations"
)

func CsvToResources(resourcesFile string) ([]map[string]interface{}, error) {
	finalResources := make([]map[string]interface{}, 0)
	csvFile, err := os.Open(resourcesFile)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	keys := []string{Id, GenericName, Name, CloudProviderId, CpuCores, Memory, MemoryPretty, Category,
		CategoryPretty, SizePretty, ResourceType, Active}

	// Parse resources from CSV
	for _, record := range records {
		resource := make(map[string]interface{})
		for i, key := range keys {
			resource[key] = record[i]
		}
		finalResources = append(finalResources, resource)
	}

	if err := validations.ValidateMachineTypes(finalResources); err != nil {
		return nil, err
	}

	return finalResources, nil
}

func ResourcesToYamlMachineTypes(resources []map[string]interface{}, outputFile string) error {

	// Create output file if 1) it does not exist and 2) if the file is intended to be placed there
	if strings.HasPrefix(outputFile, "output/") {
		err := os.Mkdir("output", 0777)
		if err != nil && !strings.Contains(err.Error(), "mkdir output: file exists") {
			return err
		}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	// Machine types
	i := 0
	_, err = file.Write([]byte("---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  " +
	"annotations:\n    " + "qontract.recycle: \"true\"\n  " + "name: cloud-resources-config\ndata:\n  instance-types.yaml: |\n    instance_types:\n"))
	if err != nil {
		panic(err)
	}
	for _, resource := range resources {
		if resource[ResourceType] == "compute.node" {
			// Create block for this machine type
			id := helper.AssignStringValue(resource["id"])
			machineTypeString := []byte(fmt.Sprintf(
				"    - id: %s\n"+
					"      name: %s\n"+
					"      cloud_provider_id: %s\n"+
					"      cpu_cores: %s\n"+
					"      memory: %s\n"+
					"      category: %s\n"+
					"      size: %s\n"+
					"      generic_name: %s\n",
				id, resource[Name].(string), resource[CloudProviderId].(string), resource[CpuCores].(string),
				resource[Memory].(string), resource[Category].(string), resource[SizePretty].(string),
				resource[GenericName].(string)))
			_, err = file.Write(machineTypeString)
			if err != nil {
				panic(err)
			}
			i++
		}
	}
	return nil
}

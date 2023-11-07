package convert

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/openshift-online/ocm-csv-parser/pkg/helper"
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

	// Parse resources from CSV
	for _, record := range records {
		resource := make(map[string]interface{})
		resource["id"] = record[0]
		resource["generic_name"] = record[1]
		resource["name"] = record[2]
		resource["cloud_provider_id"] = record[3]
		resource["cpu_cores"] = record[4]
		resource["memory"] = record[5]
		resource["memory_pretty"] = record[6]
		resource["category"] = record[7]
		resource["category_pretty"] = record[8]
		resource["size_pretty"] = record[9]
		resource["resource_type"] = record[10]
		resource["active"] = record[11]
		finalResources = append(finalResources, resource)
	}

	return finalResources, nil
}

func ResourcesToYamlMachineTypes(resources []map[string]interface{}, outputFile string) error {

	// Create output file if 1) it does not exist and 2) if the file is intended to be placed there
	if strings.Contains(outputFile, "output/") {
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
		"name: cloud-resources-config\ndata:\n  instance-types.yaml: |\n    instance_types:\n"))
	if err != nil {
		panic(err)
	}
	for _, resource := range resources {
		if resource["resource_type"] == "compute.node" {
			// Create block for this machine type
			id := helper.AssignStringValue(resource["id"])
			machineTypeString := []byte("    - id: " + id +
				"\n      name: " + resource["name"].(string) + "\n      cloud_provider_id: " +
				resource["cloud_provider_id"].(string) + "\n      cpu_cores: " +
				resource["cpu_cores"].(string) + "\n      memory: " + resource["memory"].(string) +
				"\n      category: " +
				resource["category"].(string) + "\n      size: " + resource["size_pretty"].(string) +
				"\n      generic_name: " +
				resource["generic_name"].(string) + "\n")
			_, err = file.Write(machineTypeString)
			if err != nil {
				panic(err)
			}
			i++
		}
	}
	return nil
}

func ResourcesToConstraintMap(resources []map[string]interface{}, outputFile string) error {
	if strings.Contains(outputFile, "output/") {
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
		"name: cloud-resources-constraints\ndata:\n  constraints.yaml: |\n    instance_types:\n"))
	if err != nil {
		panic(err)
	}
	for _, m := range resources {
		if m["resource_type"] == "compute.node" {
			id := helper.AssignStringValue(m["id"])
			machineTypeString := []byte("    - id: " + id + "\n      enabled: " + strings.ToLower(m["active"].(string)) + "\n")
			_, err = file.Write(machineTypeString)
			if err != nil {
				panic(err)
			}
			i++
		}
	}
	return nil
}

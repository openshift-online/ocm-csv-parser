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

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, rec := range records {
		resource := make(map[string]interface{})
		resource["id"] = rec[0]
		resource["generic_name"] = rec[1]
		resource["name"] = rec[2]
		resource["cloud_provider_id"] = rec[3]
		resource["cpu_cores"] = rec[4]
		resource["memory"] = rec[5]
		resource["memory_pretty"] = rec[6]
		resource["category"] = rec[7]
		resource["category_pretty"] = rec[8]
		resource["size_pretty"] = rec[9]
		resource["resource_type"] = rec[10]
		resource["active"] = rec[11]
		finalResources = append(finalResources, resource)
	}

	return finalResources, nil
}

func ResourcesToYaml(resources []map[string]interface{}, outputFile string) error {

	if strings.Contains(outputFile, "output/") {
		err := os.Mkdir("output", 0777)
		if err != nil && !strings.Contains(err.Error(), "mkdir output: file exists") {
			return err
		}
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	// Machine types
	i := 0
	_, err = f.Write([]byte("---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  " +
		"name: cloud-resources-config\ndata:\n  instance-types.yaml: |\n    instance_types:\n"))
	if err != nil {
		panic(err)
	}
	for _, m := range resources {
		if m["resource_type"] == "compute.node" {
			id := helper.AssignValue(m["id"])
			mString := []byte("    - id: " + id.(string) +
				"\n      name: " + m["name"].(string) + "\n      cloud_provider_id: " +
				m["cloud_provider_id"].(string) + "\n      cpu_cores: " +
				m["cpu_cores"].(string) + "\n      memory: " + m["memory"].(string) + "\n      category: " +
				m["category"].(string) + "\n      size: " + m["size_pretty"].(string) + "\n      generic_name: " +
				m["generic_name"].(string) + "\n")
			_, err = f.Write(mString)
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

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	// Machine types
	i := 0
	_, err = f.Write([]byte("---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  " +
		"name: cloud-resources-constraints\ndata:\n  constraints.yaml: |\n    instance_types:\n"))
	if err != nil {
		panic(err)
	}
	for _, m := range resources {
		if m["resource_type"] == "compute.node" {
			id := helper.AssignValue(m["id"])
			mString := []byte("    - id: " + id.(string) + "\n      enabled: " + strings.ToLower(m["active"].(string)) + "\n")
			_, err = f.Write(mString)
			if err != nil {
				panic(err)
			}
			i++
		}
	}
	return nil
}

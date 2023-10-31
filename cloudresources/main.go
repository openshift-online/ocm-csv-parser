package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	resources []map[string]interface{}
)

func csvToResources() {
	csv_file, err := os.Open("file.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csv_file.Close()

	r := csv.NewReader(csv_file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		resources = append(resources, resource)
	}
}

func resourcesToYaml() {

	f, err := os.Create("cloud-resources.configmap.yaml")
	if err != nil {
		fmt.Println("Error creating yaml file")
		return
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
			id := assignValue(m["id"])
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

}

func resourcesToConstraintMap() {

	f, err := os.Create("cloud-resource-constraints.configmap.yaml")
	if err != nil {
		fmt.Println("Error creating constraints yaml file")
		return
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
			id := assignValue(m["id"])
			mString := []byte("    - id: " + id.(string) + "\n      enabled: " + strings.ToLower(m["active"].(string)) + "\n")
			_, err = f.Write(mString)
			if err != nil {
				panic(err)
			}
			i++
		}
	}

}

type ccsOnly struct {
	ID      string `json:"id"`
	CcsOnly string `json:"ccs_only"`
}

func assignValue(value interface{}) interface{} {
	if value == nil {
		return ""
	}
	return value
}

func main() {
	start := time.Now()
	csvToResources()
	resourcesToYaml()
	//resourcesToConstraintMap()
	fmt.Printf("Process complete: ", time.Since(start))
}

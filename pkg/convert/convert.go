package convert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	. "github.com/openshift-online/ocm-csv-parser/pkg/fieldnames"
	"github.com/openshift-online/ocm-csv-parser/pkg/helper"
	"github.com/openshift-online/ocm-csv-parser/pkg/validations"
)

func MachineTypesCsvToResources(resourcesFile string) ([]map[string]interface{}, error) {
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

func RegionsCsvToResources(resourcesFile string) ([]map[string]interface{}, error) {
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

	keys := []string{Id, CloudProviderId, DisplayName, SupportsMultiAz}

	// Parse resources from CSV
	for _, record := range records {
		resource := make(map[string]interface{})
		for i, key := range keys {
			resource[key] = record[i]
		}
		finalResources = append(finalResources, resource)
	}

	if err := validations.ValidateRegions(finalResources); err != nil {
		return nil, err
	}

	return finalResources, nil
}

func ResourcesToYamlRegions(resources []map[string]interface{}, outputFile string, appInterfaceDir string) error {

	if appInterfaceDir == "" {
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

		// Regions
		_, err = file.Write([]byte("---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  " +
			"annotations:\n    " + "qontract.recycle: \"true\"\n  " + "name: cloud-resources-config\ndata:\n  cloud-regions.yaml: |\n    cloud_regions:\n"))
		if err != nil {
			panic(err)
		}
		for i, resource := range resources {
			if i == 0 {
				continue
			}
			// Create block for this region
			id := helper.AssignStringValue(resource["id"])
			regionString := []byte(fmt.Sprintf(
				"    - id: %s\n"+
					"      cloud_provider_id: %s\n"+
					"      display_name: %s\n"+
					"      supports_multi_az: %s\n",
				id, resource[CloudProviderId].(string), resource[DisplayName].(string), resource[SupportsMultiAz].(string)))
			_, err = file.Write(regionString)
			if err != nil {
				panic(err)
			}
		}
		file.Close()
	} else {
		if strings.HasSuffix(appInterfaceDir, "/") {
			appInterfaceDir = appInterfaceDir[:len(appInterfaceDir)-1]
		}
		appInterfaceConfigmapFile := fmt.Sprintf("%s/resources/services/ocm/cloud-resources.configmap.yaml", appInterfaceDir)
		file, err := os.Open(appInterfaceConfigmapFile)
		if err != nil {
			return err
		}

		defer file.Close()

		fmt.Println(fmt.Sprintf("Rewriting app-interface configmap (%s) with new data supplied from CSV...\n", appInterfaceConfigmapFile))

		scanner := bufio.NewScanner(file)
		finalText := ""
		newConfigmap := ""
		for i, resource := range resources {
			if i == 0 {
				continue
			}

			id := helper.AssignStringValue(resource["id"])

			data := fmt.Sprintf(
				"    - id: %s\n"+
					"      cloud_provider_id: %s\n"+
					"      display_name: %s\n"+
					"      supports_multi_az: %s\n",
				id, resource[CloudProviderId].(string), resource[DisplayName].(string), resource[SupportsMultiAz].(string))

			newConfigmap += string(data)
		}

		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, "cloud_regions:") {
				finalText += fmt.Sprintf("%s\n%s", text, newConfigmap)
				if !strings.HasSuffix(newConfigmap, "\n") {
					finalText += "\n"
				}
				break
			} else {
				finalText += fmt.Sprintf("%s\n", text)
			}
		}

		if err := os.Truncate(appInterfaceConfigmapFile, 0); err != nil {
			return err
		}

		newFile, err := os.OpenFile(appInterfaceConfigmapFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		defer newFile.Close()

		_, err = newFile.Write([]byte(finalText))
		if err != nil {
			return err
		}

	}
	return nil
}

func ResourcesToYamlMachineTypes(resources []map[string]interface{}, outputFile string, appInterfaceDir string) error {

	if appInterfaceDir == "" {
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
	} else {
		if appInterfaceDir[len(appInterfaceDir)-1] == '/' {
			appInterfaceDir = appInterfaceDir[:len(appInterfaceDir)-1]
		}
		appInterfaceConfigmapFile := fmt.Sprintf("%s/resources/services/ocm/cloud-resources.configmap.yaml", appInterfaceDir)
		file, err := os.Open(appInterfaceConfigmapFile)
		if err != nil {
			return err
		}

		defer file.Close()

		fmt.Println(fmt.Sprintf("Rewriting app-interface configmap (%s) with new data supplied from CSV...\n", appInterfaceConfigmapFile))

		scanner := bufio.NewScanner(file)
		finalText := ""
		newConfigmap := ""
		for i, resource := range resources {
			if resource[ResourceType] == "compute.node" {
				if i == 0 {
					continue
				}

				id := helper.AssignStringValue(resource["id"])

				data := []byte(fmt.Sprintf(
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

				newConfigmap += string(data)
			}
		}

		skip := false
		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, "instance_types:") {
				finalText += fmt.Sprintf("%s\n%s", text, newConfigmap)
				skip = true
			} else if !skip {
				finalText += fmt.Sprintf("%s\n", text)
			} else if strings.Contains(text, "cloud-regions.yaml: |") {
				skip = false
				finalText += fmt.Sprintf("%s\n", text)
			}
		}

		if err := os.Truncate(appInterfaceConfigmapFile, 0); err != nil {
			return err
		}

		newFile, err := os.OpenFile(appInterfaceConfigmapFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		defer newFile.Close()

		_, err = newFile.Write([]byte(finalText))
		if err != nil {
			return err
		}
	}
	return nil
}

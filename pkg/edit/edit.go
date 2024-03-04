package edit

import (
	"fmt"
	"os"
	"path"
	"strings"

	errors "github.com/zgalor/weberr"
)

type machineType struct {
	id      string
	ccsOnly string
	enabled string
}

type cloudRegion struct {
	id       string
	enabled  string
	govcloud string
	ccsOnly  string
}

type ConstraintMap struct {
	entireFile   []string
	machineTypes []machineType
	cloudRegions []cloudRegion
	Path         string
}

func NewConstraintMap(env, appInterfaceDir string) (ConstraintMap, error) {
	constraintMap := ConstraintMap{}
	err := constraintMap.readConstraints(env, appInterfaceDir)
	return constraintMap, err
}

func (c *ConstraintMap) readConstraints(env string, appInterfaceDir string) error {
	c.Path = path.Join(appInterfaceDir,
		"resources/services/ocm/"+env+"/cloud-resource-constraints.configmap.yaml")
	yamlFile, err := os.ReadFile(c.Path)
	if err != nil {
		return errors.Errorf("File read err: '%v'", err)
	}

	for _, line := range strings.Split(strings.TrimSuffix(string(yamlFile), "\n"), "\n") {
		c.entireFile = append(c.entireFile, line)
	}

	return nil
}

func (c *ConstraintMap) EditConstraint(id string, ccsOnly *bool, enabled *bool, govcloud *bool) error {
	edit := false
	changeMade := false
	for i, line := range c.entireFile {
		if strings.Contains(line, "id:") {
			edit = false
			if strings.Contains(line, id) {
				edit = true
			}
		}
		if edit {
			if ccsOnly != nil && strings.Contains(line, "ccs_only") {
				s := "false"
				if *ccsOnly {
					s = "true"
				}
				changeMade = true
				c.entireFile[i] = strings.Split(line, ": ")[0] + ": " + s
			} else if enabled != nil && strings.Contains(line, "enabled") {
				s := "false"
				if *enabled {
					s = "true"
				}
				changeMade = true
				c.entireFile[i] = strings.Split(line, ": ")[0] + ": " + s
			} else if govcloud != nil && strings.Contains(line, "govcloud") {
				s := "false"
				if *govcloud {
					s = "true"
				}
				changeMade = true
				c.entireFile[i] = strings.Split(line, ": ")[0] + ": " + s
			}
		}
	}
	if !changeMade {
		fmt.Printf("\nNo changes were made, writing same data back into constraintmap " +
			"(use '--help' to see what you can change via flags)")
	}
	err := c.saveFile()
	return err
}

func (c *ConstraintMap) saveFile() error {

	if err := os.Truncate(c.Path, 0); err != nil {
		return err
	}

	newFile, err := os.OpenFile(c.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer newFile.Close()

	_, err = newFile.Write([]byte(strings.Join(c.entireFile, "\n")))
	if err != nil {
		return err
	}

	return nil
}

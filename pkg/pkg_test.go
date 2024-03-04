package pkg_test

import (
	"io"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-csv-parser/pkg/convert"
	"github.com/openshift-online/ocm-csv-parser/pkg/edit"
	"github.com/openshift-online/ocm-csv-parser/pkg/errormsgs"
	. "github.com/openshift-online/ocm-csv-parser/pkg/fieldnames"
	"github.com/openshift-online/ocm-csv-parser/pkg/validations"
	errors "github.com/zgalor/weberr"
)

const (
	testOutputFile = `---
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    qontract.recycle: "true"
  name: cloud-resources-config
data:
  instance-types.yaml: |
    instance_types:
    - id: m7a.12xlarge
      name: m7a.12xlarge - General purpose
      cloud_provider_id: aws
      cpu_cores: 48
      memory: 206158430208
      category: general_purpose
      size: 12xlarge
      generic_name: standard-48-m7a
`
)

var _ = Describe("Pkg", func() {
	Context("Functions should return the correct output", Ordered, func() {
		var machineTypeId string
		var machineTypes []map[string]interface{}
		var outputFile string
		Describe("parse command", func() {
			BeforeEach(func() {
				machineTypeId = "m7a.12xlarge"
				outputFile = "test/test.configmap.yaml"
			})
			AfterAll(func() {
				os.Remove("test/test.configmap.yaml")
			})
			It("Convert single resource from CSV into YAML", func() {
				var err error
				machineTypes, err = convert.MachineTypesCsvToResources("test/test.csv")
				Expect(err).ToNot(HaveOccurred())
				machineType := machineTypes[0]
				Expect(machineType[Id]).To(Equal(machineTypeId))
				Expect(machineType[SizePretty]).To(Equal("12xlarge"))
				Expect(machineType[CpuCores]).To(Equal("48"))
			})
			It("Convert parsed resource into yaml", func() {
				err := convert.ResourcesToYamlMachineTypes(machineTypes, outputFile, "")
				Expect(err).ToNot(HaveOccurred())
				testFile, err := os.ReadFile(outputFile)

				Expect(err).ToNot(HaveOccurred())
				Expect(testFile).To(ContainSubstring(machineTypeId))
				Expect(string(testFile)).To(Equal(testOutputFile))
			})
			It("Test output file creation", func() {
				err := convert.ResourcesToYamlMachineTypes(machineTypes, "output/test.configmap.yaml", "")
				Expect(err).ToNot(HaveOccurred())
				file, err := os.ReadFile("output/test.configmap.yaml")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(file)).To(Equal(testOutputFile))
				err = os.Remove("output/test.configmap.yaml")
				Expect(err).ToNot(HaveOccurred())
				err = os.Remove("output")
				Expect(err).ToNot(HaveOccurred())
				// File creation should only happen for the 'default' option (same dir, inside a dir called 'output')
				err = convert.ResourcesToYamlMachineTypes(machineTypes, "otpt/test.configmap.yaml", "")
				Expect(err).To(HaveOccurred())
			})
			It("Fail for each field supplied with the wrong data type", func() {
				// Check function which is to be used by the parse command
				machineTypes, err := convert.MachineTypesCsvToResources("test/badDataTypes.csv")
				Expect(err).To(HaveOccurred())
				Expect(machineTypes).To(BeNil())
				// Check each specific field validation
				machineTypes, err = convert.MachineTypesCsvToResources("test/test.csv")
				Expect(err).ToNot(HaveOccurred())
				machineType := machineTypes[0]
				machineTypeId := machineType[Id].(string)
				Expect(makeFakeMachineTypeDataType(machineType, GenericName, 1234).Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, GenericName).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, SizePretty, 1234).Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, SizePretty).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, Category, 1234).Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, Category).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, Memory, "aaa").Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, Memory).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, CpuCores, "aaa").Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, CpuCores).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, CloudProviderId, 1234).Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, CloudProviderId).Error()))
				Expect(makeFakeMachineTypeDataType(machineType, Name, 1234).Error()).
					To(Equal(fakeDataTypeValidationError(machineTypeId, Name).Error()))
			})
			It("Fail for missing cloud provider ID with regions", func() {
				cloudRegions, err := convert.RegionsCsvToResources("test/testMissingCloudProvider.csv")
				Expect(err).To(HaveOccurred())
				Expect(cloudRegions).To(BeNil())
				Expect(err.Error()).To(Equal(fakeMissingDataError("eastasia", CloudProviderId).Error()))
			})
			It("Pass for non-missing cloud provider ID with regions", func() {
				cloudRegions, err := convert.RegionsCsvToResources("test/testCloudRegions.csv")
				Expect(err).ToNot(HaveOccurred())
				Expect(cloudRegions).ToNot(BeNil())
				Expect(cloudRegions[0][Id].(string)).To(Equal("eastasia"))
				Expect(cloudRegions[0][CloudProviderId].(string)).To(Equal("azure"))
				Expect(cloudRegions[0][DisplayName].(string)).To(Equal("East Asia"))
				Expect(cloudRegions[0][SupportsMultiAz]).To(Equal("true"))
				Expect(cloudRegions[1][Id].(string)).To(Equal("eastasia2"))
				Expect(cloudRegions[1][CloudProviderId].(string)).To(Equal("gcp"))
				Expect(cloudRegions[1][DisplayName].(string)).To(Equal("East Asia 2"))
				Expect(cloudRegions[1][SupportsMultiAz]).To(Equal("false"))
				Expect(cloudRegions[2][Id].(string)).To(Equal("eastasia3"))
				Expect(cloudRegions[2][CloudProviderId].(string)).To(Equal("aws"))
				Expect(cloudRegions[2][DisplayName].(string)).To(Equal("East Asia 3"))
				Expect(cloudRegions[2][SupportsMultiAz]).To(Equal("true"))
			})
			It("Test full app-interface flow", func() {
				var err error
				machineTypes, err = convert.MachineTypesCsvToResources("test/app-interface-machine-types-test.csv")
				Expect(err).ToNot(HaveOccurred())
				err = convert.ResourcesToYamlMachineTypes(machineTypes, "", "test/")
				Expect(err).ToNot(HaveOccurred())

				cloudRegions, err := convert.RegionsCsvToResources("test/app-interface-regions-test.csv")
				Expect(err).ToNot(HaveOccurred())
				err = convert.ResourcesToYamlRegions(cloudRegions, "", "test/")
				Expect(err).ToNot(HaveOccurred())

				output, err := os.ReadFile("test/resources/services/ocm/cloud-resources.configmap.yaml")
				Expect(err).ToNot(HaveOccurred())

				expected, err := os.ReadFile("test/after-testing-app-interface.configmap.yaml")
				Expect(err).ToNot(HaveOccurred())

				Expect(string(output)).To(Equal(string(expected)))
			})
		})
	})

	Context("Test editing region and machine type in constraintmap in app-interface", func() {
		Describe("edit command", Ordered, func() {
			originalFilePath := "test/resources/services/ocm/stage/cloud-resource-constraints.configmap.yaml"
			originalFileDataPath := "test/original-edit.configmap.yaml"
			afterEditFilePath := "test/after-testing-edit.configmap.yaml"
			AfterAll(func() {
				// Return test data to original data
				err := os.Remove(originalFilePath)
				Expect(err).ToNot(HaveOccurred())
				destination, err := os.Open(originalFileDataPath)
				Expect(err).ToNot(HaveOccurred())
				defer destination.Close()
				source, err := os.Create(originalFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer source.Close()
				_, err = io.Copy(source, destination)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Test edit machinetype and edit region from constraints", func() {

				// Machine types

				var constraintMap edit.ConstraintMap

				ccsOnly := false
				enabled := false
				govcloud := true

				appInterfaceDir := "test"

				err := constraintMap.ReadConstraints("stage", appInterfaceDir)
				Expect(err).ToNot(HaveOccurred())

				machineTypeId := "c5.24xlarge"

				Expect(constraintMap.Path).To(Equal(originalFilePath))

				err = constraintMap.EditConstraint(machineTypeId, &ccsOnly, &enabled, nil)
				Expect(err).ToNot(HaveOccurred())

				err = constraintMap.SaveFile()
				Expect(err).ToNot(HaveOccurred())

				// Regions

				regionId := "me-central-1"

				Expect(constraintMap.Path).To(Equal(originalFilePath))

				err = constraintMap.EditConstraint(regionId, &ccsOnly, &enabled, &govcloud)
				Expect(err).ToNot(HaveOccurred())

				err = constraintMap.SaveFile()
				Expect(err).ToNot(HaveOccurred())

				if originalEditedFile, err := os.ReadFile(constraintMap.Path); Expect(err).ToNot(HaveOccurred()) {
					if afterEditFile, err := os.ReadFile(
						afterEditFilePath); Expect(err).ToNot(HaveOccurred()) {
						after := strings.Split(string(afterEditFile), "\n")
						originalEdited := strings.Split(string(originalEditedFile), "\n")
						for i, _ := range after {
							if i < len(after) && i < len(originalEdited) {
								Expect(after[i]).To(Equal(originalEdited[i]))
							}
						}
					}
				}
			})
		})
	})
})

func makeFakeMachineTypeDataType(machineType map[string]interface{}, field string, value interface{}) error {
	machineType[field] = value
	err := validations.ValidateMachineTypes([]map[string]interface{}{machineType})
	return err
}

func fakeDataTypeValidationError(resourceId, resourceField string) error {
	return errors.UserErrorf(errormsgs.SyntaxDataType, resourceId, resourceField)
}

func fakeMissingDataError(resourceId, resourceField string) error {
	return errors.UserErrorf(errormsgs.EmptyRequiredField, resourceField, resourceId)
}

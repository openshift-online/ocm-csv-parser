package pkg_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-csv-parser/pkg/convert"
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
		BeforeEach(func() {
			machineTypeId = "m7a.12xlarge"
			outputFile = "test/test.configmap.yaml"
		})
		AfterAll(func() {
			os.Remove("test/test.configmap.yaml")
		})
		It("Convert single resource from CSV into YAML", func() {
			var err error
			machineTypes, err = convert.CsvToResources("test/test.csv")
			Expect(err).ToNot(HaveOccurred())
			machineType := machineTypes[0]
			Expect(machineType[Id]).To(Equal(machineTypeId))
			Expect(machineType[SizePretty]).To(Equal("12xlarge"))
			Expect(machineType[CpuCores]).To(Equal("48"))
		})
		It("Convert parsed resource into yaml", func() {
			err := convert.ResourcesToYamlMachineTypes(machineTypes, outputFile)
			Expect(err).ToNot(HaveOccurred())
			testFile, err := os.ReadFile(outputFile)

			Expect(err).ToNot(HaveOccurred())
			Expect(testFile).To(ContainSubstring(machineTypeId))
			Expect(string(testFile)).To(Equal(testOutputFile))
		})
		It("Test output file creation", func() {
			err := convert.ResourcesToYamlMachineTypes(machineTypes, "output/test.configmap.yaml")
			Expect(err).ToNot(HaveOccurred())
			file, err := os.ReadFile("output/test.configmap.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(string(file)).To(Equal(testOutputFile))
			err = os.Remove("output/test.configmap.yaml")
			Expect(err).ToNot(HaveOccurred())
			err = os.Remove("output")
			Expect(err).ToNot(HaveOccurred())
			// File creation should only happen for the 'default' option (same dir, inside a dir called 'output')
			err = convert.ResourcesToYamlMachineTypes(machineTypes, "otpt/test.configmap.yaml")
			Expect(err).To(HaveOccurred())
		})
		It("Fail for each field supplied with the wrong data type", func() {
			// Check function which is to be used by the parse command
			machineTypes, err := convert.CsvToResources("test/badDataTypes.csv")
			Expect(err).To(HaveOccurred())
			Expect(machineTypes).To(BeNil())
			// Check each specific field validation
			machineTypes, err = convert.CsvToResources("test/test.csv")
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

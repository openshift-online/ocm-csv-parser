package pkg_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-csv-parser/pkg/convert"
)

var _ = Describe("Pkg", func() {
	Context("Functions should return the correct output", func() {
		var machineTypeId string
		var machineTypes []map[string]interface{}
		var outputFile string
		BeforeEach(func() {
			machineTypeId = "m7a.12xlarge"
			outputFile = "test/test.configmap.yaml"
		})
		It("Convert single resource from CSV into YAML", func() {
			var err error
			machineTypes, err = convert.CsvToResources("test/test.csv")
			Expect(err).To(BeNil())
			machineType := machineTypes[0]
			Expect(machineType["id"]).To(Equal(machineTypeId))
			Expect(machineType["size_pretty"]).To(Equal("12xlarge"))
			Expect(machineType["cpu_cores"]).To(Equal("48"))
		})
		It("Convert parsed resource into yaml", func() {
			err := convert.ResourcesToYamlMachineTypes(machineTypes, outputFile)
			Expect(err).To(BeNil())
			testFile, err := os.ReadFile(outputFile)

			Expect(testFile).To(ContainSubstring(machineTypeId))
			Expect(testFile).To(ContainSubstring("12xlarge"))
			Expect(testFile).To(ContainSubstring("48"))
		})
	})
})

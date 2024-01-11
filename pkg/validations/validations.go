package validations

import (
	"reflect"
	"strconv"

	"github.com/openshift-online/ocm-csv-parser/pkg/errormsgs"
	. "github.com/openshift-online/ocm-csv-parser/pkg/fieldnames"
	errors "github.com/zgalor/weberr"
)

func validateDataType(resourceId, resourceField string, resourceValue interface{}, dataType reflect.Type) error {
	if reflect.TypeOf(resourceValue) != dataType {
		return errors.UserErrorf(errormsgs.SyntaxDataType, resourceId, resourceField)
	}
	return nil
}

func validateMachineTypeDataTypes(resource map[string]interface{}) error {
	resourceId := resource[Id].(string)
	if err := validateDataType(resourceId, Name, resource[Name], reflect.TypeOf("")); err != nil {
		return err
	}
	if err := validateDataType(resourceId, CloudProviderId, resource[CloudProviderId],
		reflect.TypeOf("")); err != nil {
		return err
	}
	cores, err := strconv.Atoi(resource[CpuCores].(string))
	if err != nil {
		return errors.UserErrorf(errormsgs.SyntaxDataType, resourceId, CpuCores)
	}
	if err := validateDataType(resourceId, CpuCores, cores, reflect.TypeOf(0)); err != nil {
		return err
	}
	memory, err := strconv.ParseInt(resource[Memory].(string), 10, 64)
	if err != nil {
		return errors.UserErrorf(errormsgs.SyntaxDataType, resourceId, Memory)
	}
	if err := validateDataType(resourceId, Memory, memory, reflect.TypeOf(int64(0))); err != nil {
		return err
	}
	if err := validateDataType(resourceId, Category, resource[Category], reflect.TypeOf("")); err != nil {
		return err
	}
	if err := validateDataType(resourceId, SizePretty, resource[SizePretty], reflect.TypeOf("")); err != nil {
		return err
	}
	if err := validateDataType(resourceId, GenericName, resource[GenericName], reflect.TypeOf("")); err != nil {
		return err
	}

	return nil
}

func validateRegionDataTypes(resource map[string]interface{}) error {
	resourceId := resource[Id].(string)
	if resource[CloudProviderId] == "" {
		return errors.UserErrorf(errormsgs.EmptyRequiredField, CloudProviderId, resourceId)
	}
	if err := validateDataType(resourceId, CloudProviderId, resource[CloudProviderId],
		reflect.TypeOf("")); err != nil {
		return err
	}
	if err := validateDataType(resourceId, DisplayName, resource[DisplayName],
		reflect.TypeOf("")); err != nil {
		return err
	}
	if err := validateDataType(resourceId, SupportsMultiAz, resource[SupportsMultiAz],
		reflect.TypeOf("")); err != nil {
		return err
	}

	return nil
}

func ValidateMachineTypes(resources []map[string]interface{}) error {
	for _, resource := range resources {
		if resource[ResourceType] == "compute.node" {
			if err := validateMachineTypeDataTypes(resource); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateRegions(resources []map[string]interface{}) error {
	for _, resource := range resources {
		if err := validateRegionDataTypes(resource); err != nil {
			return err
		}
	}
	return nil
}

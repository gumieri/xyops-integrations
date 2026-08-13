package provider

import (
	"context"
	"log"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringSliceFromSet converts a framework List of strings to a Go slice.
func stringSliceFromSet(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]types.String, 0, len(list.Elements()))
	diags := list.ElementsAs(context.Background(), &elems, false)
	if diags.HasError() {
		log.Printf("warning: stringSliceFromSet: dropped diagnostics: %v", diags.Errors())
		return nil
	}

	out := make([]string, 0, len(elems))
	for _, e := range elems {
		out = append(out, e.ValueString())
	}
	return out
}

// stringSetFromSlice converts a Go slice of strings to a framework List.
func stringSetFromSlice(slice []string) types.List {
	if slice == nil {
		return types.ListNull(types.StringType)
	}

	elems := make([]types.String, 0, len(slice))
	for _, s := range slice {
		elems = append(elems, types.StringValue(s))
	}

	list, diags := types.ListValueFrom(context.Background(), types.StringType, elems)
	if diags.HasError() {
		log.Printf("warning: stringSetFromSlice: dropped diagnostics: %v", diags.Errors())
		return types.ListNull(types.StringType)
	}
	return list
}

// intSliceFromSet converts a framework List of int64 to a Go slice.
func intSliceFromSet(list types.List) []int {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]types.Int64, 0, len(list.Elements()))
	diags := list.ElementsAs(context.Background(), &elems, false)
	if diags.HasError() {
		log.Printf("warning: intSliceFromSet: dropped diagnostics: %v", diags.Errors())
		return nil
	}

	out := make([]int, 0, len(elems))
	for _, e := range elems {
		out = append(out, int(e.ValueInt64()))
	}
	return out
}

// intSetFromSlice converts a Go slice of ints to a framework List.
func intSetFromSlice(slice []int) types.List {
	if slice == nil {
		return types.ListNull(types.Int64Type)
	}

	elems := make([]types.Int64, 0, len(slice))
	for _, s := range slice {
		elems = append(elems, types.Int64Value(int64(s)))
	}

	list, diags := types.ListValueFrom(context.Background(), types.Int64Type, elems)
	if diags.HasError() {
		log.Printf("warning: intSetFromSlice: dropped diagnostics: %v", diags.Errors())
		return types.ListNull(types.Int64Type)
	}
	return list
}

// frameworkMapToBoolMap converts a framework Map of bools to sdk.BoolMap.
func frameworkMapToBoolMap(m types.Map) sdk.BoolMap {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}

	out := make(sdk.BoolMap, len(m.Elements()))
	for key, elem := range m.Elements() {
		if b, ok := elem.(types.Bool); ok {
			out[key] = b.ValueBool()
		}
	}
	return out
}

// boolMapToFrameworkMap converts sdk.BoolMap to a framework Map of bools.
func boolMapToFrameworkMap(m sdk.BoolMap) types.Map {
	if m == nil {
		return types.MapNull(types.BoolType)
	}

	elems := make(map[string]types.Bool, len(m))
	for k, v := range m {
		elems[k] = types.BoolValue(v)
	}

	value, diags := types.MapValueFrom(context.Background(), types.BoolType, elems)
	if diags.HasError() {
		log.Printf("warning: boolMapToFrameworkMap: dropped diagnostics: %v", diags.Errors())
		return types.MapNull(types.BoolType)
	}
	return value
}

// intToBool converts 1/0 to a bool.
func intToBool(v int) bool {
	return v != 0
}

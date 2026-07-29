package provider

import (
	"context"
	"log"

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

// boolMapToIntMap converts a framework Map of bools to a Go map[string]int (1/0).
func boolMapToIntMap(m types.Map) map[string]int {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}

	out := make(map[string]int, len(m.Elements()))
	for key, elem := range m.Elements() {
		if b, ok := elem.(types.Bool); ok {
			if b.ValueBool() {
				out[key] = 1
			} else {
				out[key] = 0
			}
		}
	}
	return out
}

// intMapToBoolMap converts a Go map[string]int (1/0) to a framework Map of bools.
func intMapToBoolMap(m map[string]int) types.Map {
	if m == nil {
		return types.MapNull(types.BoolType)
	}

	elems := make(map[string]types.Bool, len(m))
	for key, v := range m {
		elems[key] = types.BoolValue(v == 1)
	}

	value, diags := types.MapValueFrom(context.Background(), types.BoolType, elems)
	if diags.HasError() {
		log.Printf("warning: intMapToBoolMap: dropped diagnostics: %v", diags.Errors())
		return types.MapNull(types.BoolType)
	}
	return value
}

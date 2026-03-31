package conftagz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvFieldSubstitutionFromMap(t *testing.T) {
	mystruct := MyStruct{"Value1", "Value2", 3, 0, 0, 0, false}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"VeryImportant":      "NewValue2",
		"ExtremelyImportant": "123",
		"Field4":             "1", // bool
	}

	expected := []string{"Field1", "Field2", "Field3", "Field3a", "Field3b"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, &EnvFieldSubstOpts{ThrowErrorIfEnvMissing: true}, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "NewValue1")
	assert.Equal(t, mystruct.Field2, "NewValue2")
	assert.Equal(t, mystruct.Field3, 123)

}

func TestEnvFieldWithPrivateFieldTagShouldFail(t *testing.T) {
	mystruct := MyStructWithPrivateAndTag{"Value1", 0}

	envMap := map[string]string{
		"Important": "NewValue1",
		"WONTWORK":  "NewValue2",
	}

	expected := []string{"Field1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, &EnvFieldSubstOpts{ThrowErrorIfEnvMissing: true}, envMap)
	if err != nil {
		t.Errorf("Error %v", result)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// assumea PATH is there
func TestEnvFieldSubstitutionFromEnviron(t *testing.T) {
	mystruct := MyStruct2{"Value1"}
	expected := []string{"Path"}

	result, err := EnvFieldSubstitution(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	assert.NotEqual(t, mystruct.Path, "Value1")
}

func TestEnvFieldSubstitutionFromMap3(t *testing.T) {
	mystruct := MyStruct{"Value1", "Value2", 3, 0, 0, 0, false}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"NOPE":               "NewValue2",
		"ExtremelyImportant": "123",
	}

	expected := []string{"Field1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, &EnvFieldSubstOpts{ThrowErrorIfEnvMissing: true}, envMap)
	if err != nil {
		// expected
	} else {
		t.Errorf("Expected error, but got %v", result)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "NewValue1", mystruct.Field1)
	// assert.Equal(t, 123, mystruct.Field3)
}
func TestEnvFieldSubstitutionStructWithStruct(t *testing.T) {
	mystruct := MyStructWithStruct{"Value1", "Value2", 3, "", nil, nil, nil, InnerStruct{"InnerValue1", nil, innerStruct{0}}}

	envMap := map[string]string{
		"ENV1":   "NewValue1",
		"ENV2":   "NewValue2",
		"ENV3":   "123",
		"INNER1": "InnerValue1",
	}

	expected := []string{"Field1", "Field2", "Field3", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "NewValue1")
	assert.Equal(t, mystruct.Field3, 123)
	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerValue1")
}

func TestEnvFieldSubstitutionStruct3(t *testing.T) {
	mystruct := MyStruct3{"Value1", "Value2", 3, nil, nil, nil}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"VeryImportant":      "NewValue2",
		"ExtremelyImportant": "123",
		"STUFF1":             "InnerValue1",
	}

	expected := []string{"Field1", "Field2", "Field3", "Inner2Ptr.Stuff1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "NewValue1")
	assert.Equal(t, mystruct.Field3, 123)
	assert.Equal(t, mystruct.Inner2Ptr.Stuff1, "InnerValue1")
	// test 'skip' tag
	assert.Nil(t, mystruct.InnerPtr)
}

func TestEnvFieldSubstitutionStruct3NonNilStructs(t *testing.T) {
	mystruct := MyStruct3{"Value1", "Value2", 3, &InnerStruct2{"blah"}, &InnerStruct2{"blah"}, nil}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"VeryImportant":      "NewValue2",
		"ExtremelyImportant": "123",
		"STUFF1":             "InnerValue1",
	}

	expected := []string{"Field1", "Field2", "Field3", "Inner2Ptr.Stuff1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "NewValue1")
	assert.Equal(t, mystruct.Field3, 123)
	assert.Equal(t, mystruct.Inner2Ptr.Stuff1, "InnerValue1")
	assert.Equal(t, mystruct.InnerPtr.Stuff1, "blah")
}

// Test env vars with comma-separated values for []string slices
func TestEnvFieldsWithSlice(t *testing.T) {
	mystruct := MyStructWithSlice{"", nil}

	envMap := map[string]string{
		"Important": "NewValue1",
		"SLICE":     "One,Two,Three",
	}

	expected := []string{"Field1", "SliceField"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "NewValue1", mystruct.Field1)
	assert.Equal(t, []string{"One", "Two", "Three"}, mystruct.SliceField)
}

// Test env vars with single value for []string slice
func TestEnvFieldsWithSliceSingleValue(t *testing.T) {
	mystruct := MyStructWithSlice{"", nil}

	envMap := map[string]string{
		"Important": "NewValue1",
		"SLICE":     "OnlyOne",
	}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	assert.NoError(t, err)
	assert.Contains(t, result, "SliceField")
	assert.Equal(t, []string{"OnlyOne"}, mystruct.SliceField)
}

// Test env vars for int, uint, float slices
func TestEnvFieldsWithNumericSlices(t *testing.T) {
	mystruct := MyStructWithEnvSlices{}

	envMap := map[string]string{
		"ENV_STRINGS": "a,b,c",
		"ENV_INTS":    "1, 2, 3",
		"ENV_INT32S":  "10,20,30",
		"ENV_UINTS":   "100,200,300",
		"ENV_UINT64S": "1000,2000,3000",
		"ENV_FLOATS":  "1.1, 2.2, 3.3",
	}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	assert.NoError(t, err)
	assert.Len(t, result, 6)

	assert.Equal(t, []string{"a", "b", "c"}, mystruct.StringSlice)
	assert.Equal(t, []int{1, 2, 3}, mystruct.IntSlice)
	assert.Equal(t, []int32{10, 20, 30}, mystruct.Int32Slice)
	assert.Equal(t, []uint{100, 200, 300}, mystruct.UintSlice)
	assert.Equal(t, []uint64{1000, 2000, 3000}, mystruct.Uint64Slice)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, mystruct.Float64Slice)
}

// Test env var with bad int value in slice
func TestEnvFieldsWithSliceBadInt(t *testing.T) {
	mystruct := MyStructWithEnvSlices{}

	envMap := map[string]string{
		"ENV_INTS": "1,notanumber,3",
	}

	_, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	assert.Error(t, err)
}

// Test that missing env var for slice is ignored (no error when ThrowErrorIfEnvMissing is false)
func TestEnvFieldsWithSliceMissingEnv(t *testing.T) {
	mystruct := MyStructWithSlice{"", nil}

	envMap := map[string]string{
		"Important": "NewValue1",
	}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Field1"}, result)
	assert.Nil(t, mystruct.SliceField)
}

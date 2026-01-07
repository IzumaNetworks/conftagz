package conftagz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test structs for GetDefaultOnlyStruct

type DefaultOnlyBasic struct {
	StringField  string  `default:"hello"`
	IntField     int     `default:"42"`
	UintField    uint    `default:"100"`
	FloatField   float64 `default:"3.14"`
	BoolField    bool    `default:"true"`
	NoDefault    string
	SkippedField string `conf:"skip" default:"should_not_appear"`
}

type DefaultOnlyInner struct {
	InnerString string `default:"inner_default"`
	InnerInt    int    `default:"99"`
}

type DefaultOnlyNested struct {
	TopField string           `default:"top"`
	Inner    DefaultOnlyInner // embedded struct
}

type DefaultOnlyNestedPtr struct {
	TopField string            `default:"top_ptr"`
	InnerPtr *DefaultOnlyInner // pointer to struct
}

type DefaultOnlyWithSlice struct {
	StringSlice []string  `default:"a,b,c"`
	IntSlice    []int     `default:"1,2,3"`
	FloatSlice  []float64 `default:"1.1,2.2,3.3"`
}

type DefaultOnlyDeepNested struct {
	Level1 string `default:"level1"`
	Inner  *struct {
		Level2     string `default:"level2"`
		DeepInner  *DefaultOnlyInner
		DeepInner2 DefaultOnlyInner
	}
}

type DefaultOnlyWithCustomFunc struct {
	CustomField string `default:"$(testDefaultFunc)"`
	NormalField string `default:"normal"`
}

type DefaultOnlyWithPtrToFundamental struct {
	StringPtr *string `default:"ptr_string"`
	IntPtr    *int    `default:"123"`
}

func TestGetDefaultOnlyStruct_Basic(t *testing.T) {
	input := &DefaultOnlyBasic{
		StringField: "existing_value", // should be ignored
		IntField:    999,              // should be ignored
	}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	output := result.(*DefaultOnlyBasic)

	// Verify defaults are applied
	assert.Equal(t, "hello", output.StringField)
	assert.Equal(t, 42, output.IntField)
	assert.Equal(t, uint(100), output.UintField)
	assert.Equal(t, 3.14, output.FloatField)
	assert.Equal(t, true, output.BoolField)
	assert.Equal(t, "", output.NoDefault) // no default tag
	assert.Equal(t, "", output.SkippedField) // skipped

	// Verify input is not modified
	assert.Equal(t, "existing_value", input.StringField)
	assert.Equal(t, 999, input.IntField)
}

func TestGetDefaultOnlyStruct_NestedEmbedded(t *testing.T) {
	input := &DefaultOnlyNested{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyNested)

	assert.Equal(t, "top", output.TopField)
	assert.Equal(t, "inner_default", output.Inner.InnerString)
	assert.Equal(t, 99, output.Inner.InnerInt)
}

func TestGetDefaultOnlyStruct_NestedPointer(t *testing.T) {
	input := &DefaultOnlyNestedPtr{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyNestedPtr)

	assert.Equal(t, "top_ptr", output.TopField)
	assert.NotNil(t, output.InnerPtr)
	assert.Equal(t, "inner_default", output.InnerPtr.InnerString)
	assert.Equal(t, 99, output.InnerPtr.InnerInt)
}

func TestGetDefaultOnlyStruct_Slices(t *testing.T) {
	input := &DefaultOnlyWithSlice{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyWithSlice)

	assert.Equal(t, []string{"a", "b", "c"}, output.StringSlice)
	assert.Equal(t, []int{1, 2, 3}, output.IntSlice)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, output.FloatSlice)
}

func TestGetDefaultOnlyStruct_DeepNested(t *testing.T) {
	input := &DefaultOnlyDeepNested{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyDeepNested)

	assert.Equal(t, "level1", output.Level1)
	assert.NotNil(t, output.Inner)
	assert.Equal(t, "level2", output.Inner.Level2)
	assert.NotNil(t, output.Inner.DeepInner)
	assert.Equal(t, "inner_default", output.Inner.DeepInner.InnerString)
	assert.Equal(t, 99, output.Inner.DeepInner.InnerInt)
	assert.Equal(t, "inner_default", output.Inner.DeepInner2.InnerString)
	assert.Equal(t, 99, output.Inner.DeepInner2.InnerInt)
}

func TestGetDefaultOnlyStruct_CustomFunc(t *testing.T) {
	// Register a custom default function
	RegisterDefaultFunc("testDefaultFunc", func(fieldname string) interface{} {
		return "custom_value_from_func"
	})
	defer func() {
		delete(defaultFuncs, "testDefaultFunc")
	}()

	input := &DefaultOnlyWithCustomFunc{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyWithCustomFunc)

	assert.Equal(t, "custom_value_from_func", output.CustomField)
	assert.Equal(t, "normal", output.NormalField)
}

func TestGetDefaultOnlyStruct_PtrToFundamental(t *testing.T) {
	input := &DefaultOnlyWithPtrToFundamental{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*DefaultOnlyWithPtrToFundamental)

	assert.NotNil(t, output.StringPtr)
	assert.Equal(t, "ptr_string", *output.StringPtr)
	assert.NotNil(t, output.IntPtr)
	assert.Equal(t, 123, *output.IntPtr)
}

func TestGetDefaultOnlyStruct_NotPointer(t *testing.T) {
	input := DefaultOnlyBasic{}

	result, err := GetDefaultOnlyStruct(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not a pointer")
}

func TestGetDefaultOnlyStruct_NotStruct(t *testing.T) {
	input := "not a struct"

	result, err := GetDefaultOnlyStruct(&input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not a struct")
}

func TestGetDefaultOnlyStruct_InputNotModified(t *testing.T) {
	input := &DefaultOnlyBasic{
		StringField: "original",
		IntField:    12345,
		BoolField:   false,
	}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify input values are unchanged
	assert.Equal(t, "original", input.StringField)
	assert.Equal(t, 12345, input.IntField)
	assert.Equal(t, false, input.BoolField)

	// Verify output has defaults
	output := result.(*DefaultOnlyBasic)
	assert.Equal(t, "hello", output.StringField)
	assert.Equal(t, 42, output.IntField)
	assert.Equal(t, true, output.BoolField)
}

func TestGetDefaultOnlyStruct_CustomFuncReturnsStruct(t *testing.T) {
	type StructWithCustomStruct struct {
		Inner *DefaultOnlyInner `default:"$(customStructFunc)"`
	}

	RegisterDefaultFunc("customStructFunc", func(fieldname string) interface{} {
		return &DefaultOnlyInner{
			InnerString: "from_custom_func",
			InnerInt:    777,
		}
	})
	defer func() {
		delete(defaultFuncs, "customStructFunc")
	}()

	input := &StructWithCustomStruct{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*StructWithCustomStruct)

	assert.NotNil(t, output.Inner)
	assert.Equal(t, "from_custom_func", output.Inner.InnerString)
	assert.Equal(t, 777, output.Inner.InnerInt)
}

func TestGetDefaultOnlyStruct_DefaultSkip(t *testing.T) {
	type StructWithDefaultSkip struct {
		Normal  string `default:"normal_value"`
		Skipped string `default:"skip_value" conf:"defaultskip"`
	}

	input := &StructWithDefaultSkip{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*StructWithDefaultSkip)

	assert.Equal(t, "normal_value", output.Normal)
	assert.Equal(t, "", output.Skipped) // defaultskip means no default applied
}

func TestGetDefaultOnlyStruct_BoolVariants(t *testing.T) {
	type StructWithBools struct {
		Bool1 bool `default:"true"`
		Bool2 bool `default:"false"`
		Bool3 bool `default:"1"`
		Bool4 bool `default:"0"`
		Bool5 bool `default:"yes"`
		Bool6 bool `default:"no"`
	}

	input := &StructWithBools{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*StructWithBools)

	assert.True(t, output.Bool1)
	assert.False(t, output.Bool2)
	assert.True(t, output.Bool3)
	assert.False(t, output.Bool4)
	assert.True(t, output.Bool5)
	assert.False(t, output.Bool6)
}

func TestGetDefaultOnlyStruct_IntVariants(t *testing.T) {
	type StructWithInts struct {
		Int8Field  int8  `default:"127"`
		Int16Field int16 `default:"32767"`
		Int32Field int32 `default:"2147483647"`
	}

	input := &StructWithInts{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*StructWithInts)

	assert.Equal(t, int8(127), output.Int8Field)
	assert.Equal(t, int16(32767), output.Int16Field)
	assert.Equal(t, int32(2147483647), output.Int32Field)
}

func TestGetDefaultOnlyStruct_UintVariants(t *testing.T) {
	type StructWithUints struct {
		Uint8Field  uint8  `default:"255"`
		Uint16Field uint16 `default:"65535"`
		Uint32Field uint32 `default:"4294967295"`
	}

	input := &StructWithUints{}

	result, err := GetDefaultOnlyStruct(input)
	assert.NoError(t, err)

	output := result.(*StructWithUints)

	assert.Equal(t, uint8(255), output.Uint8Field)
	assert.Equal(t, uint16(65535), output.Uint16Field)
	assert.Equal(t, uint32(4294967295), output.Uint32Field)
}

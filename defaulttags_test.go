package conftagz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFields(t *testing.T) {
	mystruct := MyStruct{"Value1", "", 33}

	expected := []string{"Field2"}

	result, err := SubsistuteDefaults(&mystruct, &DefaultFieldSubstOpts{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "Value1", mystruct.Field1)
	assert.Equal(t, "Banana", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 33)
}

// assumea PATH is there
func TestDefaultFields2(t *testing.T) {
	mystruct := MyStruct2{}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	debugf("result: %v\n", result)
	assert.Equal(t, 0, len(result))
}

func TestDefaultFields3(t *testing.T) {
	mystruct := MyStruct{}
	expected := []string{"Field1", "Field2", "Field3"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", result)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "Apple", mystruct.Field1)
	assert.Equal(t, "Banana", mystruct.Field2)
	assert.Equal(t, 999, mystruct.Field3)
}
func TestDefaultFieldsWithStruct(t *testing.T) {
	mystruct := MyStructWithStruct{"", "Value2", 0, "?", nil, nil, nil, InnerStruct{""}}

	expected := []string{"Field1", "Field3", "Field5", "Field6", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Apple")
	assert.Equal(t, mystruct.Field3, 999)
	assert.Equal(t, mystruct.Field2, "Value2")
	assert.Equal(t, mystruct.Field4, "?")
	assert.Equal(t, "Eggs", *mystruct.Field5)
	assert.Equal(t, int32(701), *mystruct.Field6)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "InnerApple", mystruct.Inner.FieldInner1)
}

func TestDefaultFieldsWithStruct2(t *testing.T) {
	mystruct := MyStructWithStruct2{"", "Value2", 0, "?", nil, nil, 0, nil, InnerStruct{""}, nil, nil, nil}

	expected := []string{"Field1", "Field3", "Field5", "Field6", "Inner.FieldInner1"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Apple")
	assert.Equal(t, mystruct.Field3, 999)
	assert.Equal(t, mystruct.Field2, "Value2")
	assert.Equal(t, mystruct.Field4, "?")
	assert.Equal(t, "Eggs", *mystruct.Field5)
	assert.Equal(t, int32(701), *mystruct.Field6)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "InnerApple", mystruct.Inner.FieldInner1)
	assert.Nil(t, mystruct.InnerPtr)
}

func TestDefaultFieldsWithStructStringExists(t *testing.T) {
	newstring := "NewString"
	mystruct := MyStructWithStruct{"", "Value2", 0, "?", &newstring, nil, nil, InnerStruct{""}}

	expected := []string{"Field1", "Field3", "Field6", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Apple")
	assert.Equal(t, mystruct.Field3, 999)
	assert.Equal(t, mystruct.Field2, "Value2")
	assert.Equal(t, mystruct.Field4, "?")
	assert.Equal(t, "NewString", *mystruct.Field5)
	assert.Equal(t, int32(701), *mystruct.Field6)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "InnerApple", mystruct.Inner.FieldInner1)
}

func TestDefaultFieldWithStruct3(t *testing.T) {
	mystruct := MyStruct3{"Value1", "", 3, nil, nil}

	expected := []string{"Field2", "Inner2Ptr.Stuff1"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Value1")
	assert.Equal(t, mystruct.Field2, "Banana")
	assert.Equal(t, mystruct.Field3, 3)
	assert.Equal(t, "InnerApple2", mystruct.Inner2Ptr.Stuff1)
	// test 'skip' tag
	assert.Nil(t, mystruct.InnerPtr)
}

func TestDefaultFieldsWithStruct3NonNilStructs(t *testing.T) {
	mystruct := MyStruct3{"", "", 0, &InnerStruct2{}, &InnerStruct2{}}

	expected := []string{"Field1", "Field2", "Field3", "Inner2Ptr.Stuff1"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "Apple", mystruct.Field1)
	assert.Equal(t, 999, mystruct.Field3)

	assert.Equal(t, "InnerApple2", mystruct.Inner2Ptr.Stuff1)
	// this was a skip
	assert.Equal(t, "", mystruct.InnerPtr.Stuff1)
}

func TestDefaultFieldsWithSlice(t *testing.T) {
	mystruct := MyStructWithSlice{"", nil}

	expected := []string{"Field1", "SliceField"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	expectedslice := []string{"Apple", "Banana"}

	assert.Equal(t, mystruct.Field1, "Apple")
	if !reflect.DeepEqual(mystruct.SliceField, expectedslice) {
		t.Errorf("Expected %v, but got %v", expectedslice, mystruct.SliceField)
	}

}

func TestDefaultCustomFunc(t *testing.T) {
	mystruct := AStructWithCustom{}

	expected := []string{"Field1", "Field2", "DefaultStruct"}

	field1func := func(fieldname string) interface{} {
		return "field1funcval"
	}

	field2func := func(fieldname string) interface{} {
		str := "field2funcval"
		return &str
	}

	fieldstructfund := func(fieldname string) interface{} {
		return &InnerStruct2{
			Stuff1: "specialsauce",
		}
	}

	RegisterDefaultFunc("field1default", field1func)
	RegisterDefaultFunc("field2default", field2func)
	RegisterDefaultFunc("fieldinnerstruct2", fieldstructfund)

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", result)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "field1funcval", mystruct.Field1)
	assert.Equal(t, "field2funcval", *mystruct.Field2)
	assert.Equal(t, "specialsauce", mystruct.DefaultStruct.Stuff1)
}

func TestDefaultSliceOfPointersToStruct(t *testing.T) {
	//	mystruct := MyStructWithSliceOfPointersToStruct{"APPPPLE", nil, []int{1, 2, 3}}
	mystruct := MyStructWithSliceOfPointersToStruct{"", nil, nil, nil, nil}

	inner1 := InnerStruct{""}
	inner2 := InnerStruct{""}
	inner1_2 := InnerStruct{""}

	slice := []*InnerStruct{&inner1, &inner2}
	mystruct.SliceField = slice

	slice2 := []InnerStruct{inner1_2}
	mystruct.SliceField2 = slice2

	innerstructcustomfunc := func(fieldname string) interface{} {
		return &InnerStruct{
			FieldInner1: "I123e",
		}
	}

	RegisterDefaultFunc("innerstructcustom", innerstructcustomfunc)

	expected := []string{"Field1", "SliceField[0].FieldInner1", "SliceField[1].FieldInner1", "SliceField2[0].FieldInner1", "SliceInts", "InnerStructCustom"}

	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Apple")
	assert.Equal(t, mystruct.SliceField[0].FieldInner1, "InnerApple")
	assert.Equal(t, mystruct.SliceField[1].FieldInner1, "InnerApple")
	assert.Equal(t, mystruct.SliceField2[0].FieldInner1, "InnerApple")
	assert.Equal(t, mystruct.InnerStructCustom.FieldInner1, "I123e")
}

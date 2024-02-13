package conftagz

import (
	"flag"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagFields(t *testing.T) {
	mystruct := MyStruct{"Value1", "", 33, 0, false}

	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-important", "Banana", "-veryimportant", "Razzles", "-extremelyimportant", "88", "-field4"}

	// expected := []string{"Field2"}

	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "Razzles", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 88)
	assert.Equal(t, true, mystruct.Field4)
}

func TestFlagFieldWithPrivateFieldTagShouldFail(t *testing.T) {

	mystruct := MyStructWithPrivateAndTag{"Value1", 0}
	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-important", "Banana", "-veryimportant", "Razzles", "-extremelyimportant", "88"}

	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	// assert.Equal(t, 0, len(result))
}

// assumea PATH is there
// func TestFlagFields2(t *testing.T) {
// 	mystruct := MyStruct2{}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	debugf("result: %v\n", result)
// 	assert.Equal(t, 0, len(result))
// }

// func TestFlagsFields3(t *testing.T) {
// 	mystruct := MyStruct{}
// 	expected := []string{"Field1", "Field2", "Field3"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", result)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

//		assert.Equal(t, "Apple", mystruct.Field1)
//		assert.Equal(t, "Banana", mystruct.Field2)
//		assert.Equal(t, 999, mystruct.Field3)
//	}
func TestFlagsFieldsWithStruct(t *testing.T) {
	mystruct := MyStructWithStruct{"", "Value2", 0, "?", nil, nil, nil, InnerStruct{"", nil, innerStruct{0}}}

	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-field1", "Ape", "-field3", "1024", "-field6", "8888", "-inner1", "Skynet"}

	//expected := []string{"Field1", "Field3", "Field6"} //, "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, mystruct.Field1, "Ape")
	assert.Equal(t, mystruct.Field3, 1024)
	assert.Equal(t, int32(8888), *mystruct.Field6)
	assert.Equal(t, "Skynet", mystruct.InnerPtr.FieldInner1)
	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
}

func TestFlagsFieldsWithStruct2(t *testing.T) {
	mystruct := MyStructWithStruct2{"", "Value2", 0, "?", nil, nil, 0, nil, InnerStruct{"", nil, innerStruct{0}}, nil, nil, nil}

	// expected := []string{"Field1", "Field3", "Field5", "Field6", "Inner.FieldInner1"}

	//	result, err := SubsistuteDefaults(&mystruct, nil)
	argz := []string{"-field3", "8080", "-field5", "Cayman", "-field6", "201", "-inner1", "Skynet"}
	err := ProcessFlagsWithFlagSet(&mystruct, flag.NewFlagSet("test", flag.ContinueOnError), argz)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, 8080, mystruct.Field3)
	assert.Equal(t, "Cayman", *mystruct.Field5)
	assert.Equal(t, int32(201), *mystruct.Field6)
	// assert.Equal(t, "Eggs", *mystruct.Field5)
	// assert.Equal(t, int32(701), *mystruct.Field6)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	//assert.Equal(t, "Skynet", mystruct.InnerPtr.FieldInner1)
	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
	assert.Nil(t, mystruct.InnerPtr)
}

func TestFlagsFieldsWithStructStringExists(t *testing.T) {
	newstring := "NewString"
	mystruct := MyStructWithStruct{"", "Value2", 0, "?", &newstring, nil, nil, InnerStruct{"", nil, innerStruct{0}}}

	//	expected := []string{"Field1", "Field3", "Field6", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-field1", "Ape", "-field3", "1024", "-field6", "8888", "-inner1", "Skynet"}

	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)
	//	result, err := SubsistuteDefaults(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Ape", mystruct.Field1)
	assert.Equal(t, 1024, mystruct.Field3)
	assert.Equal(t, int32(8888), *mystruct.Field6)
	assert.Equal(t, "NewString", *mystruct.Field5)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
}

func TestFlagsFieldWithStruct3(t *testing.T) {
	mystruct := MyStruct3{"Value1", "", 3, nil, nil, nil}

	// expected := []string{"Field2", "Inner2Ptr.Stuff1"}
	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-important", "Ape", "-field4", "-stuff1", "Giggy"}

	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Ape", mystruct.Field1)
	assert.True(t, *mystruct.Field4)
	assert.Equal(t, "Giggy", mystruct.Inner2Ptr.Stuff1)
	// test 'skip' tag
	assert.Nil(t, mystruct.InnerPtr)
}

func TestFlagsFieldsWithStruct3NonNilStructs(t *testing.T) {
	mystruct := MyStruct3{"", "", 0, &InnerStruct2{}, &InnerStruct2{}, nil}

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

func TestFlagsFieldsWithSlice(t *testing.T) {
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

func TestFlagsCustomFunc(t *testing.T) {
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

func TestFlagsSliceOfPointersToStruct(t *testing.T) {
	//	mystruct := MyStructWithSliceOfPointersToStruct{"APPPPLE", nil, []int{1, 2, 3}}
	mystruct := MyStructWithSliceOfPointersToStruct{"", nil, nil, nil, nil}

	inner1 := InnerStruct{"", nil, innerStruct{0}}
	inner2 := InnerStruct{"", nil, innerStruct{0}}
	inner1_2 := InnerStruct{"", nil, innerStruct{0}}

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

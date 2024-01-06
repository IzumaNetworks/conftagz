package conftagz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestFieldsFail(t *testing.T) {
	mystruct := MyStruct{"Value1", "Real Tomatoes", 33}

	_, err := RunTestFlags(&mystruct, nil)
	assert.EqualError(t, err, "value 33 ! >= 1024")
}
func TestTestFieldsFail2(t *testing.T) {
	mystruct := MyStruct{"Value1", "meh", 1025}

	_, err := RunTestFlags(&mystruct, nil)
	assert.EqualError(t, err, "value meh !~ regexp")
}

func TestTestFieldsPass(t *testing.T) {
	mystruct := MyStruct{"Value1", "Real Tomatoes", 1025}

	expected := []string{"Field2", "Field3"}

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, "Value1", mystruct.Field1)
	assert.Equal(t, "Real Tomatoes", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 1025)
}

// assumes PATH is there
// func TestTestFields2(t *testing.T) {
// 	mystruct := MyStruct2{}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	debugf("result: %v\n", result)
// 	assert.Equal(t, 0, len(result))
// }

// func TestTestFields3(t *testing.T) {
// 	mystruct := MyStruct{}
// 	expected := []string{"Field1", "Field2", "Field3"}

// 	result, err := SubsistuteDefaults(&mystruct, &EnvFieldSubstOpts{ThrowErrorIfEnvMissing: true})
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
func TestTestFieldsWithStruct(t *testing.T) {
	str := "Elastic"
	mystruct := MyStructWithStruct{"Appppple", "Bolivia", 1, "?", &str, nil, &InnerStruct{"123"}, InnerStruct{"LALA"}}

	expected := []string{"Field1", "Field2", "Field3", "Field5", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Appppple")
	assert.Equal(t, mystruct.Field3, 1)
	assert.Equal(t, mystruct.Field2, "Bolivia")
	assert.Equal(t, mystruct.Field4, "?")
	assert.Equal(t, "Elastic", *mystruct.Field5)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "123", mystruct.InnerPtr.FieldInner1)
	assert.Equal(t, "LALA", mystruct.Inner.FieldInner1)
}

func TestTestFieldsWithStruct2(t *testing.T) {
	str := "Elastic"
	mystruct := MyStructWithStruct2{"Appppple", "Bolivia", 1, "?", &str, nil, 0, nil, InnerStruct{"LALA"}, nil, nil, nil}

	expected := []string{"Field1", "Field2", "Field3", "Field5", "Inner.FieldInner1"}

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "Appppple")
	assert.Equal(t, mystruct.Field3, 1)
	assert.Equal(t, mystruct.Field2, "Bolivia")
	assert.Equal(t, mystruct.Field4, "?")
	assert.Equal(t, "Elastic", *mystruct.Field5)
	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
	assert.Equal(t, "LALA", mystruct.Inner.FieldInner1)
	assert.Nil(t, mystruct.InnerPtr)
}

func TestTestCustomFunc(t *testing.T) {
	ptrstuff := "ptrstuff"
	mystruct := AStructWithCustom{"stuff", &ptrstuff, &InnerStruct2{"innerstuff"}}
	expected := []string{"Field1", "Field2", "DefaultStruct"}

	var called1, called2, called3 bool

	field1func := func(val interface{}, fieldname string) bool {
		valstr, ok := val.(string)
		called1 = true
		if !ok {
			t.Errorf("Expected string, but got %v", val)
			return false
		}
		if valstr != "stuff" {
			return false
		}
		return true
	}

	field2func := func(val interface{}, fieldname string) bool {
		valstr, ok := val.(string)
		called2 = true
		if !ok {
			t.Errorf("Expected string, but got %v", val)
			return false
		}
		if valstr != "ptrstuff" {
			return false
		}
		return true
	}

	fieldstructfunc := func(val interface{}, fieldname string) bool {
		valstr, ok := val.(*InnerStruct2)
		called3 = true
		if !ok {
			t.Errorf("Expected struct, but got %v", val)
			return false
		}
		if valstr.Stuff1 != "innerstuff" {
			return false
		}
		return true
	}

	RegisterTestFunc("field1test", field1func)
	RegisterTestFunc("field2test", field2func)
	RegisterTestFunc("fieldinnerstruct2test", fieldstructfunc)

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	assert.True(t, called1)
	assert.True(t, called2)
	assert.True(t, called3)
	assert.Equal(t, "stuff", mystruct.Field1)
	assert.Equal(t, "ptrstuff", *mystruct.Field2)
	assert.Equal(t, "innerstuff", mystruct.DefaultStruct.Stuff1)
}

// func TestTestFieldsWithStructStringExists(t *testing.T) {
// 	newstring := "NewString"
// 	mystruct := MyStructWithStruct{"", "Value2", 0, "?", &newstring, nil, nil, InnerStruct{""}}

// 	expected := []string{"Field1", "Field3", "Field6", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	assert.Equal(t, mystruct.Field1, "Apple")
// 	assert.Equal(t, mystruct.Field3, 999)
// 	assert.Equal(t, mystruct.Field2, "Value2")
// 	assert.Equal(t, mystruct.Field4, "?")
// 	assert.Equal(t, "NewString", *mystruct.Field5)
// 	assert.Equal(t, int32(701), *mystruct.Field6)
// 	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
// 	assert.Equal(t, "InnerApple", mystruct.Inner.FieldInner1)
// }

// func TestTestFieldWithStruct3(t *testing.T) {
// 	mystruct := MyStruct3{"Value1", "", 3, nil, nil}

// 	expected := []string{"Field2", "Inner2Ptr.Stuff1"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	assert.Equal(t, mystruct.Field1, "Value1")
// 	assert.Equal(t, mystruct.Field2, "Banana")
// 	assert.Equal(t, mystruct.Field3, 3)
// 	assert.Equal(t, "InnerApple2", mystruct.Inner2Ptr.Stuff1)
// 	// test 'skip' tag
// 	assert.Nil(t, mystruct.InnerPtr)
// }

// func TestDefaultFieldsWithStruct3NonNilStructs(t *testing.T) {
// 	mystruct := MyStruct3{"", "", 0, &InnerStruct2{}, &InnerStruct2{}}

// 	expected := []string{"Field1", "Field2", "Field3", "Inner2Ptr.Stuff1"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	assert.Equal(t, "Apple", mystruct.Field1)
// 	assert.Equal(t, 999, mystruct.Field3)

// 	assert.Equal(t, "InnerApple2", mystruct.Inner2Ptr.Stuff1)
// 	// this was a skip
// 	assert.Equal(t, "", mystruct.InnerPtr.Stuff1)
// }

// func TestDefaultFieldsWithSlice(t *testing.T) {
// 	mystruct := MyStructWithSlice{"", nil}

// 	expected := []string{"Field1", "SliceField"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	expectedslice := []string{"Apple", "Banana"}

// 	assert.Equal(t, mystruct.Field1, "Apple")
// 	if !reflect.DeepEqual(mystruct.SliceField, expectedslice) {
// 		t.Errorf("Expected %v, but got %v", expectedslice, mystruct.SliceField)
// 	}

// }

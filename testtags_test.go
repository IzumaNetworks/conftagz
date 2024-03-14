package conftagz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// see envlookup_test.go for the struct definitions
func TestTestFieldsFail(t *testing.T) {
	mystruct := MyStruct{"Value1", "Real Tomatoes", 33, 0, 0, 0, false}

	_, err := RunTestFlags(&mystruct, nil)
	assert.EqualError(t, err, "field Field3: value 33 ! >= 1024")
}
func TestTestFieldsFail2(t *testing.T) {
	mystruct := MyStruct{"Value1", "meh", 1025, 0, 0, 0, false}

	_, err := RunTestFlags(&mystruct, nil)
	assert.EqualError(t, err, "field Field2: value \"meh\" !~ regexp R.*[Ss]{1}$")
}

func TestTestFieldWithPrivateFieldTagShouldFail(t *testing.T) {
	mystruct := MyStructWithPrivateAndTag{"Value1", 0}

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Expected no err got %v", result)
	}
}
func TestTestFieldsPass(t *testing.T) {
	mystruct := MyStruct{"Value1", "Real Tomatoes", 1025, 1025, 1025, 0, false}

	expected := []string{"Field2", "Field3", "Field3a", "Field3b"}

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

func TestTestFieldsWithStruct(t *testing.T) {
	str := "Elastic"
	mystruct := MyStructWithStruct{"Appppple", "Bolivia", 1, "?", &str, nil, &InnerStruct{"123", nil, innerStruct{0}}, InnerStruct{"LALA", nil, innerStruct{0}}}

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
	mystruct := MyStructWithStruct2{"Appppple", "Bolivia", 1, "?", &str, nil, 0, nil, InnerStruct{"LALA", nil, innerStruct{0}}, nil, nil, nil}

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

// and now for the fun stuff
func TestTestSliceOfPointersToStruct(t *testing.T) {
	mystruct := MyStructWithSliceOfPointersToStruct{"APPPPLE", nil, nil, []int{1, 2, 3}, &InnerStruct{"inner1", nil, innerStruct{0}}}

	inner1 := InnerStruct{"inner1", nil, innerStruct{0}}
	inner2 := InnerStruct{"inner2", nil, innerStruct{0}}
	inner1_2 := InnerStruct{"inner1_2", nil, innerStruct{0}}

	slice := []*InnerStruct{&inner1, &inner2}
	mystruct.SliceField = slice
	slice2 := []InnerStruct{inner1_2}
	mystruct.SliceField2 = slice2

	testslicefunc_ran := false

	testslicefunc := func(val interface{}, fieldname string) bool {
		valslice, ok := val.([]int)
		testslicefunc_ran = true
		if !ok {
			t.Errorf("Expected slice, but got %v", val)
			return false
		}
		if len(valslice) < 3 {
			return false
		}
		if !(valslice[2] > valslice[1] && valslice[1] > valslice[0]) {
			return false
		}
		return true
	}

	RegisterTestFunc("sliceintstest", testslicefunc)

	expected := []string{"Field1", "SliceField[0].FieldInner1", "SliceField[1].FieldInner1", "SliceField2[0].FieldInner1", "SliceInts", "InnerStructCustom.FieldInner1"}

	result, err := RunTestFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "APPPPLE")
	assert.Equal(t, mystruct.SliceField[0].FieldInner1, "inner1")
	assert.Equal(t, mystruct.SliceField[1].FieldInner1, "inner2")
	assert.True(t, testslicefunc_ran)

}

func TestTestSliceIntsFalseCustom(t *testing.T) {
	mystruct := MyStructWithSliceOfPointersToStruct{"APPPPLE", nil, nil, []int{1, 2, 3}, nil}

	testslicefunc_ran := false

	testslicefunc := func(val interface{}, fieldname string) bool {
		valslice, ok := val.([]int)
		testslicefunc_ran = true
		if !ok {
			t.Errorf("Expected slice, but got %v", val)
			return false
		}
		if len(valslice) < 3 {
			return false
		}
		if !(valslice[2] > valslice[1] && valslice[1] > valslice[0]) {
			return false
		}
		// returns false no matter...
		return false
	}

	RegisterTestFunc("sliceintstest", testslicefunc)

	expected := []string{"Field1", "SliceInts"}

	result, err := RunTestFlags(&mystruct, nil)
	var errstr string
	if err != nil {
		//		t.Errorf("Unexpected error: %v", err)
		errstr = err.Error()
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "APPPPLE")
	assert.Regexp(t, `value for field SliceInts.*sliceintstest`, errstr)
	assert.True(t, testslicefunc_ran)

}

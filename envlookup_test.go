package conftagz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	Field1 string `yaml:"important" env:"Important" default:"Apple"`
	Field2 string `json:"field2" env:"VeryImportant" default:"Banana" test:"~R.*[Ss]{1}$"`
	Field3 int    `env:"ExtremelyImportant" default:"999" test:">=1024"`
}
type MyStruct2 struct {
	Path string `yaml:"path" env:"PATH"`
}
type MyStructWithStruct struct {
	Field1   string `yaml:"important" env:"ENV1" default:"Apple" test:"~A.*[Ee]{1}"`
	Field2   string `json:"field2" env:"ENV2" default:"Banana" test:"~B.*a"`
	Field3   int    `env:"ENV3" default:"999" test:"<65537,>0"`
	Field4   string
	Field5   *string      `env:"ENV5" default:"Eggs" test:"~E.*"`
	Field6   *int32       `env:"ENV6" default:"701"`
	InnerPtr *InnerStruct `yaml:"inner"`
	Inner    InnerStruct  `yaml:"inner2"`
}

type UselessInterface interface {
	DoNothing()
}

type MyStructWithStruct2 struct {
	Field1          string `yaml:"important" env:"ENV1" default:"Apple" test:"~A.*[Ee]{1}"`
	Field2          string `json:"field2" env:"ENV2" default:"Banana" test:"~B.*a"`
	Field3          int    `env:"ENV3" default:"999" test:"<65537"`
	Field4          string
	Field5          *string      `env:"ENV5" default:"Eggs" test:"~E.*"`
	Field6          *int32       `env:"ENV6" default:"701"`
	Field7          int32        `test:">100" conf:"skipzero"`
	InnerPtr        *InnerStruct `yaml:"inner" conf:"skipnil"`
	Inner           InnerStruct  `yaml:"inner2"`
	RandomInterface UselessInterface
	// throw in some unusual, unsupported types
	WeirdStuff interface{}
	Weird2     *uintptr
}
type InnerStruct struct {
	FieldInner1 string `yaml:"important" env:"INNER1" default:"InnerApple" test:"~.{3,}"`
}
type MyStruct3 struct {
	Field1    string        `yaml:"important" env:"Important" default:"Apple"`
	Field2    string        `json:"field2" env:"VeryImportant" default:"Banana"`
	Field3    int           `env:"ExtremelyImportant" default:"999"`
	InnerPtr  *InnerStruct2 `yaml:"inner" conf:"skip"`
	Inner2Ptr *InnerStruct2 `yaml:"inner"`
}
type InnerStruct2 struct {
	Stuff1 string `yaml:"important" env:"STUFF1" default:"InnerApple2"`
}

type MyStructWithSlice struct {
	Field1     string   `yaml:"important" env:"Important" default:"Apple"`
	SliceField []string `yaml:"slice" default:"Apple,Banana"` // env:"SLICE"
}

type AStructWithCustom struct {
	Field1        string        `yaml:"field1" env:"FIELD1" default:"$(field1default)" test:"$(field1test)"`
	Field2        *string       `yaml:"field2" env:"FIELD2" default:"$(field2default)" test:"$(field2test)"`
	DefaultStruct *InnerStruct2 `yaml:"inner" default:"$(fieldinnerstruct2)" test:"$(fieldinnerstruct2test)"`
}

type MyStructWithSliceOfPointersToStruct struct {
	Field1            string         `yaml:"important" env:"Important" default:"Apple" test:"~A.*[Ee]{1}"`
	SliceField        []*InnerStruct `yaml:"slice"`
	SliceField2       []InnerStruct  `yaml:"slice"`
	SliceInts         []int          `yaml:"sliceints" default:"1,2,3" test:"$(sliceintstest)"`
	InnerStructCustom *InnerStruct   `yaml:"innerstructcustom" default:"$(innerstructcustom)"`
}

func TestEnvFieldSubstitutionFromMap(t *testing.T) {
	mystruct := MyStruct{"Value1", "Value2", 3}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"VeryImportant":      "NewValue2",
		"ExtremelyImportant": "123",
	}

	expected := []string{"Field1", "Field2", "Field3"}

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
	mystruct := MyStruct{"Value1", "Value2", 3}

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
	mystruct := MyStructWithStruct{"Value1", "Value2", 3, "", nil, nil, nil, InnerStruct{"InnerValue1"}}

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
	mystruct := MyStruct3{"Value1", "Value2", 3, nil, nil}

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
	mystruct := MyStruct3{"Value1", "Value2", 3, &InnerStruct2{"blah"}, &InnerStruct2{"blah"}}

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

// Slice is not supported for env vars. so should just ignore slice field
func TestEnvFieldsWithSlice(t *testing.T) {
	mystruct := MyStructWithSlice{"", nil}

	envMap := map[string]string{
		"Important":          "NewValue1",
		"VeryImportant":      "NewValue2",
		"ExtremelyImportant": "123",
		"STUFF1":             "InnerValue1",
	}

	expected := []string{"Field1"}

	result, err := EnvFieldSubstitutionFromMap(&mystruct, nil, envMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	assert.Equal(t, mystruct.Field1, "NewValue1")

}

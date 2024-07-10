package conftagz

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// each has a cobra varation
type MyStructCobra struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple" cflag:"important" cobra:"root"`
	Field2       string `json:"field2" env:"VeryImportant" default:"Banana" test:"~R.*[Ss]{1}$" cflag:"veryimportant" cobra:"root"`
	Field3       int    `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant" cobra:"root"`
	Field3a      uint   `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant_a" cobra:"root"`
	Field3b      uint64 `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant_b" cobra:"root"`
	privateField int
	Field4       bool `env:"Field4" cflag:"field4" cobra:"root"`
}

func TestCobraFields(t *testing.T) {
	ResetGlobals()
	mystruct := MyStructCobra{"Value1", "", 33, 0, 0, 0, false}

	//	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"--important", "Banana", "--veryimportant", "Razzles", "--extremelyimportant", "88", "--field4"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}

	RegisterCobraCmd("root", rootCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running root command %+v\n", args)
		return nil
	}
	err := PreProcessCobraFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error in PreProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.ParseFlags(argz)
	if err != nil {
		t.Errorf("Unexpected error in ParseFlagas: %v", err)
		return
	}
	err = PostProcessCobraFlags()
	if err != nil {
		t.Errorf("Unexpected error in PostProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error in Execute: %v", err)
		return
	}

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "Razzles", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 88)
	assert.Equal(t, mystruct.Field3a, uint(0))
	assert.Equal(t, mystruct.Field3b, uint64(0))
	assert.Equal(t, true, mystruct.Field4)
}
func TestCobraFieldsNoFlagProvided(t *testing.T) {
	ResetGlobals()
	mystruct := MyStructCobra{"Value1", "", 33, 0, 0, 0, false}

	argz := []string{"--important", "Banana", "--extremelyimportant", "88", "--field4"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}

	RegisterCobraCmd("root", rootCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running root command %+v\n", args)
		return nil
	}
	err := PreProcessCobraFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error PreProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.ParseFlags(argz)
	if err != nil {
		t.Errorf("Unexpected error ParseFlags(argz): %v", err)
		return
	}
	err = PostProcessCobraFlags()
	if err != nil {
		t.Errorf("Unexpected error PostProcessCobraFlags(): %v", err)
		return
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error in Execute(): %v", err)
		return
	}

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 88)
	assert.Equal(t, true, mystruct.Field4)
}

type MyStructCobra2 struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple" cflag:"important" cobra:"root"`
	Field2       string `json:"field2" env:"VeryImportant" default:"Banana" test:"~R.*[Ss]{1}$" cflag:"veryimportant" cobra:"root"`
	Verbose      bool   `env:"VERBOSE" default:"false" cflag:"verbose,v" cobra:"root,persistent"`
	Field3       int    `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant" cobra:"secondcmd"`
	Field3a      uint   `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant_a" cobra:"secondcmd"`
	Field3b      uint64 `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant_b" cobra:"secondcmd"`
	privateField int
	Field4       bool `env:"Field4" cflag:"field4" cobra:"root,persistent"`
}

func TestCobraFieldsPersistent(t *testing.T) {
	ResetGlobals()
	mystruct := MyStructCobra2{"Value1", "", false, 33, 0, 0, 0, false}

	argz := []string{"secondcmd", "-v", "--extremelyimportant", "88", "--field4"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}

	var secondCmd = &cobra.Command{
		Use:   "secondcmd",
		Short: "A simple CLI second command",
	}

	RegisterCobraCmd("root", rootCmd)
	RegisterCobraCmd("secondcmd", secondCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running root command %+v\n", args)
		return nil
	}
	secondCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running second command %+v\n", args)
		return nil
	}

	err := PreProcessCobraFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	rootCmd.AddCommand(secondCmd)
	rootCmd.SetArgs(argz)
	err = rootCmd.ParseFlags(argz)
	if err != nil {
		t.Errorf("Unexpected error in ParseFlags(argz = %v): %v", argz, err)
		return
	}
	secondCmd.SetArgs(argz)
	err = secondCmd.ParseFlags(argz)
	if err != nil {
		t.Errorf("Unexpected error in ParseFlags: %v", err)
		return
	}
	//	secondCmd.ParseFlags(argz)
	err = PostProcessCobraFlags()
	if err != nil {
		t.Errorf("Unexpected error in PostProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error in Execute: %v", err)
		return
	}
	assert.Equal(t, "Value1", mystruct.Field1)
	assert.Equal(t, true, mystruct.Verbose)
	assert.Equal(t, "", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 88)
	assert.Equal(t, mystruct.Field3a, uint(0))
	assert.Equal(t, mystruct.Field3b, uint64(0))
	assert.Equal(t, true, mystruct.Field4)
}

type MyStructWithPrivateAndTagCobra struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple"`
	privateField int    `env:"WONTWORK" default:"123" test:">=1024" cobra:"root" cflag:"important"`
}

func TestCobraFieldWithPrivateFieldTagShouldFail(t *testing.T) {
	ResetGlobals()

	mystruct := MyStructWithPrivateAndTagCobra{"Value1", 0}
	argz := []string{"--important", "Banana", "--veryimportant", "Razzles", "--extremelyimportant", "88"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}
	err := PreProcessCobraFlags(&mystruct, nil)
	if err != nil {
		t.Logf("Expected error: %v", err)
	} else {
		t.Errorf("Should have had error.")
	}
	rootCmd.ParseFlags(argz)
	err = PostProcessCobraFlags()
	if err != nil {
		t.Errorf("Unexpected error in PostProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error in Execute: %v", err)
		return
	}
	// assert.Equal(t, 0, len(result))
}

// // assumea PATH is there
// // func TestCobraFields2(t *testing.T) {
// // 	mystruct := MyStruct2{}

// // 	result, err := SubsistuteDefaults(&mystruct, nil)
// // 	if err != nil {
// // 		t.Errorf("Unexpected error: %v", err)
// // 	}
// // 	debugf("result: %v\n", result)
// // 	assert.Equal(t, 0, len(result))
// // }

// // func TestCobrasFields3(t *testing.T) {
// // 	mystruct := MyStruct{}
// // 	expected := []string{"Field1", "Field2", "Field3"}

// // 	result, err := SubsistuteDefaults(&mystruct, nil)
// // 	if err != nil {
// // 		t.Errorf("Unexpected error: %v", result)
// // 	}

// // 	if !reflect.DeepEqual(result, expected) {
// // 		t.Errorf("Expected %v, but got %v", expected, result)
// // 	}

// //		assert.Equal(t, "Apple", mystruct.Field1)
// //		assert.Equal(t, "Banana", mystruct.Field2)
// //		assert.Equal(t, 999, mystruct.Field3)
// //	}

type MyStructWithStructCobra struct {
	Field1   string `yaml:"important" env:"ENV1" default:"Apple" test:"~A.*[Ee]{1}" cflag:"field1" usage:"Usage for field1" cobra:"root"`
	Field2   string `json:"field2" env:"ENV2" default:"Banana" test:"~B.*a"`
	Field3   int    `env:"ENV3" default:"999" test:"<65537,>0" cflag:"field3" cobra:"root" usage:"Usage for field3"`
	Field4   string
	Field5   *string           `env:"ENV5" default:"Eggs" test:"~E.*"`
	Field6   *int32            `env:"ENV6" default:"701" cflag:"field6" usage:"Usage for field6" cobra:"root"`
	InnerPtr *InnerStructCobra `yaml:"inner"`
	Inner    InnerStructCobra  `yaml:"inner2"`
}

type InnerStructCobra struct {
	FieldInner1   string `yaml:"important" env:"INNER1" default:"InnerApple" test:"~.{3,}" cflag:"inner1" usage:"inner1 usage" cobra:"root"`
	privateField  *innerStruct
	privateField2 innerStruct
}

func TestCobrasFieldsWithStruct(t *testing.T) {
	mystruct := MyStructWithStructCobra{"", "Value2", 0, "?", nil, nil, nil, InnerStructCobra{"", nil, innerStruct{0}}}

	argz := []string{"--field1", "Ape", "--field3", "1024", "--field6", "8888", "--inner1", "Skynet"}

	//expected := []string{"Field1", "Field3", "Field6"} //, "InnerPtr.FieldInner1", "Inner.FieldInner1"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}

	RegisterCobraCmd("root", rootCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running root command %+v\n", args)

		return nil
	}
	err := PreProcessCobraFlags(&mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	err = rootCmd.ParseFlags(argz)
	if err != nil {
		t.Errorf("Unexpected error in ParseFlags: %v", err)
		return
	}
	err = PostProcessCobraFlags()
	if err != nil {
		t.Errorf("Unexpected error in PostProcessCobraFlags: %v", err)
		return
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error in Execute: %v", err)
		return
	}
	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, mystruct.Field1, "Ape")
	assert.Equal(t, mystruct.Field3, 1024)
	assert.Equal(t, int32(8888), *mystruct.Field6)
	//	assert.Equal(t, "Skynet", mystruct.InnerPtr.FieldInner1)
	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
}

// func TestCobrasFieldsWithStruct2(t *testing.T) {
// 	mystruct := MyStructWithStruct2{"", "Value2", 0, "?", nil, nil, 0, nil, InnerStruct{"", nil, innerStruct{0}}, nil, nil, nil}

// 	// expected := []string{"Field1", "Field3", "Field5", "Field6", "Inner.FieldInner1"}

// 	//	result, err := SubsistuteDefaults(&mystruct, nil)
// 	argz := []string{"-field3", "8080", "-field5", "Cayman", "-field6", "201", "-inner1", "Skynet"}
// 	err := ProcessFlagsWithFlagSet(&mystruct, flag.NewFlagSet("test", flag.ContinueOnError), argz)

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	// if !reflect.DeepEqual(result, expected) {
// 	// 	t.Errorf("Expected %v, but got %v", expected, result)
// 	// }

// 	assert.Equal(t, 8080, mystruct.Field3)
// 	assert.Equal(t, "Cayman", *mystruct.Field5)
// 	assert.Equal(t, int32(201), *mystruct.Field6)
// 	// assert.Equal(t, "Eggs", *mystruct.Field5)
// 	// assert.Equal(t, int32(701), *mystruct.Field6)
// 	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
// 	//assert.Equal(t, "Skynet", mystruct.InnerPtr.FieldInner1)
// 	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
// 	assert.Nil(t, mystruct.InnerPtr)
// }

// func TestCobrasFieldsWithStructStringExists(t *testing.T) {
// 	newstring := "NewString"
// 	mystruct := MyStructWithStruct{"", "Value2", 0, "?", &newstring, nil, nil, InnerStruct{"", nil, innerStruct{0}}}

// 	//	expected := []string{"Field1", "Field3", "Field6", "InnerPtr.FieldInner1", "Inner.FieldInner1"}

// 	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
// 	argz := []string{"-field1", "Ape", "-field3", "1024", "-field6", "8888", "-inner1", "Skynet"}

// 	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)
// 	//	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	// if !reflect.DeepEqual(result, expected) {
// 	// 	t.Errorf("Expected %v, but got %v", expected, result)
// 	// }

// 	assert.Equal(t, "Ape", mystruct.Field1)
// 	assert.Equal(t, 1024, mystruct.Field3)
// 	assert.Equal(t, int32(8888), *mystruct.Field6)
// 	assert.Equal(t, "NewString", *mystruct.Field5)
// 	//	assert.Equal(t, mystruct.InnerPtr.FieldInner1, "InnerApple")
// 	assert.Equal(t, "Skynet", mystruct.Inner.FieldInner1)
// }

// func TestCobrasFieldWithStruct3(t *testing.T) {
// 	mystruct := MyStruct3{"Value1", "", 3, nil, nil, nil}

// 	// expected := []string{"Field2", "Inner2Ptr.Stuff1"}
// 	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
// 	argz := []string{"-important", "Ape", "-field4", "-stuff1", "Giggy"}

// 	err := ProcessFlagsWithFlagSet(&mystruct, flagset, argz)

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	// if !reflect.DeepEqual(result, expected) {
// 	// 	t.Errorf("Expected %v, but got %v", expected, result)
// 	// }

// 	assert.Equal(t, "Ape", mystruct.Field1)
// 	assert.True(t, *mystruct.Field4)
// 	assert.Equal(t, "Giggy", mystruct.Inner2Ptr.Stuff1)
// 	// test 'skip' tag
// 	assert.Nil(t, mystruct.InnerPtr)
// }

// func TestCobrasFieldsWithStruct3NonNilStructs(t *testing.T) {
// 	mystruct := MyStruct3{"", "", 0, &InnerStruct2{}, &InnerStruct2{}, nil}

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

// func TestCobrasFieldsWithSlice(t *testing.T) {
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

// func TestCobrasCustomFunc(t *testing.T) {
// 	mystruct := AStructWithCustom{}

// 	expected := []string{"Field1", "Field2", "DefaultStruct"}

// 	field1func := func(fieldname string) interface{} {
// 		return "field1funcval"
// 	}

// 	field2func := func(fieldname string) interface{} {
// 		str := "field2funcval"
// 		return &str
// 	}

// 	fieldstructfund := func(fieldname string) interface{} {
// 		return &InnerStruct2{
// 			Stuff1: "specialsauce",
// 		}
// 	}

// 	RegisterDefaultFunc("field1default", field1func)
// 	RegisterDefaultFunc("field2default", field2func)
// 	RegisterDefaultFunc("fieldinnerstruct2", fieldstructfund)

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", result)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	assert.Equal(t, "field1funcval", mystruct.Field1)
// 	assert.Equal(t, "field2funcval", *mystruct.Field2)
// 	assert.Equal(t, "specialsauce", mystruct.DefaultStruct.Stuff1)
// }

// func TestCobrasSliceOfPointersToStruct(t *testing.T) {
// 	//	mystruct := MyStructWithSliceOfPointersToStruct{"APPPPLE", nil, []int{1, 2, 3}}
// 	mystruct := MyStructWithSliceOfPointersToStruct{"", nil, nil, nil, nil}

// 	inner1 := InnerStruct{"", nil, innerStruct{0}}
// 	inner2 := InnerStruct{"", nil, innerStruct{0}}
// 	inner1_2 := InnerStruct{"", nil, innerStruct{0}}

// 	slice := []*InnerStruct{&inner1, &inner2}
// 	mystruct.SliceField = slice

// 	slice2 := []InnerStruct{inner1_2}
// 	mystruct.SliceField2 = slice2

// 	innerstructcustomfunc := func(fieldname string) interface{} {
// 		return &InnerStruct{
// 			FieldInner1: "I123e",
// 		}
// 	}

// 	RegisterDefaultFunc("innerstructcustom", innerstructcustomfunc)

// 	expected := []string{"Field1", "SliceField[0].FieldInner1", "SliceField[1].FieldInner1", "SliceField2[0].FieldInner1", "SliceInts", "InnerStructCustom"}

// 	result, err := SubsistuteDefaults(&mystruct, nil)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %v, but got %v", expected, result)
// 	}

// 	assert.Equal(t, mystruct.Field1, "Apple")
// 	assert.Equal(t, mystruct.SliceField[0].FieldInner1, "InnerApple")
// 	assert.Equal(t, mystruct.SliceField[1].FieldInner1, "InnerApple")
// 	assert.Equal(t, mystruct.SliceField2[0].FieldInner1, "InnerApple")
// 	assert.Equal(t, mystruct.InnerStructCustom.FieldInner1, "I123e")
// }

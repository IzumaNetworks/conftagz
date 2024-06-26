package conftagz

import (
	"flag"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type ConfStruct struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple" flag:"important"`
	Field2       string `json:"field2" env:"VeryImportant" default:"Randoms" test:"~R.*[Ss]{1}$" flag:"veryimportant"`
	Field3       int    `env:"ExtremelyImportant" default:"999" test:">=1024" flag:"extremelyimportant"`
	privateField int
	Field4       bool `env:"Field4" flag:"field4"`
}

func TestProcessSelfParseFlags(t *testing.T) {

	mystruct := &ConfStruct{"Value1", "", 1111, 0, false}
	// assume conf file already read
	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-important", "Banana", "-veryimportant", "Razzles", "-extremelyimportant", "8888", "-field4"}

	flagtagopts := &FlagFieldSubstOpts{
		UseFlags: flagset,
		//		Args:     argz,
	}
	processed, err := ProcessFlagTags(mystruct, flagtagopts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	flagtagopts.Tags = processed
	flagset.Parse(argz)

	// expected := []string{"Field2"}

	err = Process(&ConfTagOpts{
		FlagTagOpts: flagtagopts,
	}, mystruct)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "Razzles", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 8888)
	assert.Equal(t, true, mystruct.Field4)
}

func TestProcessSelfParseFlags2(t *testing.T) {

	mystruct := &ConfStruct{"Value1", "", 1111, 0, false}
	// assume conf file already read
	flagset := flag.NewFlagSet("test", flag.ContinueOnError)
	argz := []string{"-important", "Banana", "-extremelyimportant", "8888"}

	flagtagopts := &FlagFieldSubstOpts{
		UseFlags: flagset,
		//		Args:     argz,
	}
	processed, err := ProcessFlagTags(mystruct, flagtagopts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	flagset.Parse(argz)
	flagtagopts.Tags = processed
	// expected := []string{"Field2"}

	err = Process(&ConfTagOpts{
		FlagTagOpts: flagtagopts,
	}, mystruct)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "Randoms", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 8888)
	assert.Equal(t, false, mystruct.Field4)
}

// cobra version of above
type ConfStructCobra struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple" cflag:"important" cobra:"root"`
	Field2       string `json:"field2" env:"VeryImportant" default:"Randoms" test:"~R.*[Ss]{1}$" cflag:"veryimportant" cobra:"root"`
	Field3       int    `env:"ExtremelyImportant" default:"999" test:">=1024" cflag:"extremelyimportant" cobra:"root"`
	privateField int
	Field4       bool `env:"Field4" cflag:"field4" cobra:"root"`
}

func TestProcessCobraTags(t *testing.T) {

	mystruct := &ConfStructCobra{"Value1", "", 1111, 0, false}
	// assume conf file already read
	argz := []string{"--important", "Banana", "--extremelyimportant", "8888"}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
	}

	RegisterCobraCmd("root", rootCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running root command %+v\n", args)
		return nil
	}
	err := PreProcessCobraFlags(mystruct, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	rootCmd.ParseFlags(argz)
	//	PostProcessCobraFlags()

	err = Process(&ConfTagOpts{
		//	FlagTagOpts: flagtagopts,
	}, mystruct)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	rootCmd.Execute()
	// if !reflect.DeepEqual(result, expected) {
	// 	t.Errorf("Expected %v, but got %v", expected, result)
	// }

	assert.Equal(t, "Banana", mystruct.Field1)
	assert.Equal(t, "Randoms", mystruct.Field2)
	assert.Equal(t, mystruct.Field3, 8888)
	assert.Equal(t, false, mystruct.Field4)
}

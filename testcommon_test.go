package conftagz

type MyStruct struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple" flag:"important"`
	Field2       string `json:"field2" env:"VeryImportant" default:"Banana" test:"~R.*[Ss]{1}$" flag:"veryimportant"`
	Field3       int    `env:"ExtremelyImportant" default:"999" test:">=1024" flag:"extremelyimportant"`
	Field3a      uint   `env:"ExtremelyImportant" default:"999" test:">=1024" flag:"extremelyimportant_a"`
	Field3b      uint64 `env:"ExtremelyImportant" default:"999" test:">=1024" flag:"extremelyimportant_b"`
	privateField int
	Field4       bool `env:"Field4" flag:"field4"`
}

type MyStructWithPrivateAndTag struct {
	Field1       string `yaml:"important" env:"Important" default:"Apple"`
	privateField int    `env:"WONTWORK" default:"123" test:">=1024"`
}

type MyStruct2 struct {
	Path string `yaml:"path" env:"PATH"`
}
type MyStructWithStruct struct {
	Field1   string `yaml:"important" env:"ENV1" default:"Apple" test:"~A.*[Ee]{1}" flag:"field1" usage:"Usage for field1"`
	Field2   string `json:"field2" env:"ENV2" default:"Banana" test:"~B.*a"`
	Field3   int    `env:"ENV3" default:"999" test:"<65537,>0" flag:"field3" usage:"Usage for field3"`
	Field4   string
	Field5   *string      `env:"ENV5" default:"Eggs" test:"~E.*"`
	Field6   *int32       `env:"ENV6" default:"701" flag:"field6" usage:"Usage for field6"`
	InnerPtr *InnerStruct `yaml:"inner"`
	Inner    InnerStruct  `yaml:"inner2"`
}

type UselessInterface interface {
	DoNothing()
}

type MyStructWithStruct2 struct {
	Field1          string `yaml:"important" env:"ENV1" default:"Apple" test:"~A.*[Ee]{1}"`
	Field2          string `json:"field2" env:"ENV2" default:"Banana" test:"~B.*a"`
	Field3          int    `env:"ENV3" default:"999" test:"<65537" flag:"field3" usage:"Usage for field3"`
	Field4          string
	Field5          *string      `env:"ENV5" default:"Eggs" test:"~E.*" flag:"field5" usage:"Usage for field5"`
	Field6          *int32       `env:"ENV6" default:"701" flag:"field6" usage:"Usage for field6"`
	Field7          int32        `test:">100" conf:"skipzero"`
	InnerPtr        *InnerStruct `yaml:"inner" conf:"skipnil"`
	Inner           InnerStruct  `yaml:"inner2"`
	RandomInterface UselessInterface
	// throw in some unusual, unsupported types
	WeirdStuff interface{}
	Weird2     *uintptr
}

type innerStruct struct {
	nothing int
}

type InnerStruct struct {
	FieldInner1   string `yaml:"important" env:"INNER1" default:"InnerApple" test:"~.{3,}" flag:"inner1" usage:"inner1 usage"`
	privateField  *innerStruct
	privateField2 innerStruct
}
type MyStruct3 struct {
	Field1    string        `yaml:"important" env:"Important" default:"Apple" flag:"important"`
	Field2    string        `json:"field2" env:"VeryImportant" default:"Banana"`
	Field3    int           `env:"ExtremelyImportant" default:"999"`
	InnerPtr  *InnerStruct2 `yaml:"inner" conf:"skip"`
	Inner2Ptr *InnerStruct2 `yaml:"inner"`
	Field4    *bool         `env:"Field4" flag:"field4"`
}
type InnerStruct2 struct {
	Stuff1 string `yaml:"important" env:"STUFF1" default:"InnerApple2" flag:"stuff1"`
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

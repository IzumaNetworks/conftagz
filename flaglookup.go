package conftagz

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

const FLAGFIELD = "flag"
const FLAGFIELDUSAGE = "usage"

// func FlagToMap() map[string]string {
// 	flagMap := make(map[string]string)

// 	// for _, env := range os.Environ() {
// 	// 	splitEnv := strings.SplitN(env, "=", 2)
// 	// 	envMap[splitEnv[0]] = splitEnv[1]
// 	// }
// 	return envMap
// }

// EnvFieldSubstitution is a function that takes a pointer to a struct
// and looks at each field. If the field has a ENVFIELD tag ("env" by default)
// then it will look up the value of the field in the environment variables
// and replace the field with the value.
// It returns a list of the names of the fields that were substituted - as
// a list of string
// If there is an error, it returns an error
// func FlagFieldSubstitution(somestruct interface{}, opts *EnvFieldSubstOpts) (ret []string, err error) {
// //	m := FlagToMap()
// 	return FlagFieldSubstitutionFromMap(somestruct, opts)
// }

// func StringToInt64(s string) (int64, error) {
// 	i, err := strconv.ParseInt(s, 10, 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}
// 	return i, nil
// }

// func StringToFloat64(s string) (float64, error) {
// 	i, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}
// 	return i, nil
// }

// func addParentPath(parentpath string, fieldname string) string {
// 	if len(parentpath) > 0 {
// 		return parentpath + "." + fieldname
// 	}
// 	return fieldname
// }

// this is used to set the flag before calling flag.Parse to get all the cmd line options to parse
type flagRetrieverFunc func(flagname string, self *flagSetRetriever) error

func (r *flagSetRetriever) retrieve(flagname string) (err error) {
	for _, retriever := range r.retrievers {
		err = retriever(flagname, r)
		if err != nil {
			return
		}
	}
	return
}

type flagSetRetriever struct {
	fieldName  string
	fieldValue reflect.Value
	retrievers []flagRetrieverFunc
	val        interface{}
	touched    bool
}

type ProcessedFlagTags struct {
	needflags map[string]*flagSetRetriever
	// true if we have ran ProcessAllFlagTags
	flagsProcessed bool
	fieldsTouched  []string
}

func (p *ProcessedFlagTags) GetFlagsFound() (ret []string) {
	ret = make([]string, 0)
	for k := range p.needflags {
		ret = append(ret, k)
	}
	return
}

func (p *ProcessedFlagTags) GetFieldsTouched() (ret []string) {
	return p.fieldsTouched
}

type FlagFieldSubstOpts struct {
	// throws an error if the environment variable is not found
	//	ThrowErrorIfEnvMissing bool
	UseFlags *flag.FlagSet
	Args     []string
	Tags     *ProcessedFlagTags
}

func ProcessFlagTags(somestruct interface{}, opts *FlagFieldSubstOpts) (ret *ProcessedFlagTags, err error) {
	ret = &ProcessedFlagTags{}
	if opts == nil {
		opts = &FlagFieldSubstOpts{}

	}
	if opts.Tags == nil {
		opts.Tags = ret
	}
	if opts.Tags.flagsProcessed {
		return opts.Tags, nil
	}
	ret.flagsProcessed = true
	// this will record all the flags we need to set
	ret.needflags = make(map[string]*flagSetRetriever)

	if opts == nil {
		opts = &FlagFieldSubstOpts{}
	}
	myflags := opts.UseFlags

	if myflags == nil {
		myflags = flag.CommandLine
	}

	setFlagVal := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string, usagetag string, existing *flagSetRetriever) (retriever *flagSetRetriever, err error) {
		k := fieldValue.Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.Bool:
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					s, ok := r.val.(*bool)
					if ok {
						fieldValue.SetBool(*s)
					}
				}
				// else {
				// 	return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
				// }
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				// myflags.BoolVar(s, tag, false, usagetag)
				myflags.BoolFunc(tag, usagetag, func(s string) error {
					v := new(bool)
					retriever.val = v
					retriever.touched = true
					b, err := strconv.ParseBool(s)
					if err != nil {
						return err
					}
					*v = b
					return nil

				})

			}

		case reflect.String:
			// Change the value of the field to the tag value
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*string)
						if ok {
							fieldValue.SetString(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					v := new(string)
					retriever.touched = true
					retriever.val = v
					*v = s
					return nil
				})

			}

			// fieldValue.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*int64)
						if ok {
							fieldValue.SetInt(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					v := new(int64)
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					retriever.touched = true
					retriever.val = v
					i, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return err
					}
					*v = i
					return nil
				})

			}
		case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:

			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*uint64)
						if ok {
							fieldValue.SetUint(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					v := new(uint64)
					retriever.val = v
					retriever.touched = true
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					i, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return err
					}
					*v = i
					return nil
				})
			}
		default:
			return nil, fmt.Errorf("(flag) %s underlying type unsupported (setFlagVal)", fieldValue.Type().String())
		}
		return retriever, nil
	}

	setflagValPtr := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string, usagetag string, existing *flagSetRetriever) (retriever *flagSetRetriever, err error) {
		k := fieldValue.Elem().Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.Bool:
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					s, ok := r.val.(*bool)
					if ok {
						fieldValue.Elem().SetBool(*s)
					}
				}
				// else {
				// 	return fmt.Errorf("flag %s underlying interface{} type coercsion failed", tag)
				// }
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				//				myflags.BoolVar(s, tag, false, usagetag)
				myflags.BoolFunc(tag, usagetag, func(s string) error {
					v := new(bool)
					retriever.val = v
					retriever.touched = true
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					b, err := strconv.ParseBool(s)
					if err != nil {
						return err
					}
					*v = b
					return nil

				})

			}
		case reflect.String:
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*string)
						if ok {
							fieldValue.Elem().SetString(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					v := new(string)
					retriever.val = v
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					retriever.touched = true
					*v = s
					return nil
				})

			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*int64)
						if ok {
							fieldValue.Elem().SetInt(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					v := new(int64)
					retriever.val = v
					retriever.touched = true
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					i, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return err
					}
					*v = i
					return nil
				})
			}
		case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			// Change the value of the field to the tag value
			retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
				if r.touched {
					if r.val != nil {
						s, ok := r.val.(*uint64)
						if ok {
							fieldValue.Elem().SetUint(*s)
						} else {
							return fmt.Errorf("flag %s underlying interface{} type coercion failed", tag)
						}
					}
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)

				myflags.Func(tag, usagetag, func(s string) error {
					v := new(uint64)
					retriever.val = v
					retriever.touched = true
					ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
					i, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return err
					}
					*v = i
					return nil
				})
			}
		default:
			return nil, fmt.Errorf("map (env) val for %s underlying type unsupported (setFlagValPtr)", fieldValue.Type().String())
		}
		return retriever, nil
	}

	var findFlags func(parentpath string, somestruct interface{}) (err error)

	findFlags = func(parentpath string, somestruct interface{}) (err error) {
		// Get the value of the input. This will be a reflect.Value
		valuePtr := reflect.ValueOf(somestruct)
		if valuePtr.Kind() != reflect.Ptr {
			return fmt.Errorf("not a pointer to a struct")
		}
		inputValue := valuePtr.Elem()
		// Get the type of the input. This will be a reflect.Type
		inputType := inputValue.Type()

		// verify that input is a struct
		if inputType.Kind() != reflect.Struct {
			return fmt.Errorf("not a struct")
		}

		// Iterate over all fields of the input
		for i := 0; i < inputType.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			field := inputType.Field(i)

			// Get the field tag value
			tag := field.Tag.Get(FLAGFIELD)
			//			defaultval := field.Tag.Get("default")
			conftags := field.Tag.Get(CONFFIELD)
			usagetag := field.Tag.Get(FLAGFIELDUSAGE)
			confops := processConfTagOptsValues(conftags)
			// check if confops has a 'skip' key
			if skipField(confops) {
				continue
			}
			debugf("flag: Field Name: %s, flag val: %s\n", field.Name, tag)
			// if len(defaultval) > 0 {
			// Get the field value
			fieldValue := inputValue.FieldByName(field.Name)
			// Only do substitution if the field value can be changed
			if !field.IsExported() {
				debugf("flag: Field %s is not exported\n", field.Name)
				continue
			}
			if field.Type.Kind() == reflect.Ptr {
				// recurse
				fieldValue := inputValue.FieldByName(field.Name)
				// used to create a temp Value of any kind, based on the underlying
				// type of the Ptr - we use this to inspect the underlying type
				// if the type is a struct, we always create it
				// if the type is a string or number we create it only if we have a default
				// value
				t := fieldValue.Type()
				if fieldValue.IsNil() {
					debugf("flag: Field %s is nil\n", field.Name)
					if skipIfNil(confops) {
						continue
					}
					switch t.Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Bool:
						debugf("flag: Ptr: Underlying fundamental type: %s\n", t.Elem().Kind().String())
						// Does an env refenced exist?
						// if _, ok := m[tag]; ok {
						fieldValue.Set(reflect.New(t.Elem()))
						// }
					default:
						debugf("flag: Got a NON-fundamental type: %s %s which is a %s\n", t.Kind().String(), t.Elem().String(), t.Elem().Kind().String())
						if fieldValue.CanSet() {
							if t.Elem().Kind() == reflect.Struct {
								fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
							} else {
								if len(tag) > 0 {
									err = fmt.Errorf("flag tag for %s underlying type unsupported", field.Name)
								}
								return
							}
						} else {
							debugf("flag: Field %s cannot be set (private ?)\n", field.Name)
						}
					}
				}

				if !fieldValue.IsNil() {
					// is this a Ptr to a struct?
					if t.Elem().Kind() == reflect.Struct {
						err := findFlags(addParentPath(parentpath, field.Name), fieldValue.Elem().Addr().Interface())
						if err != nil {
							return err
						}
					} else {
						// nope then its just a fundamental type
						if len(tag) > 0 {
							existing, ok := ret.needflags[tag] // check if we already have a retriever for this flag
							if ok {
								_, err = setflagValPtr(parentpath, field.Name, fieldValue, tag, usagetag, existing)
								if err != nil {
									return
								}
							} else {
								var retriever *flagSetRetriever
								retriever, err = setflagValPtr(parentpath, field.Name, fieldValue, tag, usagetag, nil)
								if err != nil {
									return
								}
								ret.needflags[tag] = retriever
							}
						}

					}
				}
			} else if field.Type.Kind() == reflect.Struct {
				// recurse
				fieldValue := inputValue.FieldByName(field.Name)
				// is this a Ptr to a struct?
				err := findFlags(addParentPath(parentpath, field.Name), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			} else if fieldValue.CanSet() {
				if len(tag) > 0 {
					existing, ok := ret.needflags[tag] // check if we already have a retriever for this flag
					if ok {
						_, err = setFlagVal(parentpath, field.Name, fieldValue, tag, usagetag, existing)
						if err != nil {
							return
						}
					} else {
						var retriever *flagSetRetriever
						retriever, err = setFlagVal(parentpath, field.Name, fieldValue, tag, usagetag, nil)
						if err != nil {
							return
						}
						ret.needflags[tag] = retriever
					}
				}
			} else {
				if len(tag) > 0 {
					return fmt.Errorf("env for %s cannot be set", field.Name)
				}
			}
			// }

		}
		return nil
	}

	err = findFlags("", somestruct)

	// if err != nil {
	// 	err = innerSubst("", somestruct)
	// }

	return ret, err

}

// RunFlagFlags is a function that takes a pointer to a struct

func FinalizeFlags(tags *ProcessedFlagTags) (err error) {
	if tags == nil {
		err = fmt.Errorf("tags is nil")
		return
	}
	//	flag.Parse()
	for k, v := range tags.needflags {
		// retreve sets the value of fieldValue associate with the field after a Parse
		err = v.retrieve(k)
		if err != nil {
			return fmt.Errorf("error retrieving flag value %s: %v", k, err)
		}
	}
	return nil
}

// This is a convenience function that run the FlagFieldSubstitution and then calls flag.Parse()
// and then calls FinalizeFlags
func ProcessFlags(somestruct interface{}, opts *FlagFieldSubstOpts) (err error) {
	var processed *ProcessedFlagTags

	_, ok := preprocessedStructFlags[somestruct]
	if !ok {
		processed, err = ProcessFlagTags(somestruct, opts)
		if err != nil {
			return
		}
		preprocessedStructFlags[somestruct] = processed
	}
	if opts == nil || opts.UseFlags == nil {
		if !flag.Parsed() {
			flag.Parse()
		}
	} else {
		if !opts.UseFlags.Parsed() {
			argz := opts.Args
			if argz == nil || len(argz) < 1 {
				argz = os.Args[1:]
			}
			err = opts.UseFlags.Parse(argz)
			if err != nil {
				log.Printf("Error parsing flags: %v\n", err)
				return
			}
		}
	}
	// err = FinalizeFlags(processed)
	// if err != nil {
	// 	return
	// }
	for _, v := range preprocessedStructFlags {
		err = FinalizeFlags(v)
		if err != nil {
			return
		}
	}
	return
}

// This is a convenience function that run the FlagFieldSubstitution and then calls flag.Parse()
// and then calls FinalizeFlags - but with a flag.FlagSet passed in
func ProcessFlagsWithFlagSet(somestruct interface{}, set *flag.FlagSet, argz []string) (err error) {
	//	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	var processed *ProcessedFlagTags
	_, ok := preprocessedStructFlags[somestruct]
	if !ok {
		processed, err = ProcessFlagTags(somestruct, &FlagFieldSubstOpts{
			UseFlags: set,
		})
		if err != nil {
			return
		}
		preprocessedStructFlags[somestruct] = processed
	}
	// processed, err = ProcessFlagTags(somestruct, &FlagFieldSubstOpts{
	// 	UseFlags: set,
	// })
	// if err != nil {
	// 	return
	// }
	err = set.Parse(argz)
	if err != nil {
		log.Printf("Error parsing flags: %v\n", err)
		return
	}
	// err = FinalizeFlags(processed)
	// if err != nil {
	// 	return
	// }
	for _, v := range preprocessedStructFlags {
		err = FinalizeFlags(v)
		if err != nil {
			return
		}
	}
	return
}

var preprocessedStructFlags map[interface{}]*ProcessedFlagTags

func init() {
	preprocessedStructFlags = make(map[interface{}]*ProcessedFlagTags)
}

// Use this to add in all flags from a given struct *before* flags.Parse() is called anywhere.
// flags.Parse could be handled by the caller later, OR it may be handled by the conftagz.Process()
// function.
func PreProcessFlagsWithFlagSet(somestruct interface{}, set *flag.FlagSet) (err error) {
	//	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	var processed *ProcessedFlagTags
	processed, err = ProcessFlagTags(somestruct, &FlagFieldSubstOpts{
		UseFlags: set,
	})
	if err != nil {
		return
	}
	preprocessedStructFlags[somestruct] = processed
	return
}

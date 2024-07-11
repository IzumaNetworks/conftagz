package conftagz

import (
	//	"flag"
	"fmt"
	"reflect"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

const COBRAFIELD = "cflag"
const COBRAFIELDUSAGE = "usage"
const COBRACMDFIELD = "cobra"

// this is used to set the flag before calling flag.Parse to get all the cmd line options to parse
type cobraFlagRetrieverFunc func(flagname string, self *cobraFlagSetRetriever) error

func (r *cobraFlagSetRetriever) retrieve(flagname string) (err error) {
	for _, retriever := range r.retrievers {
		err = retriever(flagname, r)
		if err != nil {
			return
		}
	}
	return
}

type cobraFlagSetRetriever struct {
	fieldName  string
	fieldValue reflect.Value
	retrievers []cobraFlagRetrieverFunc
	//val        interface{}

	touched bool
	// unfortunately spf13/pflag does not implement the flag.Func() functions since it's like almost
	// never updated, so we resort to just using its function which return pointers to vars if the flag is seen
	varstr  string
	varbool bool
	varint  int64
	varuint uint64
}

type ProcessedCobraTags struct {
	needflags map[string]*cobraFlagSetRetriever
	// true if we have ran ProcessAllFlagTags
	flagsProcessed bool
	fieldsTouched  []string
}

func (p *ProcessedCobraTags) GetFlagsFound() (ret []string) {
	ret = make([]string, 0)
	for k := range p.needflags {
		ret = append(ret, k)
	}
	return
}

func (p *ProcessedCobraTags) GetFieldsTouched() (ret []string) {
	return p.fieldsTouched
}

type CobraFieldSubstOpts struct {
	// throws an error if the environment variable is not found
	//	ThrowErrorIfEnvMissing bool
	//	UseFlags *flag.FlagSet
	Args []string
	Tags *ProcessedCobraTags
}

var cobraCommands map[string]*cobra.Command

func RegisterCobraCmd(name string, cmd *cobra.Command) {
	if cobraCommands == nil {
		cobraCommands = make(map[string]*cobra.Command)
	}
	cobraCommands[name] = cmd
}

func ProcessCobraTags(somestruct interface{}, opts *CobraFieldSubstOpts) (ret *ProcessedCobraTags, err error) {
	ret = &ProcessedCobraTags{}
	if opts == nil {
		opts = &CobraFieldSubstOpts{}

	}
	if opts.Tags == nil {
		opts.Tags = ret
	}
	if opts.Tags.flagsProcessed {
		return opts.Tags, nil
	}
	ret.flagsProcessed = true
	// this will record all the flags we need to set
	ret.needflags = make(map[string]*cobraFlagSetRetriever)

	if opts == nil {
		opts = &CobraFieldSubstOpts{}
	}
	// myflags := opts.UseFlags

	// if myflags == nil {
	// 	myflags = flag.CommandLine
	// }

	setFlagVal := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string, stag string, usagetag string, existing *cobraFlagSetRetriever, myflags *flag.FlagSet) (retriever *cobraFlagSetRetriever, err error) {
		k := fieldValue.Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.Bool:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varbool {
					fieldValue.SetBool(r.varbool)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				retriever.varbool = false
				if len(stag) > 0 {
					myflags.BoolVarP(&retriever.varbool, tag, stag, false, usagetag)
				} else {
					myflags.BoolVar(&retriever.varbool, tag, false, usagetag)
				}

			}

		case reflect.String:
			// Change the value of the field to the tag value
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if len(r.varstr) > 0 {
					fieldValue.SetString(r.varstr)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.StringVarP(&retriever.varstr, tag, stag, "", usagetag)
				} else {
					myflags.StringVar(&retriever.varstr, tag, "", usagetag)
				}
			}

			// fieldValue.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varint != 0 {
					fieldValue.SetInt(r.varint)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.Int64VarP(&retriever.varint, tag, stag, 0, usagetag)
				} else {
					myflags.Int64Var(&retriever.varint, tag, 0, usagetag)
				}
			}
		case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varint != 0 {
					fieldValue.SetUint(r.varuint)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.Uint64VarP(&retriever.varuint, tag, stag, 0, usagetag)
				} else {
					myflags.Uint64Var(&retriever.varuint, tag, 0, usagetag)
				}
			}
		default:
			return nil, fmt.Errorf("(flag) %s underlying type unsupported (setFlagVal)", fieldValue.Type().String())
		}
		return retriever, nil
	}

	setflagValPtr := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string, stag string, usagetag string, existing *cobraFlagSetRetriever, myflags *flag.FlagSet) (retriever *cobraFlagSetRetriever, err error) {
		k := fieldValue.Elem().Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.Bool:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varbool {
					fieldValue.Elem().SetBool(r.varbool)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				retriever.varbool = false
				if len(stag) > 0 {
					myflags.BoolVarP(&retriever.varbool, tag, stag, false, usagetag)
				} else {
					myflags.BoolVar(&retriever.varbool, tag, false, usagetag)
				}
			}
			// retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
			// 	if r.touched {
			// 		s, ok := r.val.(*bool)
			// 		if ok {
			// 			fieldValue.Elem().SetBool(*s)
			// 		}
			// 	}
			// 	// else {
			// 	// 	return fmt.Errorf("flag %s underlying interface{} type coercsion failed", tag)
			// 	// }
			// 	return nil
			// }
			// if existing != nil {
			// 	existing.retrievers = append(existing.retrievers, retrieverfunc)
			// } else {
			// 	retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
			// 	retriever.retrievers = append(retriever.retrievers, retrieverfunc)

			// 	//				myflags.BoolVar(s, tag, false, usagetag)
			// 	myflags.BoolFunc(tag, usagetag, func(s string) error {
			// 		v := new(bool)
			// 		retriever.val = v
			// 		retriever.touched = true
			// 		ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
			// 		b, err := strconv.ParseBool(s)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		*v = b
			// 		return nil

			// 	})

			// }
		case reflect.String:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if len(r.varstr) > 0 {
					fieldValue.Elem().SetString(r.varstr)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.StringVarP(&retriever.varstr, tag, stag, "", usagetag)
				} else {
					myflags.StringVar(&retriever.varstr, tag, "", usagetag)
				}
			}

			// retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
			// 	if r.touched {
			// 		if r.val != nil {
			// 			s, ok := r.val.(*string)
			// 			if ok {
			// 				fieldValue.Elem().SetString(*s)
			// 			} else {
			// 				return fmt.Errorf("flag %s underlying interface{} type coercsion failed", tag)
			// 			}
			// 		}
			// 	}
			// 	return nil
			// }
			// if existing != nil {
			// 	existing.retrievers = append(existing.retrievers, retrieverfunc)
			// } else {
			// 	retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
			// 	retriever.retrievers = append(retriever.retrievers, retrieverfunc)

			// 	myflags.Func(tag, usagetag, func(s string) error {
			// 		v := new(string)
			// 		retriever.val = v
			// 		ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
			// 		retriever.touched = true
			// 		*v = s
			// 		return nil
			// 	})

			// }
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varint != 0 {
					fieldValue.Elem().SetInt(r.varint)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.Int64VarP(&retriever.varint, tag, stag, 0, usagetag)
				} else {
					myflags.Int64Var(&retriever.varint, tag, 0, usagetag)
				}
			}
		case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			retrieverfunc := func(flagname string, r *cobraFlagSetRetriever) (err error) {
				if r.varint != 0 {
					fieldValue.Elem().SetUint(r.varuint)
				}
				return nil
			}
			if existing != nil {
				existing.retrievers = append(existing.retrievers, retrieverfunc)
			} else {
				retriever = &cobraFlagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
				retriever.retrievers = append(retriever.retrievers, retrieverfunc)
				if len(stag) > 0 {
					myflags.Uint64VarP(&retriever.varuint, tag, stag, 0, usagetag)
				} else {
					myflags.Uint64Var(&retriever.varuint, tag, 0, usagetag)
				}
			}
		// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 	retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
		// 		if r.touched {
		// 			if r.val != nil {
		// 				s, ok := r.val.(*int64)
		// 				if ok {
		// 					fieldValue.Elem().SetInt(*s)
		// 				} else {
		// 					return fmt.Errorf("flag %s underlying interface{} type coercsion failed", tag)
		// 				}
		// 			}
		// 		}
		// 		return nil
		// 	}
		// 	if existing != nil {
		// 		existing.retrievers = append(existing.retrievers, retrieverfunc)
		// 	} else {
		// 		retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
		// 		retriever.retrievers = append(retriever.retrievers, retrieverfunc)

		// 		myflags.Func(tag, usagetag, func(s string) error {
		// 			v := new(int64)
		// 			retriever.val = v
		// 			retriever.touched = true
		// 			ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
		// 			i, err := strconv.ParseInt(s, 10, 64)
		// 			if err != nil {
		// 				return err
		// 			}
		// 			*v = i
		// 			return nil
		// 		})
		// 	}
		// case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		// 	// Change the value of the field to the tag value
		// 	retrieverfunc := func(flagname string, r *flagSetRetriever) (err error) {
		// 		if r.touched {
		// 			if r.val != nil {
		// 				s, ok := r.val.(*uint64)
		// 				if ok {
		// 					fieldValue.Elem().SetUint(*s)
		// 				} else {
		// 					return fmt.Errorf("flag %s underlying interface{} type coercsion failed", tag)
		// 				}
		// 			}
		// 		}
		// 		return nil
		// 	}
		// 	if existing != nil {
		// 		existing.retrievers = append(existing.retrievers, retrieverfunc)
		// 	} else {
		// 		retriever = &flagSetRetriever{fieldName: fieldName, fieldValue: fieldValue}
		// 		retriever.retrievers = append(retriever.retrievers, retrieverfunc)

		// 		myflags.Func(tag, usagetag, func(s string) error {
		// 			v := new(uint64)
		// 			retriever.val = v
		// 			retriever.touched = true
		// 			ret.fieldsTouched = append(ret.fieldsTouched, addParentPath(parentpath, fieldName))
		// 			i, err := strconv.ParseUint(s, 10, 64)
		// 			if err != nil {
		// 				return err
		// 			}
		// 			*v = i
		// 			return nil
		// 		})
		// 	}
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

			// Get the field 'ctag' value
			// stag - is the shorttag which is optional. Example: ctag:"verbose,v"
			// then stag is "v"
			var stag string
			tag := field.Tag.Get(COBRAFIELD)
			// split on comma
			allctags := strings.Split(tag, ",")
			if len(allctags) > 1 {
				tag = allctags[0]
				stag = allctags[1]
			}
			//			defaultval := field.Tag.Get("default")
			conftags := field.Tag.Get(CONFFIELD)
			usagetag := field.Tag.Get(FLAGFIELDUSAGE)
			cobracmdtag := field.Tag.Get(COBRACMDFIELD)

			// persistent flag? true is persistent, false is local
			var persist bool
			allcmdtags := strings.Split(cobracmdtag, ",")
			if len(allcmdtags) > 1 {
				cobracmdtag = allcmdtags[0]
				if allcmdtags[1] == "persistent" {
					persist = true
				} else {
					err = fmt.Errorf("field %s: invalid modifier on cobra command tag: %s", field.Name, allcmdtags[1])
					return
				}
			}

			confops := processConfTagOptsValues(conftags)
			// check if confops has a 'skip' key
			if skipField(confops) {
				continue
			}
			debugf("cflag: Field Name: %s, cflag val: %s, cobra cmd: %s\n", field.Name, tag, cobracmdtag)
			// if len(defaultval) > 0 {
			// Get the field value
			fieldValue := inputValue.FieldByName(field.Name)
			// if len(cobracmdtag) < 1 && len(tag) < 1 {
			// 	debugf("cflag: Field %s has no cflag or cobra tag\n", field.Name)
			// 	continue
			// }
			// Only do substitution if the field value can be changed
			if !field.IsExported() {
				debugf("cflag: Field %s is not exported\n", field.Name)
				if len(cobracmdtag) > 0 || len(tag) > 0 {
					err = fmt.Errorf("field %s: field is not exported", field.Name)
					return
				} else {
					continue
				}
			}
			if (len(cobracmdtag) > 0 && len(tag) < 1) || (len(cobracmdtag) < 1 && len(tag) > 0) {
				return fmt.Errorf("field %s: 'cflag' and 'cobra' tag must both be present", field.Name)
			}
			var pflags *flag.FlagSet
			// get pflags for the cobra command
			if len(cobracmdtag) > 0 {
				cmd, ok := cobraCommands[cobracmdtag]
				if !ok {
					return fmt.Errorf("field %s: cobra command %s not found", field.Name, cobracmdtag)
				}
				//				pflags = cmd.PersistentFlags()
				if persist {
					debugf("cflag: Field %s is persistent flag\n", field.Name)
					pflags = cmd.PersistentFlags()
				} else {
					pflags = cmd.Flags()
				}
			}
			// is this a Ptr to a struct?
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
					debugf("cflag: Field %s is nil\n", field.Name)
					if skipIfNil(confops) {
						debugf("cflag: Field %s is nil and skipIfNil is set\n", field.Name)
						continue
					}
					switch t.Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Bool:
						debugf("flag: Ptr: Underlying fundamental type: %s\n", t.Elem().Kind().String())
						// Does an env refenced exist?
						// if _, ok := m[tag]; ok {
						if !skipIfNil(confops) {
							fieldValue.Set(reflect.New(t.Elem()))
						} else {
							debugf("cflag: Field %s is nil and skipIfNil is set\n", field.Name)
							continue
						}
						// }
					default:
						debugf("cflag: Got a NON-fundamental type: %s %s which is a %s\n", t.Kind().String(), t.Elem().String(), t.Elem().Kind().String())
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
								_, err = setflagValPtr(parentpath, field.Name, fieldValue, tag, stag, usagetag, existing, pflags)
								if err != nil {
									return
								}
							} else {
								var retriever *cobraFlagSetRetriever
								retriever, err = setflagValPtr(parentpath, field.Name, fieldValue, tag, stag, usagetag, nil, pflags)
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
						setFlagVal(parentpath, field.Name, fieldValue, tag, stag, usagetag, existing, pflags)
						if err != nil {
							return
						}
					} else {
						var retriever *cobraFlagSetRetriever
						retriever, err = setFlagVal(parentpath, field.Name, fieldValue, tag, stag, usagetag, nil, pflags)
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

	valuePtr := reflect.ValueOf(somestruct)
	if valuePtr.IsZero() {
		return nil, fmt.Errorf("struct ptr is nil")
	}

	if valuePtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not a pointer to a struct")
	}
	inputValue := valuePtr.Elem()
	if inputValue.IsZero() {
		return nil, fmt.Errorf("struct ptr is nil")
	}

	err = findFlags("", somestruct)

	// if err != nil {
	// 	err = innerSubst("", somestruct)
	// }

	return ret, err

}

// RunFlagFlags is a function that takes a pointer to a struct

func FinalizeCobraFlags(tags *ProcessedCobraTags) (err error) {
	if tags == nil {
		err = fmt.Errorf("tags is nil")
		return
	}
	//	flag.Parse()
	for k, v := range tags.needflags {
		debugf("cobra: retrieving flag %s\n", k)
		// retreve sets the value of fieldValue associate with the field after a Parse
		err = v.retrieve(k)
		if err != nil {
			return fmt.Errorf("error retrieving flag value %s: %v", k, err)
		}
	}
	return nil
}

// Should be called before running Process. Call on each struct which may have 'cflags' or other cobra related conftagz
func PreProcessCobraFlags(somestruct interface{}, opts *CobraFieldSubstOpts) (err error) {
	var processed *ProcessedCobraTags
	usingCobraFlags = true
	_, ok := preprocessedCobraStructFlags[somestruct]
	if !ok {
		processed, err = ProcessCobraTags(somestruct, opts)
		if err != nil {
			return
		}
		preprocessedCobraStructFlags[somestruct] = processed
	}
	// if opts != nil {
	// 	if opts.Args != nil && len(opts.Args) > 0 {

	// 	}
	// }
	// if opts == nil || opts.UseFlags == nil {
	// 	if !flag.Parsed() {
	// 		flag.Parse()
	// 	}
	// } else {
	// 	if !opts.UseFlags.Parsed() {
	// 		argz := opts.Args
	// 		if argz == nil || len(argz) < 1 {
	// 			argz = os.Args[1:]
	// 		}
	// 		opts.UseFlags.Parse(argz)
	// 	}
	// }
	// err = FinalizeFlags(processed)
	// if err != nil {
	// 	return
	// }

	// for _, v := range preprocessedCobraStructFlags {
	// 	err = FinalizeCobraFlags(v)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	return
}

func PostProcessCobraFlags() (err error) {
	for _, v := range preprocessedCobraStructFlags {
		err = FinalizeCobraFlags(v)
		if err != nil {
			return
		}
	}
	return
}

// This is a convenience function that run the FlagFieldSubstitution and then calls flag.Parse()
// and then calls FinalizeFlags - but with a flag.FlagSet passed in

// func ProcessCobraFlagsWithFlagSet(somestruct interface{}, set *flag.FlagSet, argz []string) (err error) {
// 	//	flagset := flag.NewFlagSet("test", flag.ExitOnError)
// 	var processed *ProcessedCobraTags
// 	_, ok := preprocessedCobraStructFlags[somestruct]
// 	if !ok {
// 		processed, err = ProcessCobraTags(somestruct, &CobraFieldSubstOpts{
// 			UseFlags: set,
// 		})
// 		if err != nil {
// 			return
// 		}
// 		preprocessedCobraStructFlags[somestruct] = processed
// 	}
// 	// processed, err = ProcessFlagTags(somestruct, &FlagFieldSubstOpts{
// 	// 	UseFlags: set,
// 	// })
// 	// if err != nil {
// 	// 	return
// 	// }
// 	set.Parse(argz)
// 	// err = FinalizeFlags(processed)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	for _, v := range preprocessedCobraStructFlags {
// 		err = FinalizeCobraFlags(v)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

var preprocessedCobraStructFlags map[interface{}]*ProcessedCobraTags

func init() {
	preprocessedCobraStructFlags = make(map[interface{}]*ProcessedCobraTags)
}

// Use this to add in all flags from a given struct *before* flags.Parse() is called anywhere.
// flags.Parse could be handled by the caller later, OR it may be handled by the conftagz.Process()
// function.

// func PreProcessCobraFlagsWithFlagSet(somestruct interface{}, set *flag.FlagSet) (err error) {
// 	//	flagset := flag.NewFlagSet("test", flag.ExitOnError)
// 	var processed *ProcessedCobraTags
// 	processed, err = ProcessCobraTags(somestruct, &CobraFieldSubstOpts{
// 		UseFlags: set,
// 	})
// 	if err != nil {
// 		return
// 	}
// 	preprocessedCobraStructFlags[somestruct] = processed
// 	return
// }

package conftagz

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type DefaultFieldSubstOpts struct {
	// throws an error if the environment variable is not found
}

type DefaultFunc func(fieldname string) interface{}

var defaultFuncs map[string]DefaultFunc

func RegisterDefaultFunc(id string, f DefaultFunc) map[string]DefaultFunc {
	if defaultFuncs == nil {
		defaultFuncs = make(map[string]DefaultFunc)
	}
	defaultFuncs[id] = f
	return defaultFuncs
}

var matchDefaultFuncPat = `^\s*\$\(([a-z,A-Z,_]+[a-z,A-Z,0-9,\_]*)\)\s*$`

var matchDefaultFuncRE = regexp.MustCompile(matchDefaultFuncPat)

// EnvFieldSubstitutionFromMap is a function that takes a pointer to a struct
func SubsistuteDefaults(somestruct interface{}, opts *DefaultFieldSubstOpts) (ret []string, err error) {

	var innerSubst func(parentpath string, somestruct interface{}) (err error)

	setDefaultSlice := func(sliceValue reflect.Value, defaultval string) error {
		parsedVals := strings.Split(defaultval, ",")
		k := sliceValue.Type().Elem().Kind()
		switch k {
		case reflect.Ptr:
			// TODO add support for Ptr to Structs
		case reflect.String:
			for _, parsedVal := range parsedVals {
				// Change the value of the field to the tag value
				sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(parsedVal)))
			}
			// TODO add float
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for _, parsedVal := range parsedVals {
				val, err := StringToInt64(parsedVal)
				if err != nil {
					return fmt.Errorf("default value %s not a number", defaultval)
				}
				// Change the value of the field to the tag value
				if reflect.ValueOf(val).CanConvert(sliceValue.Type().Elem()) {
					sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val).Convert(sliceValue.Type().Elem())))
				} else {
					return fmt.Errorf("default value %s has number out of range for %s", defaultval, sliceValue.Type().Elem().String())
				}
			}
		case reflect.Float64, reflect.Float32:
			for _, parsedVal := range parsedVals {
				val, err := StringToFloat64(parsedVal)
				if err != nil {
					return fmt.Errorf("default value %s not a number", defaultval)
				}
				// Change the value of the field to the tag value
				if reflect.ValueOf(val).CanConvert(sliceValue.Type().Elem()) {
					sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val).Convert(sliceValue.Type().Elem())))
				} else {
					return fmt.Errorf("default value %s has number out of range for %s", defaultval, sliceValue.Type().Elem().String())
				}
			}
		default:
			debugf("default: default for %s underlying type unsupported (setDefault)", k.String())
		}
		return nil
	}

	setDefault := func(parentpath string, fieldName string, fieldValue reflect.Value, defaultval string) error {
		var f DefaultFunc
		matches := matchDefaultFuncRE.FindAllStringSubmatch(defaultval, -1)
		if len(matches) > 0 {
			if len(matches[0]) > 1 {
				debugf("default: Found a default func (setDefault): %s\n", matches[0][1])
				f = defaultFuncs[matches[0][1]]
			}
		}

		k := fieldValue.Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.String:
			if fieldValue.IsZero() {
				// Change the value of the field to the tag value
				if f != nil {
					v, ok := f(fieldName).(string)
					if ok {
						fieldValue.SetString(v)
					} else {
						return fmt.Errorf("default func %s did not return a string", matches[0][1])
					}
				} else {
					fieldValue.SetString(defaultval)
				}
				ret = append(ret, addParentPath(parentpath, fieldName))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.IsZero() {
				// Change the value of the field to the tag value
				// first convert string to int
				if f != nil {
					v, ok := f(fieldName).(int64)
					if ok {
						fieldValue.SetInt(v)
					} else {
						return fmt.Errorf("default func %s did not return an int", matches[0][1])
					}
				} else {
					val, err := StringToInt64(defaultval)
					if err != nil {
						return fmt.Errorf("default value %s not a int", defaultval)
					}
					fieldValue.SetInt(val)
				}

				ret = append(ret, addParentPath(parentpath, fieldName))
			}

		case reflect.Float32, reflect.Float64:
			if fieldValue.IsZero() {
				// Change the value of the field to the tag value
				// first convert string to int
				if f != nil {
					v, ok := f(fieldName).(float64)
					if ok {
						fieldValue.SetFloat(v)
					} else {
						return fmt.Errorf("default func %s did not return an float", matches[0][1])
					}
				} else {
					val, err := StringToFloat64(defaultval)
					if err != nil {
						return fmt.Errorf("default value %s not a float", defaultval)
					}
					fieldValue.SetFloat(val)
				}

				ret = append(ret, addParentPath(parentpath, fieldName))
			}

		default:
			return fmt.Errorf("default for %s underlying type unsupported (setDefault)", fieldValue.Type().String())
		}
		return nil
	}

	setDefaultPtr := func(parentpath string, fieldName string, fieldValue reflect.Value, defaultval string) error {
		var f DefaultFunc
		matches := matchDefaultFuncRE.FindAllStringSubmatch(defaultval, -1)
		if len(matches) > 0 {
			if len(matches[0]) > 1 {
				debugf("default: Found a default func (setDefaultPtr): %s\n", matches[0][1])
				f = defaultFuncs[matches[0][1]]
			}
		}

		k := fieldValue.Elem().Kind()
		switch k {
		// TODO - add support for Ptr to String and Ints
		case reflect.String:
			if fieldValue.Elem().IsZero() {
				if f != nil {
					v, ok := f(fieldName).(string)
					if ok {
						fieldValue.Elem().SetString(v)
					} else {
						return fmt.Errorf("default func %s did not return a string", matches[0][1])
					}
				} else {
					// Change the value of the field to the tag value
					fieldValue.Elem().SetString(defaultval)
				}
				ret = append(ret, addParentPath(parentpath, fieldName))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.Elem().IsZero() {
				if f != nil {
					v, ok := f(fieldName).(int64)
					if ok {
						fieldValue.Elem().SetInt(v)
					} else {
						return fmt.Errorf("default func %s did not return an int", matches[0][1])
					}
				} else {
					// Change the value of the field to the tag value
					// first convert string to int
					val, err := StringToInt64(defaultval)
					if err != nil {
						return fmt.Errorf("default value %s not a int", defaultval)
					}
					fieldValue.Elem().SetInt(val)
				}
				ret = append(ret, addParentPath(parentpath, fieldName))
			}
		case reflect.Float32, reflect.Float64:
			if fieldValue.Elem().IsZero() {
				if f != nil {
					v, ok := f(fieldName).(float64)
					if ok {
						fieldValue.Elem().SetFloat(v)
					} else {
						return fmt.Errorf("default func %s did not return an int", matches[0][1])
					}
				} else {
					// Change the value of the field to the tag value
					// first convert string to int
					val, err := StringToFloat64(defaultval)
					if err != nil {
						return fmt.Errorf("default value %s not a float", defaultval)
					}
					fieldValue.Elem().SetFloat(val)
				}
				ret = append(ret, addParentPath(parentpath, fieldName))
			}
		default:
			return fmt.Errorf("default for %s underlying type unsupported (setDefaultPtr)", fieldValue.Type().String())
		}
		return nil
	}

	innerSubst = func(parentpath string, somestruct interface{}) (err error) {
		// Get the value of the input. This will be a reflect.Value
		valuePtr := reflect.ValueOf(somestruct)
		if valuePtr.Kind() != reflect.Ptr {
			return fmt.Errorf("not a pointer to a struct: was: %s", valuePtr.Kind().String())
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
			defaultval := field.Tag.Get("default")
			conftags := field.Tag.Get("conf")
			confops := processConfTagOptsValues(conftags)
			// check if confops has a 'skip' key
			if skipField(confops) {
				continue
			}
			if defaultSkip(confops) {
				continue
			}
			if skipIfZero(confops) && len(defaultval) > 0 {
				err = fmt.Errorf("conf:skipzero not supported with default tag: %s", field.Name)
				return
			}

			debugf("default: Field Name: %s, Default val: %s\n", field.Name, defaultval)
			// if len(defaultval) > 0 {
			// Get the field value
			fieldValue := inputValue.FieldByName(field.Name)
			// Only do substitution if the field value can be changed
			if field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Slice {
				// recurse
				fieldValue := inputValue.FieldByName(field.Name)
				// used to create a temp Value of any kind, based on the underlying
				// type of the Ptr - we use this to inspect the underlying type
				// if the type is a struct, we always create it
				// if the type is a string or number we create it only if we have a default
				// value
				t := fieldValue.Type()
				if fieldValue.IsNil() {
					debugf("Field %s is nil\n", field.Name)
					if skipIfNil(confops) {
						debugf("default: skipping b/c of skipnil tag %s\n", field.Name)
						continue
					}
					// check if the default tag is func

					var f DefaultFunc
					matches := matchDefaultFuncRE.FindAllStringSubmatch(defaultval, -1)
					if len(matches) > 0 {
						if len(matches[0]) > 1 {
							f = defaultFuncs[matches[0][1]]
							debugf("Found a default func (if Ptr nil): %s %p\n", matches[0][1], f)
						}
					}
					// if so, then we let it do the work since this is a pointer
					var fresult interface{}
					var fresultType reflect.Type
					if f != nil {
						fresult = f(field.Name) // returns an interface{}
						fresultType = reflect.TypeOf(fresult)
					}
					// 	// verify that the func returned a value of the correct type
					// 	if reflect.TypeOf(v) == t.Elem() {
					// 	fieldValue.Set(reflect.ValueOf(v))
					// } else {
					switch t.Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if field.Type.Kind() == reflect.Ptr {
							debugf("default: Ptr: Underlying fundamental type: %s\n", t.Elem().Kind().String())
							if f != nil {
								if fresultType.Kind() == reflect.Ptr && fresultType.Elem().Kind() == t.Elem().Kind() {
									fieldValue.Set(reflect.ValueOf(fresult))
									ret = append(ret, addParentPath(parentpath, field.Name))
									continue
								} else {
									return fmt.Errorf("default func %s did not return a Ptr of the correct type: ", matches[0][1])
								}
							} else {
								if len(defaultval) > 0 {
									fieldValue.Set(reflect.New(t.Elem()))
								}
							}
						}
						if field.Type.Kind() == reflect.Slice {
							if len(defaultval) > 0 {
								debugf("default: Slice: Underlying fundamental type: %s\n", t.Elem().Kind().String())
								fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
							}
						}
					case reflect.Struct:
						debugf("Ptr: Underlying struct type: %s\n", t.Elem().Kind().String())
						if f != nil {
							if fresultType.Kind() == reflect.Ptr && fieldValue.Type().Elem() == fresultType.Elem() {
								debugf("default: Ptr: Func: Underlying struct type: %s\n", t.Elem().String())
								fieldValue.Set(reflect.ValueOf(fresult))
								ret = append(ret, addParentPath(parentpath, field.Name))
								continue
							} else {
								return fmt.Errorf("default func %s did not return a ptr to struct of the correct type: ", matches[0][1])
							}
						} else {
							// no function? ok - then if its a Ptr to a struct, we create it
							// otherwise we ignore it
							if field.Type.Kind() == reflect.Ptr {
								fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
							}
						}
					case reflect.Slice:
						debugf("default: Slice: Underlying slice type: %s\n", t.Elem().Kind().String())
						if fresultType.Kind() == reflect.Slice && fieldValue.Type().Elem() == fresultType.Elem() {
							debugf("default: Slice: Func: Underlying struct type: %s\n", t.Elem().String())
							fieldValue.Set(reflect.ValueOf(fresult))
							ret = append(ret, addParentPath(parentpath, field.Name))
							continue
						} else {
							fieldValue.Set(reflect.MakeSlice(fieldValue.Type().Elem(), 0, 0))
						}
					default:
						debugf("default: ignoring: default for %s underlying type unsupported\n", field.Name)
						continue
						// default:
						// 	debugf("Got a NON-fundamental type: %s %s which is a %s\n", t.Kind().String(), t.Elem().String(), t.Elem().Kind().String())
						// 	switch t.Elem().Kind() {
						// 	}
					}
					// }
				}

				if !fieldValue.IsNil() {
					//					debugf("Field %s is NOT nil\n", field.Name)
					// TODO - add support for Slice here
					if field.Type.Kind() == reflect.Slice {
						if fieldValue.Len() < 1 {
							err = setDefaultSlice(fieldValue, defaultval)
							if err != nil {
								return
							}
							ret = append(ret, addParentPath(parentpath, field.Name))
						} else {

							if field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem().Kind() == reflect.Struct {
								for n := 0; n < fieldValue.Len(); n++ {
									name := fmt.Sprintf("%s[%d]", field.Name, n)
									debugf("default: slice of struct ptr %s\n", name)
									err := innerSubst(addParentPath(parentpath, name), fieldValue.Index(n).Elem().Addr().Interface())
									if err != nil {
										return err
									}
								}
							} else if field.Type.Elem().Kind() == reflect.Struct {
								for n := 0; n < fieldValue.Len(); n++ {
									name := fmt.Sprintf("%s[%d]", field.Name, n)
									debugf("default: slice of struct %s\n", name)
									err := innerSubst(addParentPath(parentpath, name), fieldValue.Index(n).Addr().Interface())
									if err != nil {
										return err
									}
								}
							}
						}
						// TODO - add support for Slice here
					} else

					// is this a Ptr to a struct?
					if t.Elem().Kind() == reflect.Struct {
						err := innerSubst(addParentPath(parentpath, field.Name), fieldValue.Elem().Addr().Interface())
						if err != nil {
							return err
						}
					} else {
						// nope then its just a fundamental type
						if len(defaultval) > 0 {
							err = setDefaultPtr(parentpath, field.Name, fieldValue, defaultval)
							if err != nil {
								return
							}
						}
					}
				}
			} else if field.Type.Kind() == reflect.Struct {
				// recurse
				fieldValue := inputValue.FieldByName(field.Name)
				// is this a Ptr to a struct?
				err := innerSubst(addParentPath(parentpath, field.Name), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			} else if field.Type.Kind() == reflect.Slice {
				// recurse
				debugf("default: FIXME 2\n")

			} else if fieldValue.CanSet() {
				if len(defaultval) > 0 {
					err = setDefault(parentpath, field.Name, fieldValue, defaultval)
					if err != nil {
						return
					}
				}
			} else {
				return fmt.Errorf("default for %s cannot be set", field.Name)
			}
			// }

		}
		return nil
	}

	err = innerSubst("", somestruct)
	return ret, err
}

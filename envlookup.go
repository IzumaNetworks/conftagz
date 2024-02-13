package conftagz

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type EnvFieldSubstOpts struct {
	// throws an error if the environment variable is not found
	ThrowErrorIfEnvMissing bool
}

const ENVFIELD = "env"

func EnvToMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		splitEnv := strings.SplitN(env, "=", 2)
		envMap[splitEnv[0]] = splitEnv[1]
	}
	return envMap
}

// EnvFieldSubstitution is a function that takes a pointer to a struct
// and looks at each field. If the field has a ENVFIELD tag ("env" by default)
// then it will look up the value of the field in the environment variables
// and replace the field with the value.
// It returns a list of the names of the fields that were substituted - as
// a list of string
// If there is an error, it returns an error
func EnvFieldSubstitution(somestruct interface{}, opts *EnvFieldSubstOpts) (ret []string, err error) {
	m := EnvToMap()
	return EnvFieldSubstitutionFromMap(somestruct, opts, m)
}
func StringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return i, nil
}
func StringToUint64(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return i, nil
}
func StringToFloat64(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return i, nil
}

func addParentPath(parentpath string, fieldname string) string {
	if len(parentpath) > 0 {
		return parentpath + "." + fieldname
	}
	return fieldname
}

// EnvFieldSubstitutionFromMap is a function that takes a pointer to a struct
func EnvFieldSubstitutionFromMap(somestruct interface{}, opts *EnvFieldSubstOpts, m map[string]string) (ret []string, err error) {
	var throwErrorIfEnvMissing bool
	if opts != nil {
		throwErrorIfEnvMissing = opts.ThrowErrorIfEnvMissing
	}

	setEnvVal := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string) error {
		if val, ok := m[tag]; ok {
			k := fieldValue.Kind()
			switch k {
			// TODO - add support for Ptr to String and Ints
			case reflect.String:
				// Change the value of the field to the tag value
				fieldValue.SetString(val)
				ret = append(ret, addParentPath(parentpath, fieldName))
			case reflect.Bool:
				// if env var is anything other than empty or "0" or "false"
				// then make true
				if len(val) > 0 && val != "0" && val != "false" {
					fieldValue.SetBool(true)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// Change the value of the field to the tag value
				// first convert string to int
				nval, err := StringToInt64(val)
				if err != nil {
					return fmt.Errorf("map (env) %s value %s not a number", tag, val)
				}

				fieldValue.SetInt(nval)
				ret = append(ret, addParentPath(parentpath, fieldName))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// Change the value of the field to the tag value
				// first convert string to int
				nval, err := StringToUint64(val)
				if err != nil {
					return fmt.Errorf("map (env) %s value %s not a number", tag, val)
				}

				fieldValue.SetUint(nval)
				ret = append(ret, addParentPath(parentpath, fieldName))
			default:
				return fmt.Errorf("map (env) for %s underlying type unsupported (setEnvVal)", fieldValue.Type().String())
			}
		} else {
			if throwErrorIfEnvMissing {
				return fmt.Errorf("env %s not found", tag)
			}
		}
		return nil
	}

	setEnvValPtr := func(parentpath string, fieldName string, fieldValue reflect.Value, tag string) error {
		if val, ok := m[tag]; ok {
			k := fieldValue.Elem().Kind()
			switch k {
			// TODO - add support for Ptr to String and Ints
			case reflect.String:
				// Change the value of the field to the tag value
				fieldValue.Elem().SetString(val)
				ret = append(ret, addParentPath(parentpath, fieldName))
			case reflect.Bool:
				// if env var is anything other than empty or "0" or "false"
				// then make true
				if len(val) > 0 && val != "0" && val != "false" {
					fieldValue.Elem().SetBool(true)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// Change the value of the field to the tag value
				// first convert string to int
				nval, err := StringToInt64(val)
				if err != nil {
					return fmt.Errorf("map (env) %s value %s not a number", tag, val)
				}

				fieldValue.Elem().SetInt(nval)
				ret = append(ret, addParentPath(parentpath, fieldName))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// Change the value of the field to the tag value
				// first convert string to int
				nval, err := StringToUint64(val)
				if err != nil {
					return fmt.Errorf("map (env) %s value %s not a number", tag, val)
				}

				fieldValue.Elem().SetUint(nval)
				ret = append(ret, addParentPath(parentpath, fieldName))
			default:
				return fmt.Errorf("map (env) val for %s underlying type unsupported (setEnvValPtr)", fieldValue.Type().String())
			}
			return nil
		} else {
			if throwErrorIfEnvMissing {
				return fmt.Errorf("env %s not found", tag)
			}
		}
		return nil
	}

	var innerSubst func(parentpath string, somestruct interface{}) (err error)

	innerSubst = func(parentpath string, somestruct interface{}) (err error) {
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
			tag := field.Tag.Get("env")
			//			defaultval := field.Tag.Get("default")
			conftags := field.Tag.Get("conf")
			confops := processConfTagOptsValues(conftags)
			// check if confops has a 'skip' key
			if skipField(confops) {
				continue
			}
			if envSkip(confops) {
				continue
			}
			debugf("env: Field Name: %s, Env val: %s\n", field.Name, tag)
			// if len(defaultval) > 0 {
			// Get the field value
			fieldValue := inputValue.FieldByName(field.Name)
			// Only do substitution if the field value can be changed
			if !field.IsExported() {
				debugf("env: Field %s is not exported\n", field.Name)
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
					debugf("env: Field %s is nil\n", field.Name)
					if skipIfNil(confops) {
						continue
					}
					switch t.Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Bool:
						debugf("env: Ptr: Underlying fundamental type: %s\n", t.Elem().Kind().String())
						// Does an env refenced exist?
						if _, ok := m[tag]; ok {
							fieldValue.Set(reflect.New(t.Elem()))
						}
					default:
						debugf("env: Got a NON-fundamental type: %s %s which is a %s\n", t.Kind().String(), t.Elem().String(), t.Elem().Kind().String())
						if fieldValue.CanSet() {
							if t.Elem().Kind() == reflect.Struct {
								fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
							} else {
								err = fmt.Errorf("default for %s underlying type unsupported", field.Name)
								return
							}
						} else {
							debugf("env: Field %s cannot be set (private ?)\n", field.Name)
						}
					}
				}

				if !fieldValue.IsNil() {
					// is this a Ptr to a struct?
					if t.Elem().Kind() == reflect.Struct {
						err := innerSubst(addParentPath(parentpath, field.Name), fieldValue.Elem().Addr().Interface())
						if err != nil {
							return err
						}
					} else {
						// nope then its just a fundamental type
						if len(tag) > 0 {
							err = setEnvValPtr(parentpath, field.Name, fieldValue, tag)
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
			} else if fieldValue.CanSet() {
				if len(tag) > 0 {
					err = setEnvVal(parentpath, field.Name, fieldValue, tag)
					if err != nil {
						return
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

	err = innerSubst("", somestruct)
	return ret, err

}

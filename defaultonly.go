package conftagz

import (
	"fmt"
	"reflect"
	"strings"
)

// GetDefaultOnlyStruct creates a new struct populated with default values from `default` struct tags.
// Unlike SubsistuteDefaults, this function:
// - Returns a NEW struct (does not modify the input)
// - Ignores existing values in the input struct (always applies defaults)
// - Only processes `default` tags (no env, flags, validation)
// - Supports nested structs (both pointers and embedded)
// - Supports custom default functions $(funcName)
//
// The input must be a pointer to a struct. The returned value is a pointer to a new struct
// of the same type with all defaults applied.
//
// Example:
//
//	result, err := conftagz.GetDefaultOnlyStruct(&MyStruct{})
//	myStruct := result.(*MyStruct)
func GetDefaultOnlyStruct(somestruct interface{}) (interface{}, error) {
	valuePtr := reflect.ValueOf(somestruct)
	if valuePtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not a pointer to a struct: was: %s", valuePtr.Kind().String())
	}
	inputValue := valuePtr.Elem()
	inputType := inputValue.Type()

	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: was: %s", inputType.Kind().String())
	}

	// Create a new instance of the same type
	newStructPtr := reflect.New(inputType)

	// Apply defaults to the new struct
	err := applyDefaultsRecursive("", newStructPtr.Interface())
	if err != nil {
		return nil, err
	}

	return newStructPtr.Interface(), nil
}

// applyDefaultsRecursive applies default values to a struct, always setting defaults
// regardless of existing values.
func applyDefaultsRecursive(parentpath string, somestruct interface{}) error {
	valuePtr := reflect.ValueOf(somestruct)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer to a struct: was: %s", valuePtr.Kind().String())
	}
	inputValue := valuePtr.Elem()
	inputType := inputValue.Type()

	if inputType.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: was: %s", inputType.Kind().String())
	}

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)

		if !field.IsExported() {
			continue
		}

		defaultval := field.Tag.Get("default")
		conftags := field.Tag.Get("conf")
		confops := processConfTagOptsValues(conftags)

		if skipField(confops) {
			continue
		}
		if defaultSkip(confops) {
			continue
		}

		fieldValue := inputValue.FieldByName(field.Name)

		switch field.Type.Kind() {
		case reflect.Ptr:
			err := handlePointerField(parentpath, field, fieldValue, defaultval)
			if err != nil {
				return err
			}

		case reflect.Struct:
			err := applyDefaultsRecursive(addParentPath(parentpath, field.Name), fieldValue.Addr().Interface())
			if err != nil {
				return err
			}

		case reflect.Slice:
			if len(defaultval) > 0 {
				err := setDefaultSliceAlways(fieldValue, defaultval)
				if err != nil {
					return err
				}
			}

		default:
			if len(defaultval) > 0 && fieldValue.CanSet() {
				err := setDefaultAlways(field.Name, fieldValue, defaultval)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// handlePointerField handles pointer fields, creating new instances and recursing for struct pointers
func handlePointerField(parentpath string, field reflect.StructField, fieldValue reflect.Value, defaultval string) error {
	elemType := field.Type.Elem()

	// Check for custom default function
	var f DefaultFunc
	matches := matchDefaultFuncRE.FindAllStringSubmatch(defaultval, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		f = defaultFuncs[matches[0][1]]
	}

	switch elemType.Kind() {
	case reflect.Struct:
		if f != nil {
			// Custom function returns the struct
			fresult := f(field.Name)
			fresultType := reflect.TypeOf(fresult)
			if fresultType.Kind() == reflect.Ptr && fieldValue.Type().Elem() == fresultType.Elem() {
				fieldValue.Set(reflect.ValueOf(fresult))
				return nil
			}
			return fmt.Errorf("default func %s did not return a ptr to struct of the correct type", matches[0][1])
		}
		// Create a new struct and recurse
		newStruct := reflect.New(elemType)
		fieldValue.Set(newStruct)
		return applyDefaultsRecursive(addParentPath(parentpath, field.Name), newStruct.Interface())

	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		if f != nil {
			fresult := f(field.Name)
			fresultType := reflect.TypeOf(fresult)
			if fresultType.Kind() == reflect.Ptr && fresultType.Elem().Kind() == elemType.Kind() {
				fieldValue.Set(reflect.ValueOf(fresult))
				return nil
			}
			return fmt.Errorf("default func %s did not return a Ptr of the correct type", matches[0][1])
		}
		if len(defaultval) > 0 {
			newPtr := reflect.New(elemType)
			fieldValue.Set(newPtr)
			return setDefaultAlways(field.Name, newPtr.Elem(), defaultval)
		}

	case reflect.Slice:
		if f != nil {
			fresult := f(field.Name)
			fresultType := reflect.TypeOf(fresult)
			if fresultType.Kind() == reflect.Slice && fieldValue.Type().Elem() == fresultType.Elem() {
				fieldValue.Set(reflect.ValueOf(fresult))
				return nil
			}
		}
		// Create empty slice if no default
		if len(defaultval) > 0 {
			newSlice := reflect.MakeSlice(field.Type.Elem(), 0, 0)
			newSlicePtr := reflect.New(field.Type.Elem())
			newSlicePtr.Elem().Set(newSlice)
			fieldValue.Set(newSlicePtr)
			return setDefaultSliceAlways(newSlicePtr.Elem(), defaultval)
		}
	}

	return nil
}

// setDefaultAlways sets a default value on a field, always overwriting (no IsZero check)
func setDefaultAlways(fieldName string, fieldValue reflect.Value, defaultval string) error {
	var f DefaultFunc
	matches := matchDefaultFuncRE.FindAllStringSubmatch(defaultval, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		f = defaultFuncs[matches[0][1]]
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		if f != nil {
			v, ok := f(fieldName).(string)
			if !ok {
				return fmt.Errorf("default func %s did not return a string", matches[0][1])
			}
			fieldValue.SetString(v)
		} else {
			fieldValue.SetString(defaultval)
		}

	case reflect.Bool:
		if f != nil {
			v, ok := f(fieldName).(bool)
			if !ok {
				return fmt.Errorf("default func %s did not return a bool", matches[0][1])
			}
			fieldValue.SetBool(v)
		} else {
			val := strings.ToLower(defaultval)
			fieldValue.SetBool(val == "true" || val == "1" || val == "yes")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f != nil {
			v, ok := f(fieldName).(int64)
			if !ok {
				return fmt.Errorf("default func %s did not return an int", matches[0][1])
			}
			fieldValue.SetInt(v)
		} else {
			val, err := StringToInt64(defaultval)
			if err != nil {
				return fmt.Errorf("default value %s not an int", defaultval)
			}
			fieldValue.SetInt(val)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if f != nil {
			v, ok := f(fieldName).(uint64)
			if !ok {
				return fmt.Errorf("default func %s did not return a uint", matches[0][1])
			}
			fieldValue.SetUint(v)
		} else {
			val, err := StringToUint64(defaultval)
			if err != nil {
				return fmt.Errorf("default value %s not a uint", defaultval)
			}
			fieldValue.SetUint(val)
		}

	case reflect.Float32, reflect.Float64:
		if f != nil {
			v, ok := f(fieldName).(float64)
			if !ok {
				return fmt.Errorf("default func %s did not return a float", matches[0][1])
			}
			fieldValue.SetFloat(v)
		} else {
			val, err := StringToFloat64(defaultval)
			if err != nil {
				return fmt.Errorf("default value %s not a float", defaultval)
			}
			fieldValue.SetFloat(val)
		}

	default:
		return fmt.Errorf("default for %s underlying type %s unsupported", fieldName, k.String())
	}

	return nil
}

// setDefaultSliceAlways sets default values on a slice field, always creating a new slice
func setDefaultSliceAlways(sliceValue reflect.Value, defaultval string) error {
	parsedVals := strings.Split(defaultval, ",")
	k := sliceValue.Type().Elem().Kind()

	switch k {
	case reflect.String:
		for _, parsedVal := range parsedVals {
			sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(parsedVal)))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, parsedVal := range parsedVals {
			val, err := StringToInt64(parsedVal)
			if err != nil {
				return fmt.Errorf("default value %s not a number", defaultval)
			}
			if reflect.ValueOf(val).CanConvert(sliceValue.Type().Elem()) {
				sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val).Convert(sliceValue.Type().Elem())))
			} else {
				return fmt.Errorf("default value %s has number out of range for %s", defaultval, sliceValue.Type().Elem().String())
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for _, parsedVal := range parsedVals {
			val, err := StringToUint64(parsedVal)
			if err != nil {
				return fmt.Errorf("default value %s not a number", defaultval)
			}
			if reflect.ValueOf(val).CanConvert(sliceValue.Type().Elem()) {
				sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val).Convert(sliceValue.Type().Elem())))
			} else {
				return fmt.Errorf("default value %s has number out of range for %s", defaultval, sliceValue.Type().Elem().String())
			}
		}

	case reflect.Float32, reflect.Float64:
		for _, parsedVal := range parsedVals {
			val, err := StringToFloat64(parsedVal)
			if err != nil {
				return fmt.Errorf("default value %s not a number", defaultval)
			}
			if reflect.ValueOf(val).CanConvert(sliceValue.Type().Elem()) {
				sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val).Convert(sliceValue.Type().Elem())))
			} else {
				return fmt.Errorf("default value %s has number out of range for %s", defaultval, sliceValue.Type().Elem().String())
			}
		}

	case reflect.Bool:
		for _, parsedVal := range parsedVals {
			val := strings.ToLower(strings.TrimSpace(parsedVal))
			sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(val == "true" || val == "1" || val == "yes")))
		}
	}

	return nil
}

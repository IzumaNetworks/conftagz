package conftagz

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	NONE     int = 0
	EQ       int = iota // =
	LT                  // <
	GT                  // >
	GTE                 // >-
	LTE                 // <=
	REGEX               // ~
	TESTFUNC            // $(funcname)
)

type testOp struct {
	Operator     int
	ValInt       int64
	ValString    string
	ValFloat     float64
	testFunc     TestFunc
	testFuncName string
	Regexp       *regexp.Regexp
}

type testConfOp struct {
	ops []*testOp
}

type TestFunc func(val interface{}, fieldname string) bool

var testFuncs map[string]TestFunc

func RegisterTestFunc(id string, f TestFunc) map[string]TestFunc {
	if testFuncs == nil {
		testFuncs = make(map[string]TestFunc)
	}
	testFuncs[id] = f
	return testFuncs
}

type TestWarnPrintf func(format string, args ...interface{})

type TestFieldSubstOpts struct {
	// throws an error if the environment variable is not found
	OnlyWarn bool
	WarnFunc TestWarnPrintf
}

var matchTestFuncPat = `^\s*\$\(([a-z,A-Z,_]+[a-z,A-Z,0-9,\_]*)\)\s*$`

var matchTestFuncRE = regexp.MustCompile(matchTestFuncPat)

func runTestFunc(op *testOp, val reflect.Value, fieldName string) (err error) {
	k := val.Kind()
	debugf("test TESTFUNC %s\n", op.testFuncName)
	switch k {
	case reflect.String:
		if !op.testFunc(val.String(), fieldName) {
			err = fmt.Errorf("value %s !$(%s)", val.String(), op.testFuncName)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !op.testFunc(val.Int(), fieldName) {
			err = fmt.Errorf("value %d !$(%s)", val.Int(), op.testFuncName)
		}
	case reflect.Float32, reflect.Float64:
		if !op.testFunc(val.Float(), fieldName) {
			err = fmt.Errorf("value %f !$(%s)", val.Float(), op.testFuncName)
		}
	case reflect.Bool:
		if !op.testFunc(val.Bool(), fieldName) {
			err = fmt.Errorf("value %t !$(%s)", val.Bool(), op.testFuncName)
		}
	case reflect.Ptr:
		if !op.testFunc(val.Interface(), fieldName) {
			err = fmt.Errorf("value for field %s !$(%s)", fieldName, op.testFuncName)
		}
	case reflect.Struct:
		if !op.testFunc(val.Interface(), fieldName) {
			err = fmt.Errorf("value for field %s !$(%s)", fieldName, op.testFuncName)
		}
	case reflect.Slice:
		if !op.testFunc(val.Interface(), fieldName) {
			err = fmt.Errorf("value for field %s !$(%s)", fieldName, op.testFuncName)
		}
	default:
		debugf("test: unsupported type for TESTFUNC\n")
		err = fmt.Errorf("value unsupported type for TESTFUNC")
	}
	return
}

func runTest(op *testConfOp, val reflect.Value, fieldName string) (err error) {
	for _, op := range op.ops {
		switch op.Operator {
		case LTE:
			k := val.Kind()
			switch k {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				debugf("test LTE %d\n", op.ValInt)
				if val.CanInt() {
					if val.Int() <= op.ValInt {
					} else {
						err = fmt.Errorf("value %d ! <= %d", val.Int(), op.ValInt)
					}
				} else {
					err = fmt.Errorf("value for field - can't get int64")
				}
			case reflect.Float32, reflect.Float64:
				debugf("test LTE %f\n", op.ValFloat)
				if val.CanInt() {
					if val.Float() <= op.ValFloat {
					} else {
						err = fmt.Errorf("value %f ! <= %f", val.Float(), op.ValFloat)
					}
				} else {
					err = fmt.Errorf("value for field - can't get float64")
				}
			default:
				err = fmt.Errorf("test operator require numeric type")
			}
		case GTE:
			k := val.Kind()
			switch k {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				debugf("test GTE %d\n", op.ValInt)
				if val.CanInt() {
					if val.Int() >= op.ValInt {
					} else {
						err = fmt.Errorf("value %d ! >= %d", val.Int(), op.ValInt)
					}
				} else {
					err = fmt.Errorf("value for field - can't get int64")
				}
			case reflect.Float32, reflect.Float64:
				debugf("test GTE %f\n", op.ValFloat)
				if val.CanInt() {
					if val.Float() >= op.ValFloat {
					} else {
						err = fmt.Errorf("value %f ! >= %f", val.Float(), op.ValFloat)
					}
				} else {
					err = fmt.Errorf("value for field - can't get float64")
				}
			default:
				err = fmt.Errorf("test operator require numeric type")
			}

		case LT:
			k := val.Kind()
			switch k {
			// TODO - add support for Ptr to String and Ints
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				debugf("test LT %d\n", op.ValInt)
				if val.CanInt() {
					if val.Int() < op.ValInt {
					} else {
						err = fmt.Errorf("value %d ! < %d", val.Int(), op.ValInt)
					}
				} else {
					err = fmt.Errorf("value for field - can't get int64")
				}
			case reflect.Float32, reflect.Float64:
				debugf("test LT %f\n", op.ValFloat)
				if val.CanInt() {
					if val.Float() < op.ValFloat {
					} else {
						err = fmt.Errorf("value %f ! < %f", val.Float(), op.ValFloat)
					}
				} else {
					err = fmt.Errorf("value for field - can't get float64")
				}
			default:
				err = fmt.Errorf("test operator require numeric type")
			}
		case GT:
			k := val.Kind()
			switch k {
			// TODO - add support for Ptr to String and Ints
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				debugf("test GT %d\n", op.ValInt)
				if val.CanInt() {
					if val.Int() > op.ValInt {
					} else {
						err = fmt.Errorf("value %d ! > %d", val.Int(), op.ValInt)
					}
				} else {
					err = fmt.Errorf("value for field - can't get int64")
				}
			case reflect.Float32, reflect.Float64:
				debugf("test GT %f\n", op.ValFloat)
				if val.CanFloat() {
					if val.Float() > op.ValFloat {
					} else {
						err = fmt.Errorf("value %f ! > %f", val.Float(), op.ValFloat)
					}
				} else {
					err = fmt.Errorf("value for field - can't get float64")
				}
			default:
				err = fmt.Errorf("test operator require numeric type")
			}
		case EQ:
			k := val.Kind()
			switch k {
			// TODO - add support for Ptr to String and Ints
			case reflect.String:
				debugf("test EQ %s\n", op.ValString)
				if val.String() == op.ValString {
					return
				} else {
					err = fmt.Errorf("value \"%s\" != \"%s\"", val.String(), op.ValString)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				debugf("test EQ %d\n", op.ValInt)
				if val.CanInt() {
					if val.Int() == op.ValInt {
					} else {
						err = fmt.Errorf("value %d != %d", val.Int(), op.ValInt)
					}
				} else {
					err = fmt.Errorf("value for field - can't get int64")
				}
			case reflect.Float32, reflect.Float64:
				debugf("test EQ %f\n", op.ValFloat)
				if val.CanFloat() {
					if val.Float() == op.ValFloat {
					} else {
						err = fmt.Errorf("value %f != %f", val.Float(), op.ValFloat)
					}
				} else {
					err = fmt.Errorf("value for field - can't get float64")
				}
			default:
				err = fmt.Errorf("test operator = unsupported type")
			}
		case REGEX:
			k := val.Kind()
			switch k {
			case reflect.String:
				debugf("test REGEX %s\n", op.ValString)
				if !op.Regexp.Match([]byte(val.String())) {
					err = fmt.Errorf("value \"%s\" !~ regexp %s", val.String(), op.ValString)
				}
			default:
				err = fmt.Errorf("value for field - REGEX test operator must be on a string or string* field")
			}
		case TESTFUNC:
			err = runTestFunc(op, val, fieldName)
			// k := val.Kind()
			// debugf("test TESTFUNC %s\n", op.testFuncName)
			// switch k {
			// case reflect.String:
			// 	if !op.testFunc(val.String(), fieldName) {
			// 		err = fmt.Errorf("value %s !$(%s)", val.String(), op.testFuncName)
			// 	}
			// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 	if !op.testFunc(val.Int(), fieldName) {
			// 		err = fmt.Errorf("value %d !$(%s)", val.Int(), op.testFuncName)
			// 	}
			// case reflect.Float32, reflect.Float64:
			// 	if !op.testFunc(val.Float(), fieldName) {
			// 		err = fmt.Errorf("value %f !$(%s)", val.Float(), op.testFuncName)
			// 	}
			// case reflect.Bool:
			// 	if !op.testFunc(val.Bool(), fieldName) {
			// 		err = fmt.Errorf("value %t !$(%s)", val.Bool(), op.testFuncName)
			// 	}
			// case reflect.Ptr:
			// 	if !op.testFunc(val.Interface(), fieldName) {
			// 		err = fmt.Errorf("value %s !$(%s)", val.String(), op.testFuncName)
			// 	}

			// }
		}
	}
	return
}

func parseTestVal(tagval string) (ret *testConfOp, err error) {
	// split tagval by ','
	var op *testOp
	tagval = strings.TrimSpace(tagval)

	// regex can contain commas and there is no point in combining it with other tests
	// so single it out
	if strings.HasPrefix(tagval, "~") {
		parts := strings.Split(tagval, "~")
		op := &testOp{Operator: REGEX, ValString: parts[1]}
		op.Regexp, err = regexp.Compile(op.ValString)
		if err != nil {
			err = fmt.Errorf("test: regexp failed to compile: %s", err.Error())
		} else {
			ret = &testConfOp{}
			ret.ops = append(ret.ops, op)
		}
	} else {

		vals := strings.Split(tagval, ",")

		for _, teststr := range vals {
			// check if this is a test function
			var f TestFunc
			matches := matchTestFuncRE.FindAllStringSubmatch(teststr, -1)
			if len(matches) > 0 {
				if len(matches[0]) > 1 {
					debugf("test: Found a default func (if Ptr nil): %s\n", matches[0][1])
					f = testFuncs[matches[0][1]]
				}
			}

			if f != nil {
				op = &testOp{Operator: TESTFUNC, testFunc: f, testFuncName: matches[0][1]}
			} else {
				// not a testfunc, so parse for other tests
				// remove leading and trailing spaces
				teststr = strings.TrimSpace(teststr)
				var opn int
				var n int
				var c rune
			parsetest:
				for n, c = range teststr {
					switch c {
					case ' ':
						continue parsetest
					case '<':
						opn = LT
						continue parsetest
					case '>':
						opn = GT
						continue parsetest
					case '~':
						err = fmt.Errorf("test: regex may not be combined with other tests")
						return
						// opn = REGEX
						// break parsetest
						// regex
					case '=':
						if opn == LT {
							opn = LTE
							break parsetest
						} else if opn == GT {
							opn = GTE
							break parsetest
						} else {
							opn = EQ
							break parsetest
						}
					default:
						switch opn {
						case LT, GT:
							n--
						default:
							err = fmt.Errorf("invalid test operation")
							return
						}
						break parsetest
					}
				}
				if n+1 > len(teststr) {
					err = fmt.Errorf("invalid test op - bad operand")
					return
				}
				op = &testOp{Operator: opn, ValString: teststr[n+1:]}
				debugf("test: ValString: %s\n", op.ValString)
				switch opn {
				case EQ:
					val, err := StringToInt64(op.ValString)
					if err == nil {
						op.ValInt = val
					}
					valf, err := StringToFloat64(op.ValString)
					if err == nil {
						op.ValFloat = valf
					}
				// case REGEX:
				// 	op.Regexp, err = regexp.Compile(op.ValString)
				// 	if err != nil {
				// 		err = fmt.Errorf("test: regexp failed to compile: %s", err.Error())
				// 	}
				case LTE, GTE, LT, GT:
					val, err := StringToInt64(op.ValString)
					if err != nil {
						return nil, fmt.Errorf("test value %s not a number", op.ValString)
					}
					op.ValInt = val
				}
			}
			if ret == nil {
				ret = &testConfOp{}
			}
			if op != nil {
				ret.ops = append(ret.ops, op)
			}
		}
	}
	return
}

// RUns through all test:"" tags to see if the current value passes the test
func RunTestFlags(somestruct interface{}, opts *TestFieldSubstOpts) (ret []string, err error) {

	var innerTest func(parentpath string, somestruct interface{}) (err error)

	innerTest = func(parentpath string, somestruct interface{}) (err error) {
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
			testval := field.Tag.Get("test")
			conftags := field.Tag.Get("conf")
			confops := processConfTagOptsValues(conftags)
			// check if confops has a 'skip' key
			if skipField(confops) {
				continue
			}
			if testSkip(confops) {
				continue
			}

			var op *testConfOp
			// parse testval
			if len(testval) > 0 {
				op, err = parseTestVal(testval)
				if err != nil {
					err = fmt.Errorf("parse error for test tag for field %s: %s (%s)", addParentPath(parentpath, field.Name), err.Error(), testval)
					return
				}
			}

			debugf("test: Field Name: %s, Test op: %s\n", field.Name, testval)
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
					debugf("test: Field %s is nil\n", field.Name)
					if skipIfNil(confops) {
						continue
					}
					switch t.Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if field.Type.Kind() == reflect.Ptr {
							debugf("test: Ptr: Underlying fundamental type: %s\n", t.Elem().Kind().String())
							if op != nil {
								fieldValue.Set(reflect.New(t.Elem()))
							}
						}
						// if field.Type.Kind() == reflect.Slice {
						// 	debugf("Slice: Underlying fundamental type: %s\n", t.Elem().Kind().String())
						// 	fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
						// }
					case reflect.Struct:
						debugf("test: Ptr: Underlying struct type: %s\n", t.Elem().Kind().String())
						if field.Type.Kind() == reflect.Ptr { // i.e. not a slice
							fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
						}
					// case reflect.Slice:
					// 	debugf("Ptr: Underlying slice type: %s\n", t.Elem().Kind().String())
					// 	fieldValue.Set(reflect.MakeSlice(fieldValue.Type().Elem(), 0, 0))
					default:
						debugf("test: test for %s underlying type unsupported\n", field.Name)
						if op != nil {
							err = fmt.Errorf("test for %s underlying type unsupported", addParentPath(parentpath, field.Name))
							return
						}
						// default:
						// 	debugf("Got a NON-fundamental type: %s %s which is a %s\n", t.Kind().String(), t.Elem().String(), t.Elem().Kind().String())
						// 	switch t.Elem().Kind() {
						// 	}
					}
				}

				if !fieldValue.IsNil() {
					if field.Type.Kind() == reflect.Slice {
						if op != nil {
							err = runTest(op, fieldValue, field.Name)
							ret = append(ret, addParentPath(parentpath, field.Name))
							if err != nil {
								err = fmt.Errorf("field %s: %s", addParentPath(parentpath, field.Name), err.Error())
								return
							}
						} else {
							for i := 0; i < fieldValue.Len(); i++ {
								debugf("test: slice: %s[%d]\n", field.Name, i)
								switch field.Type.Elem().Kind() {
								case reflect.Ptr:
									switch field.Type.Elem().Elem().Kind() {
									case reflect.Struct:
										err := innerTest(addParentPath(parentpath, fmt.Sprintf("%s[%d]", field.Name, i)), fieldValue.Index(i).Elem().Addr().Interface())
										if err != nil {
											return err
										}
									default:
										debugf("test: unsupported slice of type %s - ignoring (2)\n", field.Type.Elem().Kind().String())
									}
								case reflect.Struct:
									err := innerTest(addParentPath(parentpath, fmt.Sprintf("%s[%d]", field.Name, i)), fieldValue.Index(i).Addr().Interface())
									if err != nil {
										return err
									}
								default:
									debugf("test: unsupported slice of type %s - ignoring\n", field.Type.Elem().Kind().String())
								}
							}
						}
					} else

					// is this a Ptr to a struct?
					if t.Elem().Kind() == reflect.Struct {
						// see if there is a test tag for this struct?
						if op != nil {
							debugf("test: found test func for this struct ptr!\n")
							err = runTest(op, fieldValue, field.Name)
							ret = append(ret, addParentPath(parentpath, field.Name))
							if err != nil {
								err = fmt.Errorf("field %s: %s", addParentPath(parentpath, field.Name), err.Error())
								return err
							}
						} else {
							err := innerTest(addParentPath(parentpath, field.Name), fieldValue.Elem().Addr().Interface())
							if err != nil {
								return err
							}
						}
					} else {
						// nope then its just a fundamental type
						if op != nil {
							if fieldValue.Elem().IsZero() && skipIfZero(confops) {
								debugf("test: skip zero (test)\n")
								continue
							}
							err = runTest(op, fieldValue.Elem(), field.Name)
							if err != nil {
								err = fmt.Errorf("field %s: %s", addParentPath(parentpath, field.Name), err.Error())
								return
							}
							ret = append(ret, addParentPath(parentpath, field.Name))
						}
					}
				}
			} else if field.Type.Kind() == reflect.Struct {
				// recurse
				fieldValue := inputValue.FieldByName(field.Name)
				// is this a Ptr to a struct?
				err := innerTest(addParentPath(parentpath, field.Name), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			} else if field.Type.Kind() == reflect.Slice {
				// recurse
				debugf("test: skip slice 2\n")

			} else if fieldValue.CanSet() {
				if op != nil {
					if fieldValue.IsZero() && skipIfZero(confops) {
						debugf("test: skip zero (test) 2\n")
						continue
					}
					err = runTest(op, fieldValue, field.Name)
					if err != nil {
						err = fmt.Errorf("field %s: %s", addParentPath(parentpath, field.Name), err.Error())
						return
					}
					ret = append(ret, addParentPath(parentpath, field.Name))
				}
			} else {
				if op != nil {
					return fmt.Errorf("test for %s cannot be run", field.Name)
				}
			}
			// }

		}
		return nil
	}

	err = innerTest("", somestruct)
	return ret, err
}

package conftagz

import "strings"

// func NewConfTagOptsMap (opts []ConfTagOpts) map[string]string {

func processConfTagOptsValues(conftags string) map[string]string {
	confMap := make(map[string]string)

	if len(conftags) > 0 {
		conftagops := strings.Split(conftags, ",")
		for _, conftagop := range conftagops {
			pair := strings.SplitN(strings.TrimSpace(conftagop), "=", 2)
			if len(pair) == 2 {
				confMap[pair[0]] = pair[1]
			} else {
				confMap[pair[0]] = ""
			}
		}
	}

	return confMap
}

func skipField(confops map[string]string) bool {
	if _, ok := confops["skip"]; ok {
		// skip this field
		return true
	}
	return false
}

// if the value is nil, run the default logic anyway
// even if skipnil is set
func nilDefault(confops map[string]string) bool {
	if _, ok := confops["nildefault"]; ok {
		// skip this field
		return true
	}
	return false
}

func skipIfNil(confops map[string]string) bool {
	if _, ok := confops["skipnil"]; ok {
		// skip this field
		return true
	}
	return false
}

func skipIfZero(confops map[string]string) bool {
	if _, ok := confops["skipzero"]; ok {
		return true
	}
	return false
}

// true if the env var must exist, returns an error if it does not
func mustEnv(confops map[string]string) bool { // nolint:golint,unused
	if _, ok := confops["mustenv"]; ok {
		return true
	}
	return false
}

// The env var processing is entirely skipped
// for this field. Useful if you don't want a blank struct created
// while looking for env vars
func envSkip(confops map[string]string) bool {
	if _, ok := confops["envskip"]; ok {
		return true
	}
	return false
}

func defaultSkip(confops map[string]string) bool {
	if _, ok := confops["defaultskip"]; ok {
		return true
	}
	return false
}

func testSkip(confops map[string]string) bool {
	if _, ok := confops["testskip"]; ok {
		return true
	}
	return false
}

// true if the envvar should always replace a value
// if the env var exists
func preferEnv(confops map[string]string) bool { // nolint:golint,unused
	if _, ok := confops["preferenv"]; ok {
		return true
	}
	return false
}

// if true then env var is only used if the field has the zero value
func backupEnv(confops map[string]string) bool { // nolint:golint,unused
	if _, ok := confops["backupenv"]; ok {
		return true
	}
	return false
}

// if true, then this test will only warn and never
// cause an error to return
func testWarn(confops map[string]string) bool { // nolint:golint,unused
	if _, ok := confops["testwarn"]; ok {
		return true
	}
	return false
}

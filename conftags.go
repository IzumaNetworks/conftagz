package conftagz

// Constants to define flag tag types
const (
	_ int = iota
	FLAGTAGS
	COBRATAGS
	ENVTAGS
	DEFAULTTAGS
	TESTTAGS
)

func defaultOrderOfOps() []int {
	return []int{DEFAULTTAGS, ENVTAGS, FLAGTAGS, TESTTAGS}
}

var usingCobraFlags bool

// mostly just used for testing the library. Returns library to the state it should be on
// initial load
func ResetGlobals() {
	usingCobraFlags = false
}

// this just takes the current order of ops,
// and if it sees FLAGTAGS, it replaces it with COBRATAGS
func switchOrderOfOpsToCobra() []int {
	oo := defaultOrderOfOps()
	for i, op := range oo {
		if op == FLAGTAGS {
			oo[i] = COBRATAGS
		}
	}
	return oo
}

// CONFFIELD is the default configuration field name
const CONFFIELD = "conf"

// ConfTagOpts is a struct that holds the options for processing the tags
type ConfTagOpts struct {
	OrderOfOps   []int
	EnvOpts      *EnvFieldSubstOpts
	TestOpts     *TestFieldSubstOpts
	DefaultOpts  *DefaultFieldSubstOpts
	FlagTagOpts  *FlagFieldSubstOpts
	CobraTagOpts *CobraFieldSubstOpts
}

// Process takes a struct and processes the tags in the struct
func Process(opts *ConfTagOpts, somestruct interface{}) (err error) {
	if opts == nil {
		opts = &ConfTagOpts{}
	}
	if opts.OrderOfOps == nil {
		opts.OrderOfOps = defaultOrderOfOps()
		if usingCobraFlags {
			opts.OrderOfOps = switchOrderOfOpsToCobra()
		}
	}

	for _, op := range opts.OrderOfOps {
		switch op {

		case FLAGTAGS:
			debugf("Processing flag: tags\n")
			if opts.FlagTagOpts == nil {
				opts.FlagTagOpts = &FlagFieldSubstOpts{}
			}
			err = ProcessFlags(somestruct, opts.FlagTagOpts)
			if err != nil {
				return
			}
		case COBRATAGS:
			debugf("Processing cobra: tags\n")
			if opts.FlagTagOpts == nil {
				opts.FlagTagOpts = &FlagFieldSubstOpts{}
			}
			//			_, err = ProcessCobraTags(somestruct, opts.CobraTagOpts)
			err = PostProcessCobraFlags()
			if err != nil {
				return
			}
		case ENVTAGS:
			debugf("Processing env: tags\n")
			if opts.EnvOpts == nil {
				opts.EnvOpts = &EnvFieldSubstOpts{}
			}
			_, err = EnvFieldSubstitution(somestruct, opts.EnvOpts)
			if err != nil {
				return
			}
		case DEFAULTTAGS:
			debugf("Processing default: tags\n")
			if opts.DefaultOpts == nil {
				opts.DefaultOpts = &DefaultFieldSubstOpts{}
			}
			_, err = SubsistuteDefaults(somestruct, opts.DefaultOpts)
			if err != nil {
				return
			}
		case TESTTAGS:
			debugf("Processing test: tags\n")
			if opts.TestOpts == nil {
				opts.TestOpts = &TestFieldSubstOpts{}
			}
			_, err = RunTestFlags(somestruct, opts.TestOpts)
			if err != nil {
				return
			}
		}
	}

	return
}

package conftagz

const (
	_ int = iota
	ENVTAGS
	DEFAULTTAGS
	TESTTAGS
)

func defaultOrderOfOps() []int {
	return []int{ENVTAGS, DEFAULTTAGS, TESTTAGS}
}

type ConfTagOpts struct {
	OrderOfOps  []int
	EnvOpts     *EnvFieldSubstOpts
	TestOpts    *TestFieldSubstOpts
	DefaultOpts *DefaultFieldSubstOpts
}

func Process(opts *ConfTagOpts, somestruct interface{}) (err error) {
	if opts == nil {
		opts = &ConfTagOpts{}
	}
	if opts.OrderOfOps == nil {
		opts.OrderOfOps = defaultOrderOfOps()
	}

	for _, op := range opts.OrderOfOps {
		switch op {
		case ENVTAGS:
			if opts.EnvOpts == nil {
				opts.EnvOpts = &EnvFieldSubstOpts{}
			}
			_, err = EnvFieldSubstitution(somestruct, opts.EnvOpts)
			if err != nil {
				return
			}
		case DEFAULTTAGS:
			if opts.DefaultOpts == nil {
				opts.DefaultOpts = &DefaultFieldSubstOpts{}
			}
			_, err = SubsistuteDefaults(somestruct, opts.DefaultOpts)
			if err != nil {
				return
			}
		case TESTTAGS:
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

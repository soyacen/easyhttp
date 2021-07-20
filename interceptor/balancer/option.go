package easyhttpbalancer

type options struct {
	picker Picker
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(o *options)

func defaultOptions() *options {
	return &options{
		picker: NewFirstPicker(),
	}
}

func WithPicker(picker Picker) Option {
	return func(o *options) {
		o.picker = picker
	}
}

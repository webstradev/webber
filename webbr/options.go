package webbr

import "fmt"

type OptFunc (func(opts *Options))

type Options struct {
	DBName    string
	Extension string
	Encoder   Encoder
	Decoder   Decoder
}

func (o Options) GetDBName() string {
	return fmt.Sprintf("%s.%s", o.DBName, o.Extension)
}

func WithDBName(name string) OptFunc {
	return func(opts *Options) {
		opts.DBName = name
	}
}

func WithExtension(ext string) OptFunc {
	return func(opts *Options) {
		opts.Extension = ext
	}
}

func WithEncoder(enc Encoder) OptFunc {
	return func(opts *Options) {
		opts.Encoder = enc
	}
}

func WithDecoder(dec Decoder) OptFunc {
	return func(opts *Options) {
		opts.Decoder = dec
	}
}

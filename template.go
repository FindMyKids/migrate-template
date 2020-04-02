package template

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/golang-migrate/migrate/v4/source"
)

type driver struct {
	source.Driver
	funcs template.FuncMap
	vars  M
}

func (d *driver) parse(r io.ReadCloser) (io.ReadCloser, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("").Funcs(d.funcs).Parse(string(data))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, d.vars); err != nil {
		return nil, err
	}

	return ioutil.NopCloser(buf), nil
}

func (d *driver) Open(url string) (source.Driver, error) {
	inner, err := d.Driver.Open(url)
	if err != nil {
		return nil, err
	}
	return &driver{Driver: inner}, nil
}

func (d *driver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = d.Driver.ReadUp(version)
	if err == nil {
		r, err = d.parse(r)
	}
	return
}

func (d *driver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = d.Driver.ReadDown(version)
	if err == nil {
		r, err = d.parse(r)
	}
	return
}

type Option func(d *driver)

type M map[string]interface{}

func WithVars(vars M) Option {
	return func(d *driver) {
		for k := range vars {
			d.vars[k] = vars[k]
		}
	}
}

func WithFuncs(funcs template.FuncMap) Option {
	return func(d *driver) {
		for k := range funcs {
			d.funcs[k] = funcs[k]
		}
	}
}

func Wrap(inner source.Driver, opts ...Option) source.Driver {
	d := &driver{
		Driver: inner,
		funcs:  make(template.FuncMap),
		vars:   make(M),
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

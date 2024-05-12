package funcs

import (
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/influxdata/go-prompt"
)

type Funcs struct {
	promptOptions []prompt.Option
}

func NewFuncs() template.FuncMap {
	var funcs *Funcs
	return funcs.Funcs()
}

func (f *Funcs) Funcs() template.FuncMap {
	m := sprig.TxtFuncMap()
	m["i"] = Input
	m["f"] = File
	return m
}

// type Option func(Funcs)

// func WithGlobalPromptOption(opt prompt.Option) Option {
// 	return func(f Funcs) {
// 		f.promptOptions = append(f.promptOptions, opt)
// 	}
// }

func prompter(prefix string, completer prompt.Completer, opts ...prompt.Option) string {
	return prompt.Input(prefix, completer, opts...)
}

func Input(name string) string {
	return prompt.Input(name+": ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionPrefixTextColor(prompt.Green))
}

func File() string {
	return prompt.Input("file: ", func(d prompt.Document) []prompt.Suggest {
		dir := filepath.Dir(d.Text)
		q := d.Text
		r := []prompt.Suggest{}
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(path, q) {
				r = append(r, prompt.Suggest{
					Text: path,
				})
			}
			if len(r) >= 10 {
				return filepath.SkipAll
			}
			return nil
		})
		return r
	})
}

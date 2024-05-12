package main

import (
	"bytes"
	"debug/buildinfo"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/influxdata/go-prompt"
	"github.com/spf13/cobra"
)

func mainInternal() error {
	//nolint:wrapcheck
	return NewApp().Execute()
}

func main() {
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}

func NewApp() *cobra.Command {
	var path string
	c := cobra.Command{
		Use:   "tpl",
		Short: "template command",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := newTemplate(path)
			if err != nil {
				return err
			}
			var b bytes.Buffer
			err = t.Execute(&b)
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, &b)
			return err
		},
	}
	c.Flags().StringVarP(&path, "filepath", "f", "", "template file")

	c.AddCommand(
		NewFuncCommand(),
		NewVersionCommand(),
	)
	return &c
}

func NewFuncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "funcs",
		Short: "show help for template functions",
		Run: func(cmd *cobra.Command, args []string) {
			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "We are using github.com/Masterminds/sprig")
			fmt.Fprintln(w, "you can see document here, https://masterminds.github.io/sprig/")
			fmt.Fprintln(w)
			fmt.Fprintln(w, "We introduce some functions for template")

			// TODO
		},
	}
}

func NewVersionCommand() *cobra.Command {
	var detail bool
	c := &cobra.Command{
		Use:   "version",
		Short: "show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			w := cmd.OutOrStdout()
			info, err := buildinfo.ReadFile(os.Args[0])
			if err != nil {
				return fmt.Errorf("Cannot read buildinfo: %w", err)
			}

			fmt.Fprintf(w, "go version: %s", info.GoVersion)
			fmt.Fprintf(w, "module version: %s", info.Main.Version)
			if detail {
				fmt.Fprintln(w)
				fmt.Fprintln(w, info)
			}
			return nil
		},
	}
	c.Flags().BoolVarP(&detail, "detail", "d", false, "show details")
	return c
}

type Template struct {
	template *template.Template
}

func input(name string) string {
	return prompt.Input(name+": ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionPrefixTextColor(prompt.Green))
}

const maxSuggest = 10

func file() string {
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
			if len(r) >= maxSuggest {
				return filepath.SkipAll
			}
			return nil
		})
		return r
	}, prompt.OptionPrefixTextColor(prompt.Green))
}

func choose(name string, opts ...string) string {
	return prompt.Choose(name+": ", opts, prompt.OptionPrefixTextColor(prompt.Green))
}

func funcs() template.FuncMap {
	m := sprig.TxtFuncMap()
	m["i"] = input
	m["input"] = input
	m["f"] = file
	m["file"] = file
	m["choose"] = choose
	m["select"] = choose
	return m
}

func newTemplate(path string) (*Template, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file path:%s: %w", path, err)
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file path:%s: %w", path, err)
	}
	t, err := template.New("template").Funcs(funcs()).Parse(string(bs))
	if err != nil {
		return nil, fmt.Errorf("cannot parse template path:%s: %w", path, err)
	}
	return &Template{template: t}, nil
}

func (t *Template) Execute(w io.Writer) error {
	err := t.template.Execute(w, nil)
	if err != nil {
		return fmt.Errorf("unable to generate content: %w", err)
	}
	return nil
}

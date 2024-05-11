package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func mainInternal() error {
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
			t, err := NewTemplate(path)
			if err != nil {
				return err
			}
			return t.Execute(os.Stdout)
		},
	}
	c.Flags().StringVarP(&path, "filepath", "f", "", "template file")

	return &c
}

type Template struct {
	template *template.Template
}

func input(name string) string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . | green }}: ",
		Invalid: "{{ . | red }}: ",
		Success: "{{ . | bold }}: ",
	}

	prompt := promptui.Prompt{
		Label:     name,
		Templates: templates,
	}
	s, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func Funcs() template.FuncMap {
	m := sprig.TxtFuncMap()
	m["i"] = input
	return m
}

func NewTemplate(path string) (*Template, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(f)
	if err != nil && err != io.EOF {
		return nil, err
	}
	t, err := template.New("template").Funcs(Funcs()).Parse(string(bs))
	if err != nil {
		return nil, err
	}
	return &Template{template: t}, nil
}

func (t *Template) Execute(w io.Writer) error {
	return t.template.Execute(w, nil)
}

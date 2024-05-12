# tpl

Tiny interactive template generator using go text/template.

## Usage

```bash
$ cat template.tpl
Hello {{ input "name"}} !!
You can select {{ file }}.

$ tpl -f template.tpl
> name: John Doe
> file: testdata/example.tpl

Hello John Doe !!
You can select testdata/example.tpl.
```

## TODO

* Add functions to edit filepath
* Add help for functions
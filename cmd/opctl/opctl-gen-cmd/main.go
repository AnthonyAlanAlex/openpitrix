// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	. "openpitrix.io/openpitrix/cmd/opctl/common"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

const AllCmdTmpl = `
// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"openpitrix.io/openpitrix/test"
{{- range $index, $element := .services}}
	"openpitrix.io/openpitrix/test/client/{{snakeCase $element}}"
{{- end}}
	"openpitrix.io/openpitrix/test/models"
)

var AllCmd = []Cmd{
{{range $index, $element := .cmds}}	New{{$element.Action}}Cmd(),
{{end}}
}
`
const PreActionTmpl = `
{{- range $index, $element := .cmds}}
{{- if (gt (len $element.Body) 0)}}
{{/*   this if post action   */}}
type {{$element.Action}}Cmd struct {
	*models.Openpitrix{{$element.Action}}Request
}

func New{{$element.Action}}Cmd() Cmd {
	return &{{$element.Action}}Cmd{
		&models.Openpitrix{{$element.Action}}Request{},
	}
}

func (*{{$element.Action}}Cmd) GetActionName() string {
	return "{{$element.Action}}"
}

func (c *{{$element.Action}}Cmd) ParseFlag(f Flag) {
{{- range $name, $p := $element.Body}}
{{- if (eq $p.Type "string")}}
	f.StringVarP(&c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", "", "{{$p.Help}}")
{{- else if (eq $p.Type "[]string")}}
	f.StringSliceVarP(&c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", []string{}, "{{$p.Help}}")
{{- else if (eq $p.Type "boolean")}}
	f.BoolVarP(&c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", false, "{{$p.Help}}")
{{- end}}
{{- end}}
}

func (c *{{$element.Action}}Cmd) Run(out Out) error {
	params := {{snakeCase $element.Service}}.New{{$element.Action}}Params()
	params.WithBody(c.Openpitrix{{$element.Action}}Request)

	out.WriteRequest(params)

	client := test.GetClient(clientConfig)
	res, err := client.{{$element.Service}}.{{$element.Action}}(params)
	if err != nil {
		return err
	}

	out.WriteResponse(res.Payload)

	return nil
}
{{- else}}
{{/*   this if get action   */}}
type {{$element.Action}}Cmd struct {
	*{{snakeCase $element.Service}}.{{$element.Action}}Params
}

func New{{$element.Action}}Cmd() Cmd {
	return &{{$element.Action}}Cmd{
		{{snakeCase $element.Service}}.New{{$element.Action}}Params(),
	}
}

func (*{{$element.Action}}Cmd) GetActionName() string {
	return "{{$element.Action}}"
}

func (c *{{$element.Action}}Cmd) ParseFlag(f Flag) {
{{- range $name, $p := $element.Query}}
{{- if (eq $p.Type "string")}}
	c.{{pascalCase $name}} = new(string)
	f.StringVarP(c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", "", "{{$p.Help}}")
{{- else if (eq $p.Type "[]string")}}
	f.StringSliceVarP(&c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", []string{}, "{{$p.Help}}")
{{- else if (eq $p.Type "boolean")}}
	c.{{pascalCase $name}} = new(bool)
	f.BoolVarP(c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", false, "{{$p.Help}}")
{{- else if (eq $p.Type "integer")}}
	c.{{pascalCase $name}} = new(int64)
{{- if (eq $name "limit")}}
	f.Int64VarP(c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", 20, "{{$p.Help}}")
{{- else}}
	f.Int64VarP(c.{{pascalCase $name}}, "{{$name}}", "{{$p.Shorthand}}", 0, "{{$p.Help}}")
{{- end}}
{{- end}}
{{- end}}
}

func (c *{{$element.Action}}Cmd) Run(out Out) error {

	out.WriteRequest(c.{{$element.Action}}Params)

	client := test.GetClient(clientConfig)
	res, err := client.{{$element.Service}}.{{$element.Action}}(c.{{$element.Action}}Params)
	if err != nil {
		return err
	}

	out.WriteResponse(res.Payload)

	return nil
}
{{- end}}
{{- end}}
`

type Generator struct {
	Cmds
	buf  *bytes.Buffer
	data Data
}

type Data map[string]interface{}

func getTmpl(name, tmpl string) *template.Template {
	t, err := template.New(name).Funcs(template.FuncMap{
		"snakeCase": stringutil.CamelCaseToUnderscore,
		"pascalCase": func(str string) string {
			str = stringutil.UnderscoreToCamelCase(str)
			str = strings.Replace(str, "Id", "ID", -1)
			str = strings.Replace(str, "Url", "URL", -1)
			return str
		},
	}).Parse(tmpl)
	if err != nil {
		Error(err, fmt.Sprintf("parse template [%s]", name))
	}
	return t
}

func (g Generator) generateAllCmd() {
	err := getTmpl("all_cmd", AllCmdTmpl).Execute(g.buf, g.data)
	if err != nil {
		Error(err, "render template [all_cmd]")
	}
}

func (g Generator) generatePreAction() {
	err := getTmpl("pre_action", PreActionTmpl).Execute(g.buf, g.data)
	if err != nil {
		Error(err, "render template [pre_action]")
	}
}

func (g Generator) Generate() {
	g.buf = bytes.NewBufferString("")
	var services []string
	for _, cmd := range g.Cmds {
		services = append(services, cmd.Service)
	}
	g.data = Data{
		"cmds":     g.Cmds,
		"services": stringutil.Unique(services),
	}

	g.generateAllCmd()
	g.generatePreAction()

	g.buf.WriteTo(os.Stdout)
}

func main() {
	var filePath string
	flag.StringVarP(&filePath, "file", "f", "", "")
	flag.Parse()

	var content []byte
	var err error
	if filePath != "" {
		content, err = ioutil.ReadFile(filePath)
		if err != nil {
			Error(err, fmt.Sprintf("read file [%s]", filePath))
		}
	} else {
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			Error(err, "read stdin")
		}
	}
	var cmds Cmds
	err = yaml.Unmarshal(content, &cmds)
	if err != nil {
		Error(err, "unmarshal yaml")
	}
	g := Generator{cmds, nil, nil}
	g.Generate()
}

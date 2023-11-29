//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type PathPackage struct {
	Path    string
	Package string
}

var tmpl = `// Code generated by go generate; DO NOT EDIT.

package {{.Package}}

import (
	"context"
	"sync"
)

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (d *Component) Start(ctx context.Context, wg *sync.WaitGroup) {
	Start(ctx, wg)
}

func (d *Component) Name() string {
	return "{{.Package}}"
}
`

func GenerateComponents(components []PathPackage) error {
	fmt.Println("[*] Generating components ...")

	t := template.Must(template.New("component").Parse(tmpl))

	for _, component := range components {
		path := filepath.Join(component.Path, "component_gen.go")

		fmt.Printf("[+]    %s ...\n", path)

		if err := func() error {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()

			return t.Execute(f, component)
		}(); err != nil {
			return err
		}
	}

	return exec.Command("go", "fmt", "./...").Run()
}

func main() {
	components := []PathPackage{
		{"health", "health"},
		{"services/anonymizer", "anonymizer"},
	}

	err := GenerateComponents(components)
	if err != nil {
		panic(err)
	}
}
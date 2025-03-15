package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

func main() {
	fmt.Println("Enter the name of the domain object you want to create:")
	var name string
	fmt.Scanln(&name)

	check := isTitleCamelCase(name)
	for !check {
		fmt.Println()
		fmt.Println("Domain object names must be in title camel case, for example: MyDomainObject")
		fmt.Println("Enter the name of the domain object you want to create:")
		fmt.Scanln(&name)
		check = isTitleCamelCase(name)
	}

	fmt.Println()
	fmt.Println("Creating domain object...")
	fmt.Println()

	t := template.Must(template.New("domain").Parse(standardDomainTemplate))
	var tpl bytes.Buffer
	t.Execute(&tpl, TemplateData{Name: name})

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("%s.go", name)
	filePath := fmt.Sprintf("%s/%s/%s", currentDir, DomainDir, fileName)
	err = os.WriteFile(filePath, tpl.Bytes(), 0644)
	if err != nil {
		panic(err)
	} else {
		fmt.Println()
		fmt.Println("Domain object created successfully!")
		fmt.Println("Path:", filePath)
		fmt.Println()
	}
}

func isTitleCamelCase(s string) bool {
	re := regexp.MustCompile(`^[A-Z][a-z]*([A-Z][a-z]*)*$`)
	return re.MatchString(s)
}

const (
	DomainDir = "domain"
)

type TemplateData struct {
	Name string
}

const standardDomainTemplate = `package domain

import (
    "errors"
    "encoding/json"

    "github.com/coopersmall/subswag/utils"
)

type {{.Name}}ID utils.ID

type {{.Name}} struct {
    ID {{.Name}}ID ` + "`json:\"id\"`" + `
    {{.Name}}Data
    Metadata ` + "`json:\"metadata\"`" + `
}

type {{.Name}}Data struct {
    // Add fields here
}

func (obj {{.Name}}ID) Validate() error {
    if obj == 0 {
        return errors.New("{{.Name}}ID is required")
    }
    return nil
}

func (obj *{{.Name}}Data) Validate() error {
    // Add validation here
    return nil
}

func (obj *{{.Name}}) Validate() error {
    if err := obj.ID.Validate(); err != nil {
        return err
    }
    if err := obj.{{.Name}}Data.Validate(); err != nil {
        return err
    }
    if err := obj.Metadata.Validate(); err != nil {
        return err
    }
    return nil
}

func (obj *{{.Name}}) MarshalJSON() ([]byte, error) {
    type Alias {{.Name}}
    return json.Marshal(&struct {
        Alias
    }{
        Alias: (Alias)(*obj),
    })
}

func (obj *{{.Name}}) UnmarshalJSON(data []byte) error {
    type Alias {{.Name}}
    aux := &struct {
        *Alias
    }{
        Alias: (*Alias)(obj),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    *obj = {{.Name}}(*aux.Alias)
    if err := obj.Validate(); err != nil {
        return err
    }
    return nil
}

func (obj *{{.Name}}Data) MarshalJSON() ([]byte, error) {
    type Alias {{.Name}}Data
    return json.Marshal(&struct {
        Alias
    }{
        Alias: (Alias)(*obj),
    })
}

func (obj *{{.Name}}Data) UnmarshalJSON(data []byte) error {
    type Alias {{.Name}}Data
    aux := &struct {
        *Alias
    }{
        Alias: (*Alias)(obj),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    obj = (*{{.Name}}Data)(aux.Alias)
    if err := obj.Validate(); err != nil {
        return err
    }
    return nil
}
`

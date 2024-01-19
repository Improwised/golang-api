package worker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Handler interface {
	Handle(data []byte)
}

func (w DeleteUser) Handle(data []byte) {
	fmt.Println("Delete User")
	fmt.Println(w)
}
func (w UpdateUser) Handle(data []byte) {
	fmt.Println("Update User")
	fmt.Println(w)
}

func (w AddUser) Handle(data []byte) {
	fmt.Println("add user")
	fmt.Println(w)
}

func Process(msg *message.Message) error {
	// GobRegister()
	buf := bytes.NewBuffer(msg.Payload)
	dec := gob.NewDecoder(buf)

	var result Handler
	err := dec.Decode(&result)
	if err != nil {
		fmt.Println("Error in decoding")
		return err
	}
	result.Handle(msg.Payload)
	return nil
}

func init() {
	GenerateGobCode()
}

// -------------------------------------------------------------------------------------
// generate gob_register.go file

const outputTemplate = `package {{.PackageName}}
import (
    "encoding/gob"
	"fmt"
)
func GobRegister() {
    {{range .StructNames}}gob.Register({{.}}{})
	fmt.Println("Gob Register")
    {{end}}
}
`

func GenerateGobCode() {
	fset := token.NewFileSet()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("cannot get current file")
	}
	dir := filepath.Dir(filename)
	node, err := parser.ParseFile(fset, dir+"/struct.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}

	structNames := []string{}
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structNames = append(structNames, typeSpec.Name.Name)
		return false
	})

	outputFile, err := os.Create(dir + "/gob_register.go")
	if err != nil {
		fmt.Println(err)

	}

	defer outputFile.Close()

	t := template.Must(template.New("").Parse(outputTemplate))
	t.Execute(outputFile, map[string]interface{}{
		"PackageName": node.Name.Name,
		"StructNames": structNames,
	})

}

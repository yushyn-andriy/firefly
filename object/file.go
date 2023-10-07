package object

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type File struct {
	dict map[string]Object
	file *os.File
	path string
	mode string
}

func NewFile(path string, mode string) *File {
	f := new(File)
	f.path = path
	f.mode = mode
	f.dict = make(map[string]Object)

	f.initialize()

	return f
}

func (f *File) initialize() {
	f.dict["open"] = &Builtin{
		Fn:   fileOpen,
		Env:  nil,
		Self: f,
		Doc: `open(self)
open file in read/write file mode
`,
	}
	f.dict["close"] = &Builtin{
		Fn:   fileClose,
		Env:  nil,
		Self: f,
		Doc: `close(self)
`,
	}
	f.dict["read"] = &Builtin{
		Fn:   fileRead,
		Env:  nil,
		Self: f,
		Doc: `read(self)
`,
	}
	f.dict["write"] = &Builtin{
		Fn:   fileWrite,
		Env:  nil,
		Self: f,
		Doc: `write(self)
`,
	}
}

func (f *File) Type() ObjectType { return FILE_OBJ }
func (f *File) Inspect() string  { return fmt.Sprintf("%#v", f) }
func (f *File) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", f.Inspect(), key)}
}
func (f *File) GetAttr(key string) Object {
	if val, ok := f.dict[key]; ok {
		return val
	}
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", f.Inspect(), key)}
}

func fileOpen(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, ok := args[0].(*File)
	if !ok {
		return newError("%s", "cannot convert to type File")
	}

	var flag int = os.O_APPEND | os.O_CREATE
	switch self.mode {
	case "w":
		flag = flag | os.O_WRONLY
	case "r":
		flag = flag | os.O_RDONLY
	default:
		return newError("undefined mode '%s'", self.mode)
	}

	file, err := os.OpenFile(self.path, flag, 0777)
	if err != nil {
		return newError("%s", err)
	}
	self.file = file

	return self
}

func fileClose(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, ok := args[0].(*File)
	if !ok {
		return newError("%s", "cannot convert to type File")
	}

	if self.file == nil {
		return newError("%s", "cannot close not opened file")
	}

	err := self.file.Close()
	if err != nil {
		log.Printf("error while closing the file. %v", err)
		return newError("%s", err)
	}
	return NewNull()
}

func fileRead(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, ok := args[0].(*File)
	if !ok {
		return newError("%s", "cannot convert to type File")
	}

	if self.file == nil {
		return newError("%s", "cannot read from not opened file")
	}

	info, err := self.file.Stat()
	if err != nil {
		return newError("%s", err)
	}

	buffer := make([]byte, info.Size())
	_, err = self.file.Read(buffer)
	if err != nil {
		return newError("%s", err)
	}

	return NewString(string(buffer))
}

func fileWrite(env *Environment, args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, ok := args[0].(*File)
	if !ok {
		return newError("%s", "cannot convert to type File")
	}

	obj := args[1]
	// s, ok := args[1].(*String)
	// if !ok {
	// 	return newError("%s", "cannot convert to type String")
	// }

	if self.file == nil {
		return newError("%s", "cannot write to not opened file")
	}

	n, err := self.file.WriteString(
		fmt.Sprint(strings.ReplaceAll(obj.Inspect(), "\\n", "\n")))
	if err != nil {
		return newError("%s", err)
	}

	return NewInteger(int64(n))
}

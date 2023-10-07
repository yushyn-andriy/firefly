package object

import (
	"fmt"
	"log"
	"os"
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
	f.dict["fwrite"] = &Builtin{
		Fn:   fileFormatWrite,
		Env:  nil,
		Self: f,
		Doc: `fwrite(self)
`,
	}
	f.dict["writeln"] = &Builtin{
		Fn:   fileWriteLn,
		Env:  nil,
		Self: f,
		Doc: `fwriteln(self)
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

func (f *File) Close() {
	if f.file != nil {
		f.file.Close()
	}
}

func (f *File) Open() {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Printf("error while opening the file. %v", err)
		return
	}
	f.file = file
}

func (f *File) Write(s string) {
	if f.file != nil {
		f.file.Write([]byte(s))
	}
}

func (f *File) Fwrite(format, s string) {
	if f.file != nil {
		f.file.WriteString(fmt.Sprintf(format, s))
	}
}

func (f *File) Read() string {
	if f.file != nil {
		info, _ := f.file.Stat()
		buffer := make([]byte, info.Size())
		f.file.Read(buffer)
		return string(buffer)
	}
	return ""
}

func fileOpen(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	self.Open()
	return self
}

func fileClose(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	self.Close()
	return self
}

func fileRead(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	return NewString(self.Read())
}

func fileWrite(env *Environment, args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	s, _ := args[1].(*String)

	self.Fwrite("%s", s.Value)

	return self
}

func fileWriteLn(env *Environment, args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	s, _ := args[1].(*String)

	self.Fwrite("%s\n", s.Value)

	return self
}

func fileFormatWrite(env *Environment, args ...Object) Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*File)
	format, _ := args[1].(*String)
	s, _ := args[2].(*String)

	self.Fwrite(format.Value, s.Value)

	return self
}

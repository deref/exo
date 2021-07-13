package generate

type Module struct {
	Interfaces []Interface
	Structs    []Struct
}

type Interface struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name         string
	InputFields  []Field
	OutputFields []Field
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name     string
	Type     string
	Required bool
	Nullable bool
}

//func Validate(mod Module) error {
//	kinds := make(map[string]string)
//	for _, iface := range mod.Interfaces {
//	}
//}
//
//func GenerateTypes(w io.Writer, mod Module) error {
//}
//
//func GenerateClient(w io.Writer, mod Module, interfaceName string) error {
//}

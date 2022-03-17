package model

import "fmt"

type Package struct {
	path    string
	members []any
	names   map[string]Named
	errors  []error
}

type Named interface {
	Name() string
}

func NewPackage(path string) *Package {
	return &Package{
		path:  path,
		names: make(map[string]Named),
	}
}

func (pkg *Package) Path() string {
	return pkg.path
}

func (pkg *Package) Err() error {
	if len(pkg.errors) > 0 {
		return pkg.errors[0]
	}
	return nil
}

func (pkg *Package) DeclareInterface(name string) *Interface {
	named, ok := pkg.names[name]
	if named == nil {
		named = newInterface(pkg, name)
		pkg.names[name] = named
	}
	iface, ok := named.(*Interface)
	if !ok {
		pkg.AddError(fmt.Errorf("conflicting declarations for %q", name))
		iface = newInterface(pkg, name)
	}
	pkg.members = append(pkg.members, iface)
	return iface
}

func (pkg *Package) ReferInterface(name string) *Interface {
	named := pkg.names[name]
	iface, ok := named.(*Interface)
	if !ok {
		if named == nil {
			pkg.AddError(fmt.Errorf("undefined interface: %q", name))
		} else {
			pkg.AddError(fmt.Errorf("expected %q to be an interface", name))
		}
		iface = newInterface(pkg, name)
	}
	return iface
}

func (pkg *Package) DeclareStruct(name string) *Struct {
	named := pkg.names[name]
	if named == nil {
		named = newStruct(pkg, name)
		pkg.names[name] = named
	}
	strct, ok := named.(*Struct)
	if !ok {
		pkg.AddError(fmt.Errorf("conflicting declarations for %q", name))
		strct = newStruct(pkg, name)
	}
	pkg.members = append(pkg.members, strct)
	return strct
}

func (pkg *Package) ReferStruct(name string) *Struct {
	named := pkg.names[name]
	strct, ok := named.(*Struct)
	if !ok {
		if named == nil {
			pkg.AddError(fmt.Errorf("undefined struct: %q", name))
		} else {
			pkg.AddError(fmt.Errorf("expected %q to be a struct", name))
		}
		strct = newStruct(pkg, name)
	}
	return strct
}

func (pkg *Package) Interfaces() []*Interface {
	ifaces := make([]*Interface, 0, len(pkg.members))
	for _, named := range pkg.members {
		if iface, ok := named.(*Interface); ok {
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

func (pkg *Package) AddError(err error) {
	pkg.errors = append(pkg.errors, err)
}

func (pkg *Package) Structs() []*Struct {
	structs := make([]*Struct, 0, len(pkg.members))
	for _, named := range pkg.members {
		if strct, ok := named.(*Struct); ok {
			structs = append(structs, strct)
		}
	}
	return structs
}

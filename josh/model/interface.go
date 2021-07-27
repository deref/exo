package model

import "fmt"

type Interface struct {
	pkg     *Package
	name    string
	doc     string
	extends []*Interface
	methods []*Method
}

func newInterface(pkg *Package, name string) *Interface {
	return &Interface{
		pkg:  pkg,
		name: name,
	}
}

func (iface *Interface) AddError(err error) {
	iface.pkg.AddError(fmt.Errorf("interface %q: %w", iface.name, err))
}

func (iface *Interface) Name() string {
	return iface.name
}

func (iface *Interface) SetDoc(value string) {
	iface.doc = value
}

func (iface *Interface) Doc() string {
	return iface.doc
}

func (iface *Interface) Extends() []*Interface {
	extends := make([]*Interface, len(iface.extends))
	copy(extends, iface.extends)
	return extends
}

func (iface *Interface) Extend(extended *Interface) {
	iface.extends = append(iface.extends, extended)
}

func (iface *Interface) Methods() []*Method {
	methods := make([]*Method, len(iface.methods))
	copy(methods, iface.methods)
	return methods
}

func (iface *Interface) DeclareMethod(name string) *Method {
	method := &Method{
		iface: iface,
		name:  name,
	}
	iface.methods = append(iface.methods, method)
	return method
}

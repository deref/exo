package idl

import (
	"github.com/deref/exo/josh/model"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseFile(filename string) (*Unit, error) {
	var unit Unit
	if err := hclsimple.DecodeFile(filename, nil, &unit); err != nil {
		return nil, err
	}
	return &unit, nil
}

func LoadFile(pkg *model.Package, filePath string) {
	unit, err := ParseFile(filePath)
	if err != nil {
		pkg.AddError(err)
		return
	}

	// Declaration pass.
	for _, ifaceNode := range unit.Interfaces {
		pkg.DeclareInterface(ifaceNode.Name)
	}
	for _, structNode := range unit.Structs {
		pkg.DeclareStruct(structNode.Name)
	}

	// Definition pass.
	for _, ifaceNode := range unit.Interfaces {
		iface := pkg.ReferInterface(ifaceNode.Name)
		if ifaceNode.Doc != nil {
			iface.SetDoc(*ifaceNode.Doc)
		}
		for _, extendsName := range ifaceNode.Extends {
			extended := pkg.ReferInterface(extendsName)
			iface.Extend(extended)
		}
		for _, methodNode := range ifaceNode.Methods {
			method := iface.DeclareMethod(methodNode.Name)
			method.SetDoc(methodNode.Doc)
			for _, inputNode := range methodNode.Inputs {
				inputCfg := model.FieldConfig{
					Name: inputNode.Name,
					Type: inputNode.Type,
				}
				if inputNode.Doc != nil {
					inputCfg.Doc = *inputNode.Doc
				}
				if inputNode.Required != nil {
					inputCfg.Required = *inputNode.Required
				}
				if inputNode.Nullable != nil {
					inputCfg.Nullable = *inputNode.Nullable
				}
				method.AddInput(inputCfg)
			}
			for _, outputNode := range methodNode.Outputs {
				outputCfg := model.FieldConfig{
					Name: outputNode.Name,
					Type: outputNode.Type,
				}
				if outputNode.Doc != nil {
					outputCfg.Doc = *outputNode.Doc
				}
				if outputNode.Required != nil {
					outputCfg.Required = *outputNode.Required
				}
				if outputNode.Nullable != nil {
					outputCfg.Nullable = *outputNode.Nullable
				}
				method.AddOutput(outputCfg)
			}
		}
	}
	for _, structNode := range unit.Structs {
		strct := pkg.ReferStruct(structNode.Name)
		if structNode.Doc != nil {
			strct.SetDoc(*structNode.Doc)
		}
		for _, fieldNode := range structNode.Fields {
			fieldCfg := model.FieldConfig{
				Name: fieldNode.Name,
				Type: fieldNode.Type,
			}
			if fieldNode.Doc != nil {
				fieldCfg.Doc = *fieldNode.Doc
			}
			if fieldNode.Required != nil {
				fieldCfg.Required = *fieldNode.Required
			}
			if fieldNode.Nullable != nil {
				fieldCfg.Nullable = *fieldNode.Nullable
			}
			strct.AddField(fieldCfg)
		}
	}
}

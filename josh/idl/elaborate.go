package idl

func Elaborate(unit *Unit) {
	for _, controller := range unit.Controllers {
		methods := make([]Method, len(controller.Methods))
		for methodIndex, method := range controller.Methods {
			shiftInputs := 3
			inputs := make([]Field, len(method.Inputs)+shiftInputs)
			inputs[0] = Field{
				Name: "id",
				Type: "string",
			}
			inputs[1] = Field{
				Name: "spec",
				Type: "string",
			}
			inputs[2] = Field{
				Name: "state",
				Type: "string",
			}
			for inputIndex, input := range method.Inputs {
				inputs[shiftInputs+inputIndex] = input
			}

			shiftOutputs := 1
			outputs := make([]Field, len(method.Outputs)+shiftOutputs)
			outputs[0] = Field{
				Name: "state",
				Type: "string",
			}
			for outputIndex, output := range method.Outputs {
				inputs[shiftOutputs+outputIndex] = output
			}

			methods[methodIndex] = Method{
				Name:    method.Name,
				Doc:     method.Doc,
				Inputs:  inputs,
				Outputs: outputs,
			}
		}
		unit.Interfaces = append(unit.Interfaces, Interface{
			Name:    controller.Name,
			Doc:     controller.Doc,
			Extends: append([]string{}, controller.Extends...),
			Methods: methods,
		})
	}
}

package install

import (
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/atom"
)

type Install struct {
	deviceIDAtom atom.Atom
}

func Get(deviceIDPath string) *Install {
	return &Install{
		deviceIDAtom: atom.NewFileAtom(deviceIDPath, atom.CodecString),
	}
}

func (i *Install) GetDeviceID() (string, error) {
	var deviceID string
	err := i.deviceIDAtom.Swap(&deviceID, func() error {
		// See XXX: [ATOM JSON CODING].
		if deviceID == "null" {
			deviceID = gensym.RandomBase32()
		}
		return nil
	})
	return deviceID, err
}

package compose

import "testing"

func TestBlockioYAML(t *testing.T) {
	testYAML(t, "block_io", `
device_read_bps:
  - path: /dev/sdb
    rate: 12mb
device_write_bps:
  - path: /dev/sdb
    rate: 1024k
device_read_iops:
  - path: /dev/sdb
    rate: 120
device_write_iops:
  - path: /dev/sdb
    rate: 30
weight: 300
weight_device:
  - path: /dev/sda
    weight: 400
`, BlkioConfig{
		DeviceReadBPS: []ThrottleDevice{
			{
				Path: MakeString("/dev/sdb"),
				Rate: Bytes{
					String:   MakeString("12mb"),
					Quantity: 12,
					Unit: ByteUnit{
						Scalar: 1024 * 1024,
						Suffix: "mb",
					},
				},
			},
		},
		DeviceReadIOPS: []ThrottleDevice{
			{
				Path: MakeString("/dev/sdb"),
				Rate: Bytes{
					String:   MakeInt(120).String,
					Quantity: 120,
				},
			},
		},
		DeviceWriteBPS: []ThrottleDevice{
			{
				Path: MakeString("/dev/sdb"),
				Rate: Bytes{
					String:   MakeString("1024k"),
					Quantity: 1024,
					Unit: ByteUnit{
						Scalar: 1024,
						Suffix: "k",
					},
				},
			},
		},
		DeviceWriteIOPS: []ThrottleDevice{
			{
				Path: MakeString("/dev/sdb"),
				Rate: Bytes{
					String:   MakeInt(30).String,
					Quantity: 30,
				},
			},
		},
		Weight: MakeInt(300),
		WeightDevice: []WeightDevice{
			{
				Path:   MakeString("/dev/sda"),
				Weight: MakeInt(400),
			},
		},
	})
}

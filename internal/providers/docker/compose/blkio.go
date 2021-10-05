package compose

type BlkioConfig struct {
	DeviceReadBPS   []ThrottleDevice `yaml:"device_read_bps,omitempty"`
	DeviceWriteBPS  []ThrottleDevice `yaml:"device_write_bps,omitempty"`
	DeviceReadIOPS  []ThrottleDevice `yaml:"device_read_iops,omitempty"`
	DeviceWriteIOPS []ThrottleDevice `yaml:"device_write_iops,omitempty"`
	Weight          uint16           `yaml:"weight,omitempty"`
	WeightDevice    []WeightDevice   `yaml:"weight_device,omitempty"`
}

type ThrottleDevice struct {
	Path string `yaml:"path,omitempty"`
	Rate Bytes  `yaml:"rate,omitempty"`
}

type WeightDevice struct {
	Path   string `yaml:"path,omitempty"`
	Weight uint16 `yaml:"weight,omitempty"`
}

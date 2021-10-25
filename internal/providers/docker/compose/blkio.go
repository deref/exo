package compose

type BlkioConfig struct {
	DeviceReadBPS   []ThrottleDevice `yaml:"device_read_bps,omitempty"`
	DeviceWriteBPS  []ThrottleDevice `yaml:"device_write_bps,omitempty"`
	DeviceReadIOPS  []ThrottleDevice `yaml:"device_read_iops,omitempty"`
	DeviceWriteIOPS []ThrottleDevice `yaml:"device_write_iops,omitempty"`
	Weight          Int              `yaml:"weight,omitempty"`
	WeightDevice    []WeightDevice   `yaml:"weight_device,omitempty"`
}

type ThrottleDevice struct {
	Path String `yaml:"path,omitempty"`
	Rate Bytes  `yaml:"rate,omitempty"`
}

type WeightDevice struct {
	Path   String `yaml:"path,omitempty"`
	Weight Int    `yaml:"weight,omitempty"`
}

func (blkio *BlkioConfig) Interpolate(env Environment) error {
	return interpolateStruct(blkio, env)
}

func (td *ThrottleDevice) Interpolate(env Environment) error {
	return interpolateStruct(td, env)
}

func (wd *WeightDevice) Interpolate(env Environment) error {
	return interpolateStruct(wd, env)
}

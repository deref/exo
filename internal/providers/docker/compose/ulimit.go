package compose

type Ulimits map[string]Ulimit

type Ulimit struct {
	Hard int64
	Soft int64
}

func (u *Ulimit) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n int64
	err := unmarshal(&n)
	if err == nil {
		u.Hard = n
		u.Soft = n
		return err
	}

	var m struct {
		Hard int64 `ymal:"hard,omitempty"`
		Soft int64 `ymal:"soft,omitempty"`
	}
	if err := unmarshal(&m); err != nil {
		return err
	}
	u.Hard = m.Hard
	u.Soft = m.Soft
	return nil
}

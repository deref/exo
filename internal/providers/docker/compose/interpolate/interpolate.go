package interpolate

import (
	"fmt"
	"strconv"
)

func Interpolate(v interface{}, env Environment) error {
	r := &interpolator{
		env: env,
	}
	_, err := r.interpolate(v)
	return err
}

type interpolator struct {
	env  Environment
	path []string
}

func (r *interpolator) errorf(format string, v ...interface{}) error {
	return fmt.Errorf("at %v: "+format, append([]interface{}{r.path}, v...)...)
}

func (r *interpolator) interpolate(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case bool, int, uint, int32, uint32, int64, uint64, float32, float64:
		return v, nil

	case string:
		tmpl, err := NewTemplate(v)
		if err != nil {
			return nil, r.errorf("compiling template: %w", err)
		}
		res, err := Substitute(tmpl, r.env)
		if err != nil {
			return nil, r.errorf("substituting: %w", err)
		}
		return res, nil

	case map[string]interface{}:
		n := len(r.path)
		r.path = append(r.path, "")
		for k, x := range v {
			r.path[n] = k
			y, err := r.interpolate(x)
			if err != nil {
				return nil, err
			}
			v[k] = y
		}
		r.path = r.path[:n]
		return v, nil

	case []interface{}:
		n := len(r.path)
		r.path = append(r.path, "")
		for i, x := range v {
			r.path[n] = strconv.Itoa(i)
			y, err := r.interpolate(x)
			if err != nil {
				return nil, err
			}
			v[i] = y
		}
		r.path = r.path[:n]
		return v, nil

	default:
		return nil, r.errorf("cannot interpolate value of type: %T", v)
	}
}

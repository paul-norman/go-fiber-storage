package storage

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
)

// Result struct for returning data back to its original types
type Result struct {
	Value any
	Error error
	Missed bool
}

// Return the result back as a boolean
func (r *Result) Bool() (bool, error, bool) {
	if r.Error != nil {
		return false, r.Error, r.Missed
	}

	if r.Missed {
		return false, nil, true
	}

	if r.Value == nil {
		return false, errors.New("bool values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			return value, nil, false
		case []byte:
			bool, err := strconv.ParseBool(string(value))
			if err == nil {
				return bool, nil, false
			}
			return false, errors.New("invalid bool value (from byte slice)"), false
		case string:
			bool, err := strconv.ParseBool(value)
			if err == nil {
				return bool, nil, false
			}
			return false, errors.New("invalid bool value (from string)"), false
		case int:
			return value != 0, nil, false
		case int8:
			return value != int8(0), nil, false
		case int16:
			return value != int16(0), nil, false
		case int32:
			return value != int32(0), nil, false
		case int64:
			return value != int64(0), nil, false
		case uint:
			return value != uint(0), nil, false
		case uint8:
			return value != uint8(0), nil, false
		case uint16:
			return value != uint16(0), nil, false
		case uint32:
			return value != uint32(0), nil, false
		case uint64:
			return value != uint64(0), nil, false
		case float32:
			return value != float32(0), nil, false
		case float64:
			return value != float64(0), nil, false
	}

	return false, errors.New("invalid bool value"), false
}

// Return the result back as a bool slice
func (r *Result) BoolSlice() ([]bool, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("bool slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []bool:
			return value, nil, false
	}

	return nil, errors.New("invalid bool slice value"), false
}

// Return the result back as a byte slice
func (r *Result) Bytes() ([]byte, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("byte slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []byte:
			return value, nil, false
		case string:
			return []byte(value), nil, false
	}

	return nil, errors.New("invalid byte slice value"), false
}

// Return the result back as a byte slice
func (r *Result) ByteSlice() ([]byte, error, bool) {
	return r.Bytes()
}

// Return the result error
func (r *Result) Err() error {
	return r.Error
}

// Return the result back as a 64-bit float
func (r *Result) Float() (float64, error, bool) {
	return r.Float64()
}

// Return the result back as a 32-bit float
func (r *Result) Float32() (float32, error, bool) {
	if r.Error != nil {
		return 0, r.Error, r.Missed
	}

	if r.Missed {
		return 0, nil, true
	}

	if r.Value == nil {
		return 0, errors.New("float32 values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return 1, nil, false
			}
			return 0, nil, false
		case []byte:
			float, err := strconv.ParseFloat(string(value), 32)
			if err != nil {
				return 0, errors.New("invalid float32 value (from byte slice)"), false
			}
			return float32(float), nil, false
		case string:
			float, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, errors.New("invalid float32 value (from string)"), false
			}
			return float32(float), nil, false
		case int:
			return float32(value), nil, false
		case int8:
			return float32(value), nil, false
		case int16:
			return float32(value), nil, false
		case int32:
			return float32(value), nil, false
		case int64:
			return float32(value), nil, false
		case uint:
			return float32(value), nil, false
		case uint8:
			return float32(value), nil, false
		case uint16:
			return float32(value), nil, false
		case uint32:
			return float32(value), nil, false
		case uint64:
			return float32(value), nil, false
		case float32:
			return value, nil, false
		case float64:
			return float32(value), nil, false
	}

	return 0, errors.New("invalid float32 value"), false
}

// Return the result back as a float32 slice
func (r *Result) Float32Slice() ([]float32, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("float32 slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []float32:
			return value, nil, false
	}

	return nil, errors.New("invalid float32 slice value"), false
}

// Return the result back as a 64-bit float
func (r *Result) Float64() (float64, error, bool) {
	if r.Error != nil {
		return 0, r.Error, r.Missed
	}

	if r.Missed {
		return 0, nil, true
	}

	if r.Value == nil {
		return 0, errors.New("float64 values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return 1, nil, false
			}
			return 0, nil, false
		case []byte:
			float, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return 0, errors.New("invalid float64 value (from byte slice)"), false
			}
			return float, nil, false
		case string:
			float, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, errors.New("invalid float64 value (from string)"), false
			}
			return float, nil, false
		case int:
			return float64(value), nil, false
		case int8:
			return float64(value), nil, false
		case int16:
			return float64(value), nil, false
		case int32:
			return float64(value), nil, false
		case int64:
			return float64(value), nil, false
		case uint:
			return float64(value), nil, false
		case uint8:
			return float64(value), nil, false
		case uint16:
			return float64(value), nil, false
		case uint32:
			return float64(value), nil, false
		case uint64:
			return float64(value), nil, false
		case float32:
			return float64(value), nil, false
		case float64:
			return value, nil, false
	}

	return 0, errors.New("invalid float64 value"), false
}

// Return the result back as a float64 slice
func (r *Result) Float64Slice() ([]float64, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("float64 slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []float64:
			return value, nil, false
	}

	return nil, errors.New("invalid float64 slice value"), false
}

// Return whether the result was present in the cache
func (r *Result) Hit() bool {
	return !r.Missed
}

// Return the result back as an integer
func (r *Result) Int() (int, error, bool) {
	if r.Error != nil {
		return 0, r.Error, r.Missed
	}

	if r.Missed {
		return 0, nil, true
	}

	if r.Value == nil {
		return 0, errors.New("int values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return 1, nil, false
			}
			return 0, nil, false
		case []byte:
			integer, err := strconv.ParseInt(string(value), 10, 64)
			if err != nil {
				return 0, errors.New("invalid int value (from byte slice)"), false
			}
			return int(integer), nil, false
		case string:
			integer, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, errors.New("invalid int value (from string)"), false
			}
			return int(integer), nil, false
		case int:
			return value, nil, false
		case int8:
			return int(value), nil, false
		case int16:
			return int(value), nil, false
		case int32:
			return int(value), nil, false
		case int64:
			return int(value), nil, false
		case uint:
			return int(value), nil, false
		case uint8:
			return int(value), nil, false
		case uint16:
			return int(value), nil, false
		case uint32:
			return int(value), nil, false
		case uint64:
			return int(value), nil, false
		case float32:
			return int(math.Round(float64(value))), nil, false
		case float64:
			return int(math.Round(value)), nil, false
	}

	return 0, errors.New("invalid int value"), false
}

// Return the result back as an int slice
func (r *Result) IntSlice() ([]int, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("int slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []bool:
			new := []int{}
			for _, v := range value {
				if v {
					new = append(new, 1)
				} else {
					new = append(new, 0)
				}
			}
			return new, nil, false
		case []int:
			return value, nil, false
		case []int8:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []int16:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []int32:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []int64:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []uint:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []uint8:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []uint16:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []uint32:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []uint64:
			new := []int{}
			for _, v := range value {
				new = append(new, int(v))
			}
			return new, nil, false
		case []float32:
			new := []int{}
			for _, v := range value {
				new = append(new, int(math.Round(float64(v))))
			}
			return new, nil, false
		case []float64:
			new := []int{}
			for _, v := range value {
				new = append(new, int(math.Round(v)))
			}
			return new, nil, false
	}

	return nil, errors.New("invalid int slice value"), false
}

// Return the result back as a 64-bit integer
func (r *Result) Int64() (int64, error, bool) {
	if r.Error != nil {
		return 0, r.Error, r.Missed
	}

	if r.Missed {
		return 0, nil, true
	}

	if r.Value == nil {
		return 0, errors.New("int64 values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return 1, nil, false
			}
			return 0, nil, false
		case []byte:
			integer, err := strconv.ParseInt(string(value), 10, 64)
			if err != nil {
				return 0, errors.New("invalid int64 value (from byte slice)"), false
			}
			return integer, nil, false
		case string:
			integer, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, errors.New("invalid int64 value (from string)"), false
			}
			return integer, nil, false
		case int:
			return int64(value), nil, false
		case int8:
			return int64(value), nil, false
		case int16:
			return int64(value), nil, false
		case int32:
			return int64(value), nil, false
		case int64:
			return value, nil, false
		case uint:
			return int64(value), nil, false
		case uint8:
			return int64(value), nil, false
		case uint16:
			return int64(value), nil, false
		case uint32:
			return int64(value), nil, false
		case uint64:
			return int64(value), nil, false
		case float32:
			return int64(math.Round(float64(value))), nil, false
		case float64:
			return int64(math.Round(value)), nil, false
	}

	return 0, errors.New("invalid int64 value"), false
}

// Return the result back as an int64 slice
func (r *Result) Int64Slice() ([]int64, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("int64 slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []bool:
			new := []int64{}
			for _, v := range value {
				if v {
					new = append(new, 1)
				} else {
					new = append(new, 0)
				}
			}
			return new, nil, false
		case []int:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []int8:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []int16:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []int32:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []int64:	
			return value, nil, false
		case []uint:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []uint8:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []uint16:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []uint32:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []uint64:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(v))
			}
			return new, nil, false
		case []float32:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(math.Round(float64(v))))
			}
			return new, nil, false
		case []float64:
			new := []int64{}
			for _, v := range value {
				new = append(new, int64(math.Round(v)))
			}
			return new, nil, false
	}

	return nil, errors.New("invalid int64 slice value"), false
}

// Return the result back as an interface{}
func (r *Result) Interface() (any, error, bool) {
	return r.Value, r.Error, r.Missed
}

// Return whether the result was not present in the cache
func (r *Result) Miss() bool {
	return r.Missed
}

// Return the result back as an interface{} (Redis naming)
func (r *Result) Result() (any, error, bool) {
	return r.Value, r.Error, r.Missed
}

// Return the result back as a string
func (r *Result) String() (string, error, bool) {
	if r.Error != nil {
		return "", r.Error, r.Missed
	}

	if r.Missed {
		return "", nil, true
	}

	if r.Value == nil {
		return "", errors.New("string values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return "true", nil, false
			}
			return "false", nil, false
		case []byte:
			return string(value), nil, false
		case string:
			return value, nil, false
		case int:
			return fmt.Sprintf("%d", value), nil, false
		case int8:
			return fmt.Sprintf("%d", value), nil, false
		case int16:
			return fmt.Sprintf("%d", value), nil, false
		case int32:
			return fmt.Sprintf("%d", value), nil, false
		case int64:
			return fmt.Sprintf("%d", value), nil, false
		case uint:
			return fmt.Sprintf("%d", value), nil, false
		case uint8:
			return fmt.Sprintf("%d", value), nil, false
		case uint16:
			return fmt.Sprintf("%d", value), nil, false
		case uint32:
			return fmt.Sprintf("%d", value), nil, false
		case uint64:
			return fmt.Sprintf("%d", value), nil, false
		case float32:
			return fmt.Sprintf("%f", value), nil, false
		case float64:
			return fmt.Sprintf("%f", value), nil, false
	}

	return "", errors.New("invalid string value"), false
}

// Return the result back as a string slice
func (r *Result) StringSlice() ([]string, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("string slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []string:
			return value, nil, false
	}

	return nil, errors.New("invalid string slice value"), false
}

// Return the result back as a string (Redis naming)
func (r *Result) Text() (string, error, bool) {
	return r.String()
}

// Return the result back as a string slice (Redis naming)
func (r *Result) TextSlice() ([]string, error, bool) {
	return r.StringSlice()
}

// Return the result back as an unsigned 64-bit integer
func (r *Result) Uint64() (uint64, error, bool) {
	if r.Error != nil {
		return 0, r.Error, r.Missed
	}

	if r.Missed {
		return 0, nil, true
	}

	if r.Value == nil {
		return 0, errors.New("uint64 values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case bool:
			if value {
				return 1, nil, false
			}
			return 0, nil, false
		case []byte:
			integer, err := strconv.ParseUint(string(value), 10, 64)
			if err != nil {
				return 0, errors.New("invalid uint64 value (from byte slice)"), false
			}
			return integer, nil, false
		case string:
			integer, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return 0, errors.New("invalid uint64 value (from string)"), false
			}
			return integer, nil, false
		case int:
			return uint64(value), nil, false
		case int8:
			return uint64(value), nil, false
		case int16:
			return uint64(value), nil, false
		case int32:
			return uint64(value), nil, false
		case int64:
			return uint64(value), nil, false
		case uint:
			return uint64(value), nil, false
		case uint8:
			return uint64(value), nil, false
		case uint16:
			return uint64(value), nil, false
		case uint32:
			return uint64(value), nil, false
		case uint64:
			return value, nil, false
		case float32:
			return uint64(value), nil, false
		case float64:
			return uint64(value), nil, false
	}

	return 0, errors.New("invalid uint64 value"), false
}

// Return the result back as a uint64 slice
func (r *Result) Uint64Slice() ([]uint64, error, bool) {
	if r.Error != nil {
		return nil, r.Error, r.Missed
	}

	if r.Missed {
		return nil, nil, true
	}

	if r.Value == nil {
		return nil, errors.New("uint64 slice values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []bool:
			new := []uint64{}
			for _, v := range value {
				if v {
					new = append(new, 1)
				} else {
					new = append(new, 0)
				}
			}
			return new, nil, false
		case []int:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []int8:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []int16:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []int32:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []int64:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []uint:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []uint8:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []uint16:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []uint32:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(v))
			}
			return new, nil, false
		case []uint64:
			return value, nil, false
		case []float32:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(math.Round(float64(v))))
			}
			return new, nil, false
		case []float64:
			new := []uint64{}
			for _, v := range value {
				new = append(new, uint64(math.Round(v)))
			}
			return new, nil, false
	}

	return nil, errors.New("invalid uint64 slice value"), false
}

// Return the result back as a UUID
func (r *Result) UUID() (uuid.UUID, error, bool) {
	if r.Error != nil {
		return uuid.Nil, r.Error, r.Missed
	}

	if r.Missed {
		return uuid.Nil, nil, true
	}

	if r.Value == nil {
		return uuid.Nil, errors.New("UUID values may not be nil"), false
	}

	switch value := r.Value.(type) {
		case []byte:
			token, err := uuid.ParseBytes(value)
			if err != nil {
				return uuid.Nil, errors.New("invalid UUID value (from byte slice)"), false
			}
			return token, nil, false
		case string:
			token, err := uuid.Parse(value)
			if err != nil {
				return uuid.Nil, errors.New("invalid UUID value (from string)"), false
			}
			return token, nil, false
	}

	return uuid.Nil, errors.New("invalid UUID value"), false
}

// Return the result value back as an interface{}
func (r *Result) Val() any {
	return r.Value
}
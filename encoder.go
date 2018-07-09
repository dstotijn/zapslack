package zapslack

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

// StringObjectEncoder implements zapcore.ObjectEncoder and is used for encoding
// fields to string values.
type StringObjectEncoder struct {
	Fields map[string]string
	ns     string
}

// NewStringObjectEncoder returns a new StringObjectEncoder.
func NewStringObjectEncoder() *StringObjectEncoder {
	f := make(map[string]string)

	return &StringObjectEncoder{Fields: f}
}

// AddArray implements ObjectEncoder.
func (enc *StringObjectEncoder) AddArray(key string, v zapcore.ArrayMarshaler) error {
	// arr := &sliceArrayEncoder{}
	// err := v.MarshalLogArray(arr)
	// enc.Fields[key] = arr.elems
	// return err
	return nil
}

// AddObject implements ObjectEncoder.
func (enc *StringObjectEncoder) AddObject(k string, v zapcore.ObjectMarshaler) error {
	// newMap := NewMapObjectEncoder()
	// enc.Fields[enc.keyNS(k)] = newMap.Fields
	// return v.MarshalLogObject(newMap)
	return nil
}

// AddBinary implements ObjectEncoder.
func (enc *StringObjectEncoder) AddBinary(k string, v []byte) {
	enc.Fields[enc.keyNS(k)] = base64.StdEncoding.EncodeToString(v)
}

// AddByteString implements ObjectEncoder.
func (enc *StringObjectEncoder) AddByteString(k string, v []byte) {
	enc.Fields[enc.keyNS(k)] = string(v)
}

// AddBool implements ObjectEncoder.
func (enc *StringObjectEncoder) AddBool(k string, v bool) {
	enc.Fields[enc.keyNS(k)] = strconv.FormatBool(v)
}

// AddDuration implements ObjectEncoder.
func (enc *StringObjectEncoder) AddDuration(k string, v time.Duration) {
	enc.Fields[enc.keyNS(k)] = v.String()
}

// AddComplex128 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddComplex128(k string, v complex128) {
	r, i := float64(real(v)), float64(imag(v))
	enc.Fields[enc.keyNS(k)] = fmt.Sprintf("%v+%vi",
		strconv.FormatFloat(r, 'f', -1, 64),
		strconv.FormatFloat(i, 'f', -1, 64),
	)
}

// AddComplex64 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddComplex64(k string, v complex64) {
	enc.AddComplex128(k, complex128(v))
}

// AddFloat64 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddFloat64(k string, v float64) {
	enc.Fields[enc.keyNS(k)] = strconv.FormatFloat(v, 'f', -1, 64)
}

// AddFloat32 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddFloat32(k string, v float32) {
	enc.Fields[enc.keyNS(k)] = strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// AddInt implements ObjectEncoder.
func (enc *StringObjectEncoder) AddInt(k string, v int) {
	enc.AddInt64(k, int64(v))
}

// AddInt64 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddInt64(k string, v int64) {
	enc.Fields[enc.keyNS(k)] = strconv.FormatInt(v, 10)
}

// AddInt32 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddInt32(k string, v int32) {
	enc.AddInt64(k, int64(v))
}

// AddInt16 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddInt16(k string, v int16) {
	enc.AddInt64(k, int64(v))
}

// AddInt8 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddInt8(k string, v int8) {
	enc.AddInt64(k, int64(v))
}

// AddString implements ObjectEncoder.
func (enc *StringObjectEncoder) AddString(k string, v string) { enc.Fields[enc.keyNS(k)] = v }

// AddTime implements ObjectEncoder.
func (enc *StringObjectEncoder) AddTime(k string, v time.Time) {
	enc.Fields[enc.keyNS(k)] = v.String()
}

// AddUint implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUint(k string, v uint) {
	enc.AddUint64(k, uint64(v))
}

// AddUint64 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUint64(k string, v uint64) {
	enc.Fields[enc.keyNS(k)] = strconv.FormatUint(v, 10)
}

// AddUint32 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUint32(k string, v uint32) {
	enc.AddUint64(k, uint64(v))
}

// AddUint16 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUint16(k string, v uint16) {
	enc.AddUint64(k, uint64(v))
}

// AddUint8 implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUint8(k string, v uint8) {
	enc.AddUint64(k, uint64(v))
}

// AddUintptr implements ObjectEncoder.
func (enc *StringObjectEncoder) AddUintptr(k string, v uintptr) {
	enc.AddUint64(k, uint64(v))
}

// AddReflected implements ObjectEncoder.
func (enc *StringObjectEncoder) AddReflected(k string, v interface{}) error {
	str, ok := v.(fmt.Stringer)
	if !ok {
		return errors.New("type cannot be formatted as string")
	}
	enc.Fields[enc.keyNS(k)] = str.String()

	return nil
}

// OpenNamespace implements ObjectEncoder.
func (enc *StringObjectEncoder) OpenNamespace(k string) {
	if enc.ns == "" {
		enc.ns = k
	} else {
		enc.ns += "." + k
	}
}

func (enc *StringObjectEncoder) keyNS(k string) string {
	if enc.ns == "" {
		return k
	}

	return enc.ns + "." + k
}

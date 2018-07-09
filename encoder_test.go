package zapslack

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

type foobar struct {
	s string
}

func (f foobar) String() string { return f.s }

func TestStringObjectEncoder(t *testing.T) {
	tests := []struct {
		desc     string
		f        func(zapcore.ObjectEncoder)
		expected string
	}{
		// {
		// 	desc: "AddObject",
		// 	f: func(e ObjectEncoder) {
		// 		assert.NoError(t, e.AddObject("k", loggable{true}), "Expected AddObject to succeed.")
		// 	},
		// 	expected: map[string]interface{}{"loggable": "yes"},
		// },
		// {
		// 	desc: "AddObject (nested)",
		// 	f: func(e ObjectEncoder) {
		// 		assert.NoError(t, e.AddObject("k", turducken{}), "Expected AddObject to succeed.")
		// 	},
		// 	expected: wantTurducken,
		// },
		// {
		// 	desc: "AddArray",
		// 	f: func(e ObjectEncoder) {
		// 		assert.NoError(t, e.AddArray("k", ArrayMarshalerFunc(func(arr ArrayEncoder) error {
		// 			arr.AppendBool(true)
		// 			arr.AppendBool(false)
		// 			arr.AppendBool(true)
		// 			return nil
		// 		})), "Expected AddArray to succeed.")
		// 	},
		// 	expected: []interface{}{true, false, true},
		// },
		// {
		// 	desc: "AddArray (nested)",
		// 	f: func(e ObjectEncoder) {
		// 		assert.NoError(t, e.AddArray("k", turduckens(2)), "Expected AddArray to succeed.")
		// 	},
		// 	expected: []interface{}{wantTurducken, wantTurducken},
		// },
		{
			desc:     "AddBinary",
			f:        func(e zapcore.ObjectEncoder) { e.AddBinary("k", []byte("foobar")) },
			expected: "Zm9vYmFy",
		},
		{
			desc:     "AddBool",
			f:        func(e zapcore.ObjectEncoder) { e.AddBool("k", true) },
			expected: "true",
		},
		{
			desc:     "AddComplex128",
			f:        func(e zapcore.ObjectEncoder) { e.AddComplex128("k", 1+2i) },
			expected: "1+2i",
		},
		{
			desc:     "AddComplex64",
			f:        func(e zapcore.ObjectEncoder) { e.AddComplex64("k", 1+2i) },
			expected: "1+2i",
		},
		{
			desc:     "AddDuration",
			f:        func(e zapcore.ObjectEncoder) { e.AddDuration("k", time.Millisecond) },
			expected: "1ms",
		},
		{
			desc:     "AddFloat64",
			f:        func(e zapcore.ObjectEncoder) { e.AddFloat64("k", 3.14) },
			expected: "3.14",
		},
		{
			desc:     "AddFloat32",
			f:        func(e zapcore.ObjectEncoder) { e.AddFloat32("k", 3.14) },
			expected: "3.14",
		},
		{
			desc:     "AddInt",
			f:        func(e zapcore.ObjectEncoder) { e.AddInt("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddInt64",
			f:        func(e zapcore.ObjectEncoder) { e.AddInt64("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddInt32",
			f:        func(e zapcore.ObjectEncoder) { e.AddInt32("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddInt16",
			f:        func(e zapcore.ObjectEncoder) { e.AddInt16("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddInt8",
			f:        func(e zapcore.ObjectEncoder) { e.AddInt8("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddString",
			f:        func(e zapcore.ObjectEncoder) { e.AddString("k", "v") },
			expected: "v",
		},
		{
			desc:     "AddTime",
			f:        func(e zapcore.ObjectEncoder) { e.AddTime("k", time.Unix(0, 100).UTC()) },
			expected: "1970-01-01 00:00:00.0000001 +0000 UTC",
		},
		{
			desc:     "AddUint",
			f:        func(e zapcore.ObjectEncoder) { e.AddUint("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddUint64",
			f:        func(e zapcore.ObjectEncoder) { e.AddUint64("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddUint32",
			f:        func(e zapcore.ObjectEncoder) { e.AddUint32("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddUint16",
			f:        func(e zapcore.ObjectEncoder) { e.AddUint16("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddUint8",
			f:        func(e zapcore.ObjectEncoder) { e.AddUint8("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddUintptr",
			f:        func(e zapcore.ObjectEncoder) { e.AddUintptr("k", 42) },
			expected: "42",
		},
		{
			desc:     "AddReflected",
			f:        func(e zapcore.ObjectEncoder) { e.AddReflected("k", foobar{"foo"}) },
			expected: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			enc := NewStringObjectEncoder()
			tt.f(enc)
			actual := enc.Fields["k"]

			if actual != tt.expected {
				t.Errorf("expected `%v`, got: `%v`", tt.expected, actual)
			}
		})
	}
}

func TestStringObjectEncoderNS(t *testing.T) {
	enc := NewStringObjectEncoder()
	enc.OpenNamespace("k")
	enc.AddInt("foo", 1)
	enc.OpenNamespace("middle")
	enc.AddInt("foo", 2)
	enc.OpenNamespace("inner")
	enc.AddInt("foo", 3)

	tests := map[string]string{
		"k.foo":              "1",
		"k.middle.foo":       "2",
		"k.middle.inner.foo": "3",
	}

	for k, v := range tests {
		if enc.Fields[k] != v {
			t.Errorf("enc.Fields[%v]: expected `%v`, got: `%v`", k, v, enc.Fields[k])
		}
	}
}

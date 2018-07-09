package zapslack

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCore(t *testing.T) {
	core := NewCore(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	}), nil)

	tests := []struct {
		desc          string
		field         zapcore.Field
		expectedKey   string
		expectedValue string
	}{
		{
			desc:          "zap.String",
			field:         zap.String("k", "v"),
			expectedKey:   "k",
			expectedValue: "v",
		},
		{
			desc:          "zap.Error",
			field:         zap.Error(errors.New("v")),
			expectedKey:   "error",
			expectedValue: "v",
		},
		{
			desc:          "zap.Any",
			field:         zap.Any("k", errors.New("v")),
			expectedKey:   "k",
			expectedValue: "v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			core.With([]zapcore.Field{tt.field})
			got := core.enc.Fields[tt.expectedKey]
			if got != tt.expectedValue {
				t.Errorf("expected `%v`, got: `%v`", tt.expectedValue, got)
			}
		})
	}
}

func TestEntry(t *testing.T) {
	core := NewCore(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	}), nil)
	logger := zap.New(core).Sugar()

	tests := []struct {
		desc          string
		f             func(*zap.SugaredLogger)
		expectedKey   string
		expectedValue string
	}{
		{
			desc: "zap.Error",
			f: func(logger *zap.SugaredLogger) {
				logger.Errorw("foo",
					"error", errors.New("bar"),
				)
			},
			expectedKey:   "error",
			expectedValue: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.f(logger)
			t.Log(core.enc.Fields)
			got := core.enc.Fields[tt.expectedKey]
			if got != tt.expectedValue {
				t.Errorf("expected `%v`, got: `%v`", tt.expectedValue, got)
			}
		})
	}
}

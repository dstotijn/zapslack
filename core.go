package zapslack

import (
	"github.com/dstotijn/slackhook"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLevelColors = map[zapcore.Level]string{
	zapcore.DebugLevel:  "good",
	zapcore.InfoLevel:   "good",
	zapcore.WarnLevel:   "warning",
	zapcore.ErrorLevel:  "danger",
	zapcore.DPanicLevel: "danger",
	zapcore.PanicLevel:  "danger",
	zapcore.FatalLevel:  "danger",
}

// Core implements the zapcore.Core interface for logging.
type Core struct {
	zapcore.LevelEnabler
	MessageFn MessageFunc
	client    *slackhook.Client
}

// NewCore returns a new zap Core.
func NewCore(enab zapcore.LevelEnabler, client *slackhook.Client) zapcore.Core {
	return &Core{
		LevelEnabler: enab,
		MessageFn:    defaultMessage,
		client:       client,
	}
}

// Check decides if an entry should be logged by this Core.
func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

// Sync does nothing, as sending webhooks to Slack is non-buffered.
func (c *Core) Sync() error {
	return nil
}

// With satisfies the zapcore.Core interface but does not alter anything.
func (c *Core) With(fields []zap.Field) zapcore.Core {
	clone := *c
	return &clone
}

// Write parses a log entry and sends it to Slack.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// We don't want to block when sending messages to Slack.
	// Downside is we don't return any errors, and of course we don't want to
	// log here because that could potentially lead to an infinite loop.
	// There is a risk that log entries will end up out of order at Slack.
	// TODO: Should we do queueing?
	go c.client.SendMessage(c.MessageFn(entry, fields))

	return nil
}

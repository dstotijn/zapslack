package zapslack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"go.uber.org/zap/zapcore"
)

// MessageFunc defines a type that parses a zap log entry into a Slack message.
type MessageFunc func(zapcore.Entry, []zapcore.Field) slack.Msg

func defaultMessage(entry zapcore.Entry, fields []zapcore.Field) slack.Msg {
	msgFields := make([]slack.AttachmentField, len(fields))

	for i, field := range fields {
		msgFields[i] = slack.AttachmentField{
			Title: strings.Title(field.Key),
			Value: field.String,
		}
	}

	msgFields = append(msgFields, slack.AttachmentField{
		Title: "Caller",
		Value: entry.Caller.TrimmedPath(),
	})

	if entry.Stack != "" {
		msgFields = append(msgFields, slack.AttachmentField{
			Title: "Stack trace",
			Value: "```" + entry.Stack + "```",
		})
	}

	return slack.Msg{
		Attachments: []slack.Attachment{
			{
				Title:      fmt.Sprintf("%v: %v", entry.Level.CapitalString(), entry.Message),
				Fields:     msgFields,
				MarkdownIn: []string{"fields"},
				Footer:     os.Args[0],
				Ts:         json.Number(strconv.FormatInt(entry.Time.Unix(), 10)),
				Color:      zapLevelColors[entry.Level],
			},
		},
	}
}

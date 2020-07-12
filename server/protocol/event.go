package protocol

import (
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/packet"
)

type EventMessage struct {
	BaseMessage

	Event string

	Data string
}

func ParseEventMessage(base BaseMessage, r *packet.PacketReader) (*EventMessage, error) {
	m := &EventMessage{}

	m.BaseMessage = base

	var ok bool

	m.Event, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	m.Data, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	return m, nil
}

func (m *EventMessage) Type() MessageType {
	return MessageEvent
}

func (m *EventMessage) Content() string {
	return m.Data
}

func (m *EventMessage) EventMsg() string {
	return m.Event
}

func (m *EventMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(MessageEvent))
	builder.WriteCString(m.BaseMessage.EntityName)

	builder.WriteCString(m.Event)
	builder.WriteCString(m.Data)

	return builder.Bytes()
}

func (m *EventMessage) Plain() string {
	return replacePlaceholders(m)
}

func (m *EventMessage) Embed() *discordgo.MessageEmbed {

	phrases := []string{"@everyone", "@here"}

	m.Data = cutPhrases(m.Data, phrases)
	m.Event = cutPhrases(m.Event, phrases)

	return &discordgo.MessageEmbed{
		Color:     16777215,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: m.BaseMessage.EntityName,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  m.Event,
				Value: m.Data,
			},
		},
	}
}

func (m *EventMessage) Webhook() *discordgo.WebhookParams {

	phrases := []string{"@everyone", "@here"}
	str := replacePlaceholders(m)

	var webhook discordgo.WebhookParams

	if m.Event == "ADMIN" {
		re := regexp.MustCompile(`ID64{(.+)} NAME{(.+)} MSG{(.+)}`)
		reg := re.FindStringSubmatch(m.Data)
		if len(reg) > 1 {
			webhook.AvatarURL = getAvatarURL("https://steamcommunity.com/profiles/" + reg[1] + "?xml=1")
			webhook.Username = reg[2]
			webhook.Content = cutPhrases(reg[3], phrases)
		}
	} else {
		webhook.Username = m.EntityName
		webhook.Content = cutPhrases(str, phrases)
	}

	return &webhook
}

func replacePlaceholders(m *EventMessage) string {
	var str string

	switch m.Event {
	case "Map Start":
		str = strings.ReplaceAll(config.Config.Messages.EventFormatSimpleMapStart, "%data%", m.Data)
	case "Map Ended":
		str = strings.ReplaceAll(config.Config.Messages.EventFormatSimpleMapEnd, "%data%", m.Data)
	case "Player Connected":
		str = strings.ReplaceAll(config.Config.Messages.EventFormatSimplePlayerConnect, "%data%", m.Data)
	case "Player Disconnected":
		str = strings.ReplaceAll(config.Config.Messages.EventFormatSimplePlayerDisconnect, "%data%", m.Data)
	default:
		str = strings.ReplaceAll(strings.ReplaceAll(config.Config.Messages.EventFormatSimple, "%data%", m.Data), "%event%", m.Event)
	}
	return str
}

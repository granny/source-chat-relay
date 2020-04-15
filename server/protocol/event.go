package protocol

import (
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

func (m *EventMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(MessageEvent))
	builder.WriteCString(m.BaseMessage.EntityName)

	builder.WriteCString(m.Event)
	builder.WriteCString(m.Data)

	return builder.Bytes()
}

func (m *EventMessage) Plain() string {

	switch m.Event {
	case "Map Start":
		return strings.ReplaceAll(config.Config.Messages.EventFormatSimpleMapStart, "%data%", m.Data)
	case "Map Ended":
		return strings.ReplaceAll(config.Config.Messages.EventFormatSimpleMapEnd, "%data%", m.Data)
	case "Player Connected":
		return strings.ReplaceAll(config.Config.Messages.EventFormatSimplePlayerConnect, "%data%", m.Data)
	case "Player Disconnected":
		return strings.ReplaceAll(config.Config.Messages.EventFormatSimplePlayerDisconnect, "%data%", m.Data)
	default:
		return strings.ReplaceAll(strings.ReplaceAll(config.Config.Messages.EventFormatSimple, "%data%", m.Data), "%event%", m.Event)
	}

}

func (m *EventMessage) Embed() *discordgo.MessageEmbed {

	phrases := [2]string{"@everyone", "@here"}

	for contains(phrases, m.Data) {
		m.Data = strings.ReplaceAll(m.Data, phrases[0], "")
		m.Data = strings.ReplaceAll(m.Data, phrases[1], "")
	}

	for contains(phrases, m.Event) {
		m.Event = strings.ReplaceAll(m.Event, phrases[0], "")
		m.Event = strings.ReplaceAll(m.Event, phrases[1], "")
	}

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

	phrases := [2]string{"@everyone", "@here"}

	for contains(phrases, str) {
		str = strings.ReplaceAll(str, phrases[0], "")
		str = strings.ReplaceAll(str, phrases[1], "")
	}

	return &discordgo.WebhookParams{
		Content: str,
	}
}

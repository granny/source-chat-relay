package bot

import (
	"strings"

	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/entity"
	"github.com/rumblefrog/source-chat-relay/server/relay"
	"github.com/sirupsen/logrus"
)

func Listen() {
	for {
		select {
		case message := <-relay.Instance.Bot:
			for _, guild := range RelayBot.State.Guilds {
				for _, channel := range guild.Channels {
					tEntity, err := entity.GetEntity(channel.ID)

					if err != nil {
						continue
					}

					if channel.ID != message.Author() &&
						tEntity.CanReceiveType(message.Type()) &&
						tEntity.ReceiveIntersectsWith(entity.DeliverableSendChannels(message)) {
						if !config.Config.Bot.SimpleMessage {
							RelayBot.ChannelMessageSendEmbed(channel.ID, message.Embed())
						} else if config.Config.Bot.Webhook {
							webhooks, err := RelayBot.ChannelWebhooks(channel.ID)
							if err != nil {
								logrus.Error(err.Error)
							} else {
								var id, token string
								var webhookName strings.Builder
								webhookName.WriteString("SCR ")
								webhookName.WriteString(channel.ID)
								for _, webhook := range webhooks {
									if strings.Contains(webhook.Name, webhookName.String()) {
										id = webhook.ID
										token = webhook.Token
									}
								}
								if id == "" {
									wh, err := RelayBot.WebhookCreate(channel.ID, webhookName.String(), "")
									if err == nil {
										id = wh.ID
										token = wh.Token
									}
								}
								RelayBot.WebhookExecute(id, token, false, message.Webhook())
							}
						} else {
							content := TransformMentions(RelayBot, channel.ID, message.Plain())
							RelayBot.ChannelMessageSend(channel.ID, content)
						}
					}
				}
			}
		}
	}
}

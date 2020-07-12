package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
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

						logrus.WithFields(logrus.Fields{
							"Author":  message.Webhook().Username,
							"Message": message.Content(),
							"Type":    message.Type(),
						}).Info("Message")

						if !config.Config.Bot.SimpleMessage {
							RelayBot.ChannelMessageSendEmbed(channel.ID, message.Embed())
						} else if config.Config.Bot.Webhook {
							webhooks, err := RelayBot.ChannelWebhooks(channel.ID)
							if err != nil {
								logrus.Error(err.Error)
							} else {
								var id, token string
								webh := message.Webhook()
								if message.EventMsg() == "ADMIN" {
									webhooks, err := RelayBot.ChannelWebhooks(config.Config.General.AdminChat)
									if err == nil {
										id, token = findWebhook(webhooks, config.Config.General.AdminChat)
									}
								} else if message.EventMsg() == "ADMIN2" {
									webhooks, err := RelayBot.ChannelWebhooks(config.Config.General.AdminChat2)
									if err == nil {
										id, token = findWebhook(webhooks, config.Config.General.AdminChat2)
									}
								} else {
									id, token = findWebhook(webhooks, channel.ID)
								}
								RelayBot.WebhookExecute(id, token, false, webh)
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

func findWebhook(webhooks []*discordgo.Webhook, channelID string) (id string, token string) {
	var lid, ltoken string
	for _, webhook := range webhooks {
		if strings.Contains(webhook.Name, channelID) {
			lid = webhook.ID
			ltoken = webhook.Token
		}
	}
	if lid == "" {
		wh, err := RelayBot.WebhookCreate(channelID, "SCR "+channelID, "")
		if err == nil {
			lid = wh.ID
			ltoken = wh.Token
		}
	}
	return lid, ltoken
}

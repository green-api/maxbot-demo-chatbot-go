package scenes

import (
	"time"

	m "github.com/green-api/maxbot-api-client-go/pkg/models"
	n "github.com/green-api/maxbot-chatbot-go/pkg/notification"
	"github.com/green-api/maxbot-chatbot-go/pkg/state"
	util "github.com/green-api/maxbot-demo-chatbot-go/utils"
)

type MainMenuScene struct{}

func (s *MainMenuScene) Start(app state.BotApp) {}

func (s *MainMenuScene) Execute(n *n.Notification) {
	if n.Type() == m.TypeMessageCallback {
		n.AnswerCallback("")
	}

	stateData := n.StateManager.GetStateData(n.StateId)
	lang, ok := stateData["lang"].(string)
	if !ok {
		lang = "ru"
	}

	s.SendMainMenu(n, lang)
}

func (s *MainMenuScene) SendMainMenu(n *n.Notification, lang string) {
	n.StateManager.UpdateStateData(n.StateId, map[string]any{"lang": lang})
	senderName, _ := n.SenderName()

	time.Sleep(500 * time.Millisecond)

	imageSource := "https://drive.google.com/uc?export=download&id=1gi2bPCQVgldRZolRH7gxevs1GIMnxmi2"
	imageAttachment := m.Attachment{
		Type: m.AttachmentImage,
		Payload: m.PhotoAttachmentPayload{
			URL: imageSource,
		},
	}

	btnLabels := map[string][]string{
		"ru": {"Текст 📩", "Файл 📋", "Картинка 🖼", "Аудио 🎵", "Видео 📽", "Контакт 📱", "Геолокация 🌎", "О боте 🦎", "Стоп"},
		"en": {"Text 📩", "File 📋", "Image 🖼", "Audio 🎵", "Video 📽", "Contact 📱", "Location 🌎", "About 🦎", "Stop"},
	}
	l := btnLabels[lang]

	keyboardAttachment := m.AttachKeyboard([][]m.KeyboardButton{
		{{Type: "callback", Text: l[0], Payload: "/message"}, {Type: "callback", Text: l[1], Payload: "/file"}},
		{{Type: "callback", Text: l[2], Payload: "/image"}, {Type: "callback", Text: l[3], Payload: "/audio"}},
		{{Type: "callback", Text: l[4], Payload: "/video"}, {Type: "callback", Text: l[5], Payload: "/contact"}},
		{{Type: "callback", Text: l[6], Payload: "/location"}, {Type: "callback", Text: l[7], Payload: "/about"}},
		{{Type: "callback", Text: l[8], Payload: "/stop"}},
	})
	n.ReplyWithAttachments(
		util.T(lang, "welcome_message")+"**"+senderName+"**!"+util.T(lang, "menu"),
		m.Markdown,
		[]m.Attachment{imageAttachment, keyboardAttachment},
	)

	n.ActivateNextScene(&EndpointsScene{})
}

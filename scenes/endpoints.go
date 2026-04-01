package scenes

import (
	"strings"
	"time"

	m "github.com/green-api/maxbot-api-client-go/pkg/models"
	n "github.com/green-api/maxbot-chatbot-go/pkg/notification"
	"github.com/green-api/maxbot-chatbot-go/pkg/state"
	util "github.com/green-api/maxbot-demo-chatbot-go/utils"
)

type EndpointsScene struct{}

func (s *EndpointsScene) Start(app state.BotApp) {}

func (s *EndpointsScene) Execute(n *n.Notification) {
	if n.Type() == m.TypeMessageCallback {
		n.AnswerCallback("")
	}

	text, err := n.Text()
	if err != nil {
		return
	}

	stateData := n.StateManager.GetStateData(n.StateId)
	lang, ok := stateData["lang"].(string)
	if !ok {
		lang = "ru"
	}

	senderId, _ := n.SenderID()
	senderName, _ := n.SenderName()

	switch strings.ToLower(text) {
	case "1", "/message":
		n.ShowAction("typing_on")
		n.ReplyWithKeyboard(
			util.T(lang, "send_text_message")+util.T(lang, "links.send_text_documentation"),
			m.Markdown,
			s.getControlButtons(lang),
		)
	case "2", "/file":
		n.ShowAction("sending_file")
		time.Sleep(200 * time.Millisecond)
		n.ReplyWithMedia(
			util.T(lang, "send_file_message")+util.T(lang, "links.send_file_documentation"),
			m.Markdown,
			"https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/corgi.pdf",
			s.getControlButtons(lang),
		)
		return
	case "3", "/image":
		n.ShowAction("sending_photo")
		n.ReplyWithMedia(
			util.T(lang, "send_image_message")+util.T(lang, "links.send_file_documentation"),
			m.Markdown,
			"https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/corgi.jpg",
			s.getControlButtons(lang),
		)
		return
	case "4", "/audio":
		n.ShowAction("sending_audio")
		time.Sleep(200 * time.Millisecond)

		audioUrl := "https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/Audio_bot.mp3"
		if lang == "ru" {
			audioUrl = "https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/Audio_bot_eng.mp3"
		}

		n.ReplyWithMedia(
			util.T(lang, "send_audio_message")+util.T(lang, "links.send_file_documentation"),
			m.Markdown,
			audioUrl,
			s.getControlButtons(lang),
		)
		return
	case "5", "/video":
		n.ShowAction("sending_video")
		time.Sleep(200 * time.Millisecond)

		videoUrl := "https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/Video_bot_eng.mp4"
		if lang == "ru" {
			videoUrl = "https://storage.yandexcloud.net/sw-prod-03-test/ChatBot/Video_bot_ru.mp4"
		}

		n.ReplyWithMedia(
			util.T(lang, "send_video_message")+util.T(lang, "links.send_file_documentation"),
			m.Markdown,
			videoUrl,
			s.getControlButtons(lang),
		)
		return
	case "6", "/contact":
		n.ReplyWithKeyboard(
			util.T(lang, "send_contact_message")+util.T(lang, "links.send_contact_documentation"),
			m.Markdown,
			s.getControlButtons(lang),
		)
		time.Sleep(500 * time.Millisecond)
		n.ReplyWithContact(senderName, "", &senderId)
		return
	case "7", "/location":
		n.ReplyWithKeyboard(
			util.T(lang, "send_location_message")+util.T(lang, "links.send_location_documentation"),
			m.Markdown,
			s.getControlButtons(lang),
		)
		time.Sleep(500 * time.Millisecond)
		n.ReplyWithLocation(35.888171, 14.440230)
		return
	case "8", "/about":
		aboutMsg := util.T(lang, "about_go_chatbot") +
			util.T(lang, "link_to_max") +
			util.T(lang, "link_to_docs") +
			util.T(lang, "link_to_source_code")

		n.ReplyWithKeyboard(
			aboutMsg,
			m.Markdown,
			s.getControlButtons(lang),
		)
		return
	case "стоп", "stop", "0", "/stop":
		n.Reply(
			util.T(lang, "stop_message")+senderName+"!",
			m.Markdown)
		n.ActivateNextScene(&StartScene{})
		return
	case "menu", "меню", "/menu":
		n.ActivateNextScene(&MainMenuScene{})
		sMenu := &MainMenuScene{}
		sMenu.SendMainMenu(n, lang)
		return
	default:
		n.Reply(util.T(lang, "not_recognized_message"), m.Markdown)
	}
}

func (s *EndpointsScene) getControlButtons(lang string) [][]m.KeyboardButton {
	mText, sText := "Меню", "Стоп"
	if lang == "en" {
		mText, sText = "Menu", "Stop"
	}

	return [][]m.KeyboardButton{
		{
			{Type: "callback", Text: mText, Payload: "/menu"},
			{Type: "callback", Text: sText, Payload: "/stop"},
		},
	}
}

package scenes

import (
	"github.com/green-api/maxbot-api-client-go/pkg/models"
	n "github.com/green-api/maxbot-chatbot-go/pkg/notification"
	"github.com/green-api/maxbot-chatbot-go/pkg/state"
)

type StartScene struct{}

func (s *StartScene) Start(app state.BotApp) {}

func (s *StartScene) Execute(n *n.Notification) {
	if n.Type() == models.TypeMessageCallback {
		n.AnswerCallback("")
	}

	text, err := n.Text()

	if err != nil || text == "/start" {
		s.askLanguage(n)
		return
	}

	switch text {
	case "1", "/english":
		s.proceedToMainMenu(n, "en")
	case "2", "/russian":
		s.proceedToMainMenu(n, "ru")
	default:
		s.askLanguage(n)
	}
}

func (s *StartScene) askLanguage(n *n.Notification) {
	buttons := [][]models.KeyboardButton{
		{{Type: "callback", Text: "English", Payload: "/english"}, {Type: "callback", Text: "Русский", Payload: "/russian"}},
	}
	keyboardAttachment := models.AttachKeyboard(buttons)
	n.ReplyWithAttachments(
		"Please select your language: \nПожалуйста, выберите язык: ",
		"",
		[]models.Attachment{keyboardAttachment},
	)
}

func (s *StartScene) proceedToMainMenu(n *n.Notification, lang string) {
	nextScene := &MainMenuScene{}
	n.ActivateNextScene(nextScene)
	nextScene.SendMainMenu(n, lang)
}

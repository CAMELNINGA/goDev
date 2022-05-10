package telegram

import (
	"Yaratam/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"

	"go.uber.org/zap"
)

type adapter struct {
	config bot.Config
	logger *zap.Logger
	bot    *tgbotapi.BotAPI
}

func NewAdapter(logger logrus.FieldLogger, config *Config) (domain.Sender, error) {
	a := &adapter{
		config: config,
		logger: logger,
	}

	var (
		err error
	)
	switch config.RunMode {
	case "direct":
		a.bot, err = bot.StartNoProxy(a.config)
	case "proxy":
		a.bot, err = bot.StartWithProxy(a.config)
	case "endpoint":
		a.bot, err = bot.StartWithCustomEndPoint(a.config)
	}
	if err != nil {
		a.logger.Error("Bot init failed", zap.Error(err))
		return nil, err
	}
	a.bot.Debug = a.config.Debug

	a.logger.Info("Authorized on", zap.String("username", a.bot.Self.UserName))

	return a, nil
}

func (a *adapter) SendTelegramTextMessage(telegramID int, text string) (msgID int, err error) {
	msg := tgbotapi.NewMessage(int64(telegramID), text)
	msg.ParseMode = tgbotapi.ModeHTML

	m, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Error while sending a message!", zap.Error(err))
		return 0, err
	}

	return m.MessageID, nil
}

func (a *adapter) ShareTelegramFile(telegramID int, fileID, fileCaption string) (msgID int, err error) {
	msg := tgbotapi.NewDocumentShare(int64(telegramID), fileID)
	msg.Caption = fileCaption

	m, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Error while sending a message!", zap.Error(err))
		return 0, err
	}

	return m.MessageID, nil
}
 gitlab.sovcombank.group/scb-mobile/underwriting/gamification-bot.git/internal/app/domain
func (a *adapter) SendTelegramKeyboardMessage(telegramID int, text string, keyboardType int) (msgID int, err error) {
	msg := tgbotapi.NewMessage(int64(telegramID), text)
	msg.ParseMode = tgbotapi.ModeHTML
	switch keyboardType {
	case domain.StartKeyboardType:
		msg.ReplyMarkup = domain.StartKeyboard
	case domain.DarkAndLightKeyboardType:
		msg.ReplyMarkup = domain.DarkAndLightKeyboard
	case domain.RaceKeyboardType:
		msg.ReplyMarkup = domain.RaceKeyboard
	case domain.ArenaKeyboardType:
		msg.ReplyMarkup = domain.ArenaKeyboard
	case domain.MonsterKeyboardType:
		msg.ReplyMarkup = domain.MonsterKeyboard
	case domain.PiratesStartKeyboardType:
		msg.ReplyMarkup = domain.StartKeyboardPirates
	default:
		msg.ReplyMarkup = domain.StartKeyboard
	}

	m, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Error while sending a message!", zap.Error(err))
		return 0, err
	}

	return m.MessageID, nil
}

func (a *adapter) SendImage(telegramID int, b []byte) (msgID int, fileID string, err error) {
	msg := tgbotapi.NewPhotoUpload(int64(telegramID), tgbotapi.FileBytes{
		Name:  "img",
		Bytes: b,
	})

	m, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Error while sending a image!", zap.Error(err))
		return 0, "", err
	}

	return m.MessageID, (*m.Photo)[len(*m.Photo)-1].FileID, nil
}

func (a *adapter) ShareImage(telegramID int, fileID string) (msgID int, err error) {
	msg := tgbotapi.NewPhotoShare(int64(telegramID), fileID)

	m, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Error while sending a image!", zap.Error(err))
		return 0, err
	}

	return m.MessageID, nil
}

func (a *adapter) DeleteMessage(telegramID int, messageID int) error {
	del := tgbotapi.NewDeleteMessage(int64(telegramID), messageID)

	_, err := a.bot.DeleteMessage(del)
	if err != nil {
		a.logger.Error("Error while removing a message", zap.Error(err))
		return err
	}
	return nil
}

func (a *adapter) ActivateNewKeyboard(telegramID int, typ int) error {
	msg := tgbotapi.NewMessage(int64(telegramID), domain.RegistrationSuccessful)

	switch typ {
	case 1:
		msg.ReplyMarkup = domain.StartKeyboard
	case 2:
		msg.ReplyMarkup = domain.StartKeyboardPirates
	default:
		msg.ReplyMarkup = domain.StartKeyboard
	}

	_, err := a.bot.Send(msg)
	if err != nil {
		a.logger.Error("Send new keyboard failed!", zap.Error(err))
		return err
	}
	return nil
}

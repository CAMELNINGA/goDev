package domain

import tgbotapi "github.com/Syfaro/telegram-bot-api"

// Bot Messages

const (
	StartMsg               = `Привет! Это бот похранению и структуризации документов телеграмме`
	RegisterMsg            = "Ваш код, для привязки Телеграм: <code>%s</code>"
	RegistrationSuccessful = "Бот-помощник активирован!"
	AlreadyRegistered      = "С возвращением тебя!"
	MuteModeActivated      = "Уведомления отключены до конца дня!"
	SupportText            = "Напишите нам: \n\n@came1l"
)

// Buttons
const (
	RegisterButton        = "🔑 Регистрация"
	SupportButton         = "💬 Служба поддержки"
	FortuneCookieButton   = "🎲 Предсказание"
	SaveFilesButton       = "Сохранить файл"
	CreateDirectoryButton = "Создать папку"
	ChooseDirectoryButton = "Ввойти в папку"
	DeleteFileButton      = "Удаление файла"
	ShareDirectory        = "Поделиться папкой"
)

const (
	StartKeyboardType = iota + 1
	DarkAndLightKeyboardType
	RaceKeyboardType
	ArenaKeyboardType
	MonsterKeyboardType
	PiratesStartKeyboardType
)

var StartKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(RegisterButton),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(FortuneCookieButton),
		tgbotapi.NewKeyboardButton(SupportButton),
	))

var MainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CreateDirectoryButton),
		tgbotapi.NewKeyboardButton(ChooseDirectoryButton),
	))

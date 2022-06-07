package bot

import (
	"Yaratam/internal/domain"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func StartNoProxy(c Config) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(c.Token)
}

const (
	ErrNoSideSelected   = "Вы не выбрали сторону! Выберите используя команду - /darkside или /lightside"
	ErrNoData           = "К сожалению, у нас пока нет данных за выбранный период. Информация может доставляться с задержкой"
	ErrHasFight         = "Ты уже участвуешь или запланировал другой бой. Прежде чем начать новый бой, заверши сначала предыдущий."
	ErrNoMonsters       = "Все монстры уже повержены."
	ErrNoHP             = "У тебя недостаточно жизней, чтобы начать бой."
	ErrNoDataF          = "Произошла ошибка, попробуйте позже."
	ErrNoAlivePlayers   = "Этот бой уже завершен."
	ErrNoButtonSelected = "Ошибка, для отправки файла необходимо сначала нажать кнопку"
)

var (
	ErrNoCommandSelected = errors.New("no command selected")
)

type Adapter interface {
	StartBot() error
	StopBot()
}

type adapter struct {
	config  Config
	service domain.Service
	logger  *zap.Logger
	bot     *tgbotapi.BotAPI

	commandsCache map[int]string
}

func NewAdapter(config Config, service domain.Service, logger *zap.Logger) (Adapter, error) {
	a := &adapter{
		config:        config,
		service:       service,
		logger:        logger.Named("TelegramBotAdapter"),
		commandsCache: make(map[int]string),
	}

	return a, nil
}

func (a *adapter) StartBot() error {
	var (
		err error
	)
	a.bot, err = StartNoProxy(a.config)

	if err != nil {
		a.logger.Error("Bot init failed", zap.Error(err))
		return err
	}
	a.bot.Debug = a.config.Debug

	a.logger.Info("Authorized on", zap.String("username", a.bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	upd, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range upd {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.Type != "private" {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				ply, err := a.service.GetUserData(int(update.Message.Chat.ID))
				if err != nil || ply.UserName == "" {
					a.logger.Error("Get player data failed", zap.Error(err))
					//continue
					msg.Text = domain.StartMsg
					msg.ReplyMarkup = domain.StartKeyboard
				} else {
					msg.ReplyMarkup = domain.MainKeyboard

					msg.Text = domain.AlreadyRegistered
				}
			case domain.RegisterButton:
				ply, err := a.service.GetUserData(int(update.Message.Chat.ID))
				if err != nil || ply.UserName == "" {
					user := domain.User{
						UserName: update.Message.Chat.UserName,
						ChatID:   int(update.Message.Chat.ID),
					}
					if err = a.service.AddUser(&user); err != nil {
						msg.Text = "Анлаки чет не получилось зарегаться"
						msg.ReplyMarkup = domain.StartKeyboard
					} else {
						msg.Text = domain.RegistrationSuccessful
						msg.ReplyMarkup = domain.MainKeyboard
					}

				} else {
					msg.Text = domain.AlreadyRegistered
					msg.ReplyMarkup = domain.MainKeyboard
				}
			case domain.ChooseDirectoryButton:
				{
					ply, err := a.service.GetUserData(int(update.Message.Chat.ID))
					if err != nil || ply.UserName == "" {
						a.logger.Error("Get player data failed", zap.Error(err))
						//continue
						msg.Text = domain.StartMsg
						msg.ReplyMarkup = domain.StartKeyboard
					} else {
						paths, err := a.service.GetPaths(int(update.Message.Chat.ID))
						if err != nil {
							a.logger.Error("Get player paths failed", zap.Error(err))
							msg.ReplyMarkup = domain.MainKeyboard
							msg.Text = "Не удалось выбрать папку"
						} else {
							msgPath := tgbotapi.NewReplyKeyboard()
							for _, v := range paths {
								msgPath.Keyboard = append(msgPath.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(v.DisplayName)))
							}
							msg.ReplyMarkup = msgPath
							msg.Text = "Выбери папку"
						}

					}
				}
			default:
				if update.Message.Document != nil {
					_, err = a.bot.GetFileDirectURL(update.Message.Document.FileID)
					if err != nil {
						if err.Error() == "Bad Request: file is too big" {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла: Размер файла больше 20MB ")
						}
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла!!!")
					}
					if err = a.service.AddFile(int(update.Message.Chat.ID), update.Message.Document.FileID); err != nil {
						if err == domain.ErrInvalidInputData {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла, вы не вошли в папаку")

						} else {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла")

						}
					} else {
						msg.Text = "Файл успешно отправлен!"
					}
					continue

				} else if update.Message.Command() != "" {
					paths, err := a.service.GetPaths(int(update.Message.Chat.ID))
					if err != nil {
						a.logger.Error("Get player paths failed", zap.Error(err))
						msg.ReplyMarkup = domain.MainKeyboard
						msg.Text = "Не удалось выбрать папку"
					} else {
						for i, v := range paths {
							if v.DisplayName == update.Message.Command() {
								if err = a.service.ChangeUserPath(int(update.Message.Chat.ID), v.ID); err != nil {
									msg.ReplyMarkup = domain.MainKeyboard
									msg.Text = "Не удалось выбрать папку"

								} else {
									files, err := a.service.GetFiles(int64(update.Message.Chat.ID))
									if err != nil {
										msg.Text = "Не удалось загрузить файлы"
									}
									for _, v := range files {
										tgbotapi.NewDocumentShare(int64(update.Message.Chat.ID), v.Path)
									}

									msg.ReplyMarkup = domain.MainKeyboard
									msg.Text = domain.AddedFile
									break
								}
							}
							if i == len(paths)-1 {
								msg.ReplyMarkup = domain.MainKeyboard
								msg.Text = "Не удалось найти папку папку"
							}
						}
					}
				} else {
					msg.Text = "Я не знаю такую команду, мне известны только [/start]"
				}

			}
			_, err = a.bot.Send(msg)
			if err != nil {
				a.logger.Error("Send message failed", zap.Error(err))
				continue
			}
		}
		a.logger.Debug("got message", zap.Any("chatid", update.Message.Chat.ID), zap.String("from", update.Message.From.UserName), zap.String("message", update.Message.Text))

	}
	return nil
}

func (a *adapter) StopBot() {
	a.bot.StopReceivingUpdates()
}

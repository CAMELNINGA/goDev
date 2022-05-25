package bot

import (
	"Yaratam/internal/domain"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
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

	upd := a.bot.GetUpdatesChan(u)

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
					if err := a.service.AddUser(&user); err != nil {
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
				if update.Message.Photo != nil || update.Message.Document != nil || update.Message.Video != nil {
					_, err := a.handleFile(&update)
					if err != nil {
						a.logger.Error("Handle file from message failed", zap.Error(err))
						continue
					}
					continue
					//msg.Text = "Файл успешно отправлен!"
				} else if update.Message.Command() != "" {
					paths, err := a.service.GetPaths(int(update.Message.Chat.ID))
					if err != nil {
						a.logger.Error("Get player paths failed", zap.Error(err))
						msg.ReplyMarkup = domain.MainKeyboard
						msg.Text = "Не удалось выбрать папку"
					} else {
						for i, v := range paths {
							if v.DisplayName == update.Message.Command() {
								if err := a.service.ChangeUserPath(int(update.Message.Chat.ID), v.ID); err != nil {
									msg.ReplyMarkup = domain.MainKeyboard
									msg.Text = "Не удалось выбрать папку"

								}

								tgbotapi.FileURL().SendData()
								msg.ReplyMarkup = domain.FileKeyboard
								msg.Text = domain.AddedFile
								break
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
			_, err := a.bot.Send(msg)
			if err != nil {
				a.logger.Error("Send message failed", zap.Error(err))
				continue
			}
		}
		a.logger.Debug("got message", zap.Any("chatid", update.Message.Chat.ID), zap.String("from", update.Message.From.UserName), zap.String("message", update.Message.Text))

	}
	return nil
}
func (a *adapter) handleFile(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	ply, err := a.service.GetUserData(int(update.Message.Chat.ID))
	if err != nil || ply.UserName == "" || ply.PathID == -1 {
		a.logger.Error("Get player data failed", zap.Error(err))
		//continue
	}

	rawURL := ""
	if update.Message.Video != nil {
		rawURL, err = a.bot.GetFileDirectURL(update.Message.Video.FileID)
		if err != nil {
			if err.Error() == "Bad Request: file is too big" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла: Размер файла больше 20MB можно загрузить только через сайт")
				sentMsg, errr := a.bot.Send(msg)
				if errr != nil {
					return nil, errr
				}

				return &sentMsg, fmt.Errorf("file is too big")
			}
			return nil, err
		}
	} else if update.Message.Document != nil {
		rawURL, err = a.bot.GetFileDirectURL(update.Message.Document.FileID)
		if err != nil {
			if err.Error() == "Bad Request: file is too big" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла: Размер файла больше 20MB можно загрузить только через сайт")
				sentMsg, errr := a.bot.Send(msg)
				if errr != nil {
					return nil, errr
				}

				return &sentMsg, fmt.Errorf("file is too big")
			}
			return nil, err
		}
	} else if update.Message.Photo != nil {
		rawURL, err = a.bot.GetFileDirectURL((*update.Message.Photo)[0].FileID)
		if err != nil {
			if err.Error() == "Bad Request: file is too big" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла: Размер файла больше 20MB можно загрузить только через сайт")
				sentMsg, errr := a.bot.Send(msg)
				if errr != nil {
					return nil, errr
				}

				return &sentMsg, fmt.Errorf("file is too big")
			}
			return nil, err
		}
	}
	if rawURL == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла, можно отправить только: фото, видео, файл")
		sentMsg, errr := a.bot.Send(msg)
		if errr != nil {
			return nil, errr
		}

		return &sentMsg, fmt.Errorf("file url is empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	fileName := path.Base(u.Path)

	if strings.HasSuffix(strings.ToLower(fileName), "png") {
	} else if strings.HasSuffix(strings.ToLower(fileName), "jpg") {
	} else if strings.HasSuffix(strings.ToLower(fileName), "jpeg") {
	} else if strings.HasSuffix(strings.ToLower(fileName), "mp4") {
	} else if strings.HasSuffix(strings.ToLower(fileName), "mov") {
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка:\nДля загрузки принимаются только следующие форматы: .jpg, .jpeg, .png, .mp4, .mov\nПожалуйста, загрузите другой файл.")
		sentMsg, errr := a.bot.Send(msg)
		if errr != nil {
			return nil, errr
		}

		return &sentMsg, fmt.Errorf("file extension is not allowed: %s", strings.ToLower(fileName))
	}

	resp, err := http.Get(u.String())
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения файла, попробуйте позже")
		sentMsg, errr := a.bot.Send(msg)
		if errr != nil {
			return nil, errr
		}

		return &sentMsg, err
	}

	filLink, err := a.service.UploadMultipartFile(resp.Body, ply.UserName, strconv.Itoa(ply.ChatID), fileName)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла, попробуйте позже")
		sentMsg, errr := a.bot.Send(msg)
		if errr != nil {
			return nil, errr
		}

		return &sentMsg, err
	}
	if err := a.service.AddFile(int(update.Message.Chat.ID), filLink); err != nil {
		if err == domain.ErrInvalidInputData {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла, вы не вошли в папаку")
			sentMsg, errr := a.bot.Send(msg)
			if errr != nil {
				return nil, errr
			}
			return &sentMsg, err
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла")
			sentMsg, errr := a.bot.Send(msg)
			if errr != nil {
				return nil, errr
			}
			return &sentMsg, err
		}

	}
	defer resp.Body.Close()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Файл сохранен."))
	sentMsg, err := a.bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &sentMsg, nil
}

func (a *adapter) StopBot() {
	a.bot.StopReceivingUpdates()
}

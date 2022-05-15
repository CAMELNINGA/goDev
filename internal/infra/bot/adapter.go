package bot

import (
	"Yaratam/internal/domain"
	"errors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
	"regexp"
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

	upd, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		a.logger.Error("Bot get updates channel failed", zap.Error(err))
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
					msg.ReplyMarkup = domain.StartKeyboard

					msg.Text = domain.AlreadyRegistered
				}
			default:
				res := strings.Contains(update.Message.Command(), "look")
				if res {
					re := regexp.MustCompile("[0-9]+")
					fightIDStr := re.FindAllString(update.Message.Command(), 1)
					ply, err := a.service.GetUserData(int(update.Message.Chat.ID))
					if err != nil {
						a.logger.Error("Get player data failed", zap.Error(err))
						if err != nil {
							continue
						}
						continue
					}

					msg.Text = domain.JoinedGroupFight
				} else {
					if update.Message.Photo != nil || update.Message.Document != nil || update.Message.Video != nil {
						_, err := a.handleFile(&update)
						if err != nil {
							a.logger.Error("Handle file from message failed", zap.Error(err))
							continue
						}
						continue
						//msg.Text = "Файл успешно отправлен!"
					} else {
						msg.Text = "Я не знаю такую команду, мне известны только [/start]"
					}
				}
			}
		}

	}
}
func (a *adapter) handleFile(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	mstext := ErrNoDataF
	typ := 0
	ply, err := a.service.GetPlayerData(int(update.Message.Chat.ID))
	if err != nil || ply.Username == "" {
		a.logger.Error("Get player data failed", zap.Error(err))
		//continue
	} else {
		gm, err := a.service.GetGameName(ply.UnitID)
		if err != nil {
			a.logger.Error("Error while GetGameName!", zap.Error(err))
		}
		switch gm {
		case "pirates":
			mstext = ErrNoButtonSelected
			typ = 1
		case "star_wars":
			mstext = ErrNoSideSelected
			typ = 2
		}
	}

	vv, ok := a.commandsCache[int(update.Message.Chat.ID)]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, mstext)
		sentMsg, err := a.bot.Send(msg)
		if err != nil {
			return nil, err
		}

		return &sentMsg, fmt.Errorf("no button selected")
	}
	// reset button after usage
	defer delete(a.commandsCache, int(update.Message.Chat.ID))

	if typ == 2 {
		var tmpRussSide string
		switch vv {
		case "darkside":
			tmpRussSide = "Темная сторона"
		case "lightside":
			tmpRussSide = "Светлая сторона"
		}
		msg := tgbotapi.NewMessage(int64(266789284), fmt.Sprintf("Пользователь прислал фото\n\nВыбранная сторона: %s \nОтправитель -\nTelegramID: %v\nTelegram Username: %s\nTelegram Name: %s", tmpRussSide, update.Message.Chat.ID, update.Message.Chat.UserName, update.Message.Chat.FirstName+" "+update.Message.Chat.LastName))
		_, err = a.bot.Send(msg)
		if err != nil {
			a.logger.Error("Send message failed", zap.Error(err))
			return nil, err
		}

		fwd := tgbotapi.NewForward(int64(266789284), update.Message.Chat.ID, update.Message.MessageID)
		_, err = a.bot.Send(fwd)
		if err != nil {
			a.logger.Error("Send message failed", zap.Error(err))
			return nil, err
		}
		return nil, nil
	} else if typ == 1 {

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
		//defer resp.Body.Close()

		filLink, err := a.service.UploadMultipartFile(resp.Body, ply.Username, strconv.Itoa(ply.UnitID), fileName)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка загрузки файла, попробуйте позже")
			sentMsg, errr := a.bot.Send(msg)
			if errr != nil {
				return nil, errr
			}

			return &sentMsg, err
		}
		_, err = a.service.PostPunishmentFile(ply.UnitID, ply.Username, filLink, punishments[playerPunishmentID].PunishmentID, punishments[playerPunishmentID].ActionPeriod)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка публикации файла, попробуйте позже")
			sentMsg, errr := a.bot.Send(msg)
			if errr != nil {
				return nil, errr
			}

			return &sentMsg, err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Отлично! Файл будет проверен модератором, а пока держи свою награду: %d рубинов.", punishments[playerPunishmentID].PlayerPrize))
		sentMsg, err := a.bot.Send(msg)
		if err != nil {
			return nil, err
		}

		return &sentMsg, nil
	}
	return nil, nil
}

func (a *adapter) StopBot() {
	a.bot.StopReceivingUpdates()
}

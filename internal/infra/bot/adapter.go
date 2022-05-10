package bot

import (
	"Yaratam/internal/domain"
	"errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
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

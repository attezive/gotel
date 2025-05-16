package bot

import (
	"github.com/attezive/gotel/data"
	"github.com/attezive/gotel/internal/deleter"
	"github.com/attezive/gotel/internal/handler"
	"github.com/attezive/gotel/internal/menu"
	"github.com/attezive/gotel/internal/sender"
	"strconv"
)

type GotelBot struct {
	token   string
	handler *handler.Handler
	Sender  *sender.Sender
	deleter *deleter.Deleter
	menu    *menu.Menu
	stop    chan bool
}

func CreateBot(token ...string) *GotelBot {
	bot := new(GotelBot)
	if len(token) != 0 {
		tokenParts := ""
		for _, tokenPart := range token {
			tokenParts += tokenPart
		}
		bot.token = tokenParts
	}
	bot.stop = make(chan bool)
	bot.handler = handler.CreateHandler(&bot.token)
	bot.Sender = sender.CreateSender(&bot.token)
	bot.deleter = deleter.CreateDeleter(&bot.token)
	bot.menu = menu.CreateMenu(&bot.token)
	return bot
}

func (tBot *GotelBot) GetToken() string {
	return tBot.token
}

func (tBot *GotelBot) SetToken(token string) {
	tBot.token = token
}

func (tBot *GotelBot) AddHandleFunction(handleFunction func(*data.Update)) {
	tBot.handler.AddHandleFunction(handleFunction)
}

func (tBot *GotelBot) Start() <-chan error {
	tBot.handler.Start = true
	tBot.stop = make(chan bool)
	errChan := make(chan error, 1)
	go tBot.handler.Handle(tBot.stop, errChan)
	return errChan
}

func (tBot *GotelBot) Stop() {
	tBot.stop <- true
}

func (tBot *GotelBot) SendMessage(message *data.SendingEntity) (<-chan *data.Message, <-chan error) {
	returnedMsg := make(chan *data.Message, 1)
	errChan := make(chan error, 1)
	go tBot.Sender.SendMessage(message, returnedMsg, errChan)
	return returnedMsg, errChan
}

func (tBot *GotelBot) SendPhoto(message *data.SendingEntity, saveFileId bool) (<-chan *data.Message, <-chan error) {
	returnedMsg := make(chan *data.Message, 1)
	errChan := make(chan error, 1)
	go tBot.Sender.SendPhoto(message, returnedMsg, errChan)
	if saveFileId {
		if err := <-errChan; err != nil {
			errChan <- err
			return returnedMsg, errChan
		}
		errChan <- nil
		photo := message.Value.(*data.Photo)
		msg := <-returnedMsg
		photo.FileId = (msg).Photo[0].FileId
		returnedMsg <- msg
	}
	return returnedMsg, errChan
}

func (tBot *GotelBot) AddReaction(
	handleFunction func(*data.Update) interface{},
	sendFunction func(*data.SendingEntity, chan<- *data.Message, chan<- error)) (<-chan *data.Message, <-chan error) {
	errCh := make(chan error, 1)
	msgCh := make(chan *data.Message, 1)
	tBot.AddHandleFunction(func(update *data.Update) {
		value := handleFunction(update)
		entity := data.SendingEntity{
			ChatId: strconv.FormatInt(update.Message.Chat.Id, 10),
			Value:  value}
		sendFunction(&entity, msgCh, errCh)
	})
	return msgCh, errCh
}

func (tBot *GotelBot) DeleteMessage(chatId string, messageId string) (*data.SuccessResponse, error) {
	deleteResponse, err := tBot.deleter.DeleteMessage(chatId, messageId)
	return deleteResponse, err
}

func (tBot *GotelBot) GetCommands() (*[]data.BotCommand, error) {
	commands, err := tBot.menu.GetMyCommands()
	return commands, err
}

func (tBot *GotelBot) SetCommands(newCommands *[]data.BotCommand, saveOldCommands bool) (*data.SuccessResponse, error) {
	commandResponse, err := tBot.menu.SetMyCommands(newCommands, saveOldCommands)
	return commandResponse, err
}

package bot

import (
	"gotel_alpha/data"
	"gotel_alpha/internal/deleter"
	"gotel_alpha/internal/handler"
	"gotel_alpha/internal/menu"
	"gotel_alpha/internal/sender"
)

type GotelBot struct {
	token   string
	handler *handler.Handler
	sender  *sender.Sender
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
	bot.sender = sender.CreateSender(&bot.token)
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
	errCh := make(chan error, 1)
	go tBot.handler.Handle(tBot.stop, errCh)
	return errCh
}

func (tBot *GotelBot) Stop() {
	tBot.stop <- true
}

func (tBot *GotelBot) SendMessage(message *data.SendingEntity) (*data.Message, error) {
	returnedMsg, err := tBot.sender.SendMessage(message)
	return returnedMsg, err
}

func (tBot *GotelBot) SendPhoto(message *data.SendingEntity, saveFileId bool) (*data.Message, error) {
	returnedMsg, err := tBot.sender.SendPhoto(message)
	if saveFileId {
		if err != nil {
			return nil, err
		}
		photo := message.Value.(*data.Photo)
		photo.FileId = returnedMsg.Photo[0].FileId
	}
	return returnedMsg, err
}

// AddReaction is unsafe operation with panic when error in send request
func (tBot *GotelBot) AddReaction(
	handleFunction func(*data.Update) interface{},
	sendFunction func(*data.SendingEntity) (*data.Message, error),
	handleMessageFunction func(*data.Message)) {
	tBot.AddHandleFunction(func(update *data.Update) {
		value := handleFunction(update)
		err := tBot.sender.ReactionSend(update, value, sendFunction, handleMessageFunction)
		if err != nil {
			panic(err)
		}
	})
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

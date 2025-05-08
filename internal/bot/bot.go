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

func (tBot *GotelBot) AddHandleFunction(handleFunction func(*handler.Update)) {
	tBot.handler.AddHandleFunction(handleFunction)
}

func (tBot *GotelBot) Start() error {
	tBot.handler.Start = true
	return tBot.handler.Handle(tBot.stop)
}

func (tBot *GotelBot) Stop() error {
	tBot.StopHandling()
	return nil
}

func (tBot *GotelBot) StopHandling() {
	tBot.stop <- true
}

func (tBot *GotelBot) SendMessage(message *sender.SendingEntity) (*data.Message, error) {
	returnedMsg, err := tBot.sender.SendMessage(message)
	return returnedMsg, err
}

func (tBot *GotelBot) SendPhoto(message *sender.SendingEntity, saveFileId bool) (*data.Message, error) {
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
	handleFunction func(*handler.Update) interface{},
	sendFunction func(*sender.SendingEntity) (*data.Message, error),
	handleMessageFunction func(*data.Message)) {
	tBot.AddHandleFunction(func(update *handler.Update) {
		value := handleFunction(update)
		err := tBot.sender.ReactionSend(update, value, sendFunction, handleMessageFunction)
		if err != nil {
			panic(err)
		}
	})
}

func (tBot *GotelBot) DeleteMessage(chatId string, messageId string) (*deleter.DeleteResponse, error) {
	deleteResponse, err := tBot.deleter.DeleteMessage(chatId, messageId)
	return deleteResponse, err
}

func (tBot *GotelBot) GetCommands() (*[]menu.BotCommand, error) {
	commands, err := tBot.menu.GetMyCommands()
	return commands, err
}

func (tBot *GotelBot) SetCommands(newCommands *[]menu.BotCommand, saveOldCommands bool) (*menu.CommandResponse, error) {
	commandResponse, err := tBot.menu.SetMyCommands(newCommands, saveOldCommands)
	return commandResponse, err
}

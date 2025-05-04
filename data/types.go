package data

type MessageInfo struct {
	Ok      bool    `json:"ok"`
	Message Message `json:"result"`
}

type Message struct {
	MessageId int         `json:"message_id"`
	From      User        `json:"from,omitempty"`
	Chat      Chat        `json:"chat"`
	Date      int64       `json:"date"`
	EditDate  int64       `json:"edit_date,omitempty"`
	Text      string      `json:"text,omitempty"`
	Entities  interface{} `json:"entities,omitempty"`
	EffectId  string      `json:"effect_id,omitempty"`
	Animation interface{} `json:"animation,omitempty"`
	Audio     interface{} `json:"audio,omitempty"`
	Document  interface{} `json:"document,omitempty"`
	Photo     []Photo     `json:"photo,omitempty"`
	Sticker   interface{} `json:"sticker,omitempty"`
	Video     interface{} `json:"video,omitempty"`
	Story     interface{} `json:"story,omitempty"`
	VideoNote interface{} `json:"video_note,omitempty"`
	Voice     interface{} `json:"voice,omitempty"`
	Caption   interface{} `json:"caption,omitempty"`
	Contact   interface{} `json:"contact,omitempty"`
	Location  interface{} `json:"location,omitempty"`
}

type User struct {
	Id           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type Photo struct {
	FilePath string `json:"-"`
	FileId   string `json:"file_id"`
}

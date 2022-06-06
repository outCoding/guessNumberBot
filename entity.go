package main

type AuthRequest struct {
	Op int     `json:"op"`
	D  Intents `json:"d"`
}

type Intents struct {
	Token   string `json:"token"`
	Intents int    `json:"intents"`
	Shard   []int  `json:"shard"`
}

type HeatBeat struct {
	Op int `json:"op"`
	D  int `json:"d"`
}

type Message struct {
	Content string `json:"content"`
	MsgId   string `json:"msg_id"`
}

type Response struct {
	Op int    `json:"op"`
	S  int    `json:"s"`
	T  string `json:"t"`
	D  Data   `json:"d"`
}

type Data struct {
	Author struct {
		Avatar   string `json:"avatar"`
		Id       string `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
	ChannelId     string `json:"channel_id"`
	Content       string `json:"content"`
	DirectMessage bool   `json:"direct_message"`
	GuildId       string `json:"guild_id"`
	Id            string `json:"id"`
	Member        struct {
		JoinedAt string `json:"joined_at"`
	} `json:"member"`
	Seq          int    `json:"seq"`
	SeqInChannel string `json:"seq_in_channel"`
	SrcGuildId   string `json:"src_guild_id"`
	Timestamp    string `json:"timestamp"`
}

type GameData struct {
	Number int `json:"number"`
	Times  int `json:"times"`
}

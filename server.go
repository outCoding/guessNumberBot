package main

import (
	"encoding/json"
	"log"
	"strings"
)

func main() {
	for true {
		log.Printf("start server...")
		startServer()
	}
}

func startServer() {
	// 建立连接
	client := getWebSocketClient(websocketUrl)

	// 鉴权
	authentication(client)

	// 维持心跳
	go heatBeat(client)

	for true {
		// 接收消息
		_, msg, err := client.ReadMessage()
		if err != nil {
			log.Printf("receive err msg:[" + err.Error() + "]")
			break
		}
		resp := new(Response)
		err = json.Unmarshal(msg, resp)
		if resp.S > 0 {
			heatBeatD = resp.S
		}
		dealMessage(resp)
		log.Printf("receive msg:[" + string(msg) + "]")
	}
	defer client.Close()
}

// 处理消息
func dealMessage(resp *Response) {
	if resp.T == "AT_MESSAGE_CREATE" {
		content := strings.Replace(resp.D.Content, "<@!"+botId+">", "", 1)
		content = strings.Replace(content, " ", "", -1)
		if content == "开始猜数" {
			startGame(resp.D.ChannelId, resp.D.Id, resp.D.Author.Id, content)
		} else if isNumber(content) {
			doGame(resp.D.ChannelId, resp.D.Id, resp.D.Author.Id, content)
		} else if content == "结束猜数" {
			endGame(resp.D.ChannelId, resp.D.Id, resp.D.Author.Id, content)
		} else {
			message := Message{Content: "<@!" + resp.D.Author.Id + "> 请发送合法的指令", MsgId: resp.D.Id}
			messageData, _ := json.Marshal(message)
			sendMessage(resp.D.ChannelId, messageData)
		}
	}
}

package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var apiBaseUrl = "https://api.sgroup.qq.com"
var websocketUrl = "wss://api.sgroup.qq.com/websocket"

var token = "Bot appId.token"
var botId = "botId"

var heatBeatD int

var ticker *time.Ticker = nil

var gameDataMap = make(map[string]GameData)

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
			log.Printf("receive err msg:", err)
			break
		}
		resp := new(Response)
		err = json.Unmarshal(msg, resp)
		if resp.S > 0 {
			heatBeatD = resp.S
		}
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
		log.Printf("receive msg:[" + string(msg) + "]")
	}
	defer client.Close()
}

func startGame(channelId string, msgId string, userId string, content string) {
	_, exist := gameDataMap[channelId+":"+userId]
	if exist {
		message := Message{Content: "<@!" + userId + "> 已经存在进行中的猜数游戏了，若要重新开始请先发送“结束猜数”", MsgId: msgId}
		messageData, _ := json.Marshal(message)
		sendMessage(channelId, messageData)
		return
	}
	num := rand.Intn(100) + 1
	gameDataMap[channelId+":"+userId] = GameData{Number: num, Times: 1}
	message := Message{Content: "<@!" + userId + "> 猜数游戏开始, 开始猜第1次（数字范围1-100）", MsgId: msgId}
	messageData, _ := json.Marshal(message)
	sendMessage(channelId, messageData)
}

func doGame(channelId string, msgId string, userId string, content string) {
	gameData, exist := gameDataMap[channelId+":"+userId]
	var message = Message{}

	number, err := strconv.Atoi(content)
	if err != nil {
		return
	}

	if !exist {
		message = Message{Content: "<@!" + userId + "> 请先发送“开始猜数”", MsgId: msgId}
	} else if number == gameData.Number {
		message = Message{Content: "<@!" + userId + "> 你猜对了，正确的数字是:" + strconv.Itoa(number), MsgId: msgId}
		delete(gameDataMap, channelId+":"+userId)
	} else if number < gameData.Number {
		message = Message{Content: "<@!" + userId + "> 你第" + strconv.Itoa(gameData.Times) + "次猜错了，比这个数字大哦", MsgId: msgId}
		gameData.Times++
		gameDataMap[channelId+":"+userId] = gameData
	} else if number > gameData.Number {
		message = Message{Content: "<@!" + userId + "> 你第" + strconv.Itoa(gameData.Times) + "次猜错了，比这个数字小哦", MsgId: msgId}
		gameData.Times++
		gameDataMap[channelId+":"+userId] = gameData
	} else {
		return
	}
	messageData, _ := json.Marshal(message)
	sendMessage(channelId, messageData)
}

func endGame(channelId string, msgId string, userId string, content string) {
	delete(gameDataMap, channelId+":"+userId)
	message := Message{Content: "<@!" + userId + "> 游戏结束", MsgId: msgId}
	messageData, _ := json.Marshal(message)
	sendMessage(channelId, messageData)
}

func isNumber(content string) bool {
	for _, r := range content {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func heatBeat(client *websocket.Conn) {
	if ticker != nil {
		ticker.Stop()
		heatBeatD = 0
		log.Printf("heatBeat end")
	}
	ticker = time.NewTicker(41250 * time.Millisecond)
	for range ticker.C {
		heatBeat := HeatBeat{Op: 1, D: heatBeatD}
		heatBeatData, _ := json.Marshal(heatBeat)
		log.Printf("send headBeat msg:[%s] [%s]", strconv.Itoa(heatBeat.Op), strconv.Itoa(heatBeat.D))
		client.WriteMessage(websocket.TextMessage, heatBeatData)
	}
	return
}

func authentication(client *websocket.Conn) {
	var authRequest = AuthRequest{
		Op: 2,
		D: Intents{
			Token:   token,
			Intents: 1073741824}}
	authRequestData, _ := json.Marshal(authRequest)
	err := client.WriteMessage(websocket.TextMessage, authRequestData)
	if err != nil {
		log.Println(err)
	}
}

func getWebSocketClient(url string) *websocket.Conn {
	client, _, _ := websocket.DefaultDialer.Dial(url, nil)
	return client
}

func sendMessage(channelId string, data []byte) {
	sendRequest("POST", "/channels/"+channelId+"/messages", "application/json", data)
}

func sendRequest(method string, url string, contentType string, body []byte) (string, error) {
	url = apiBaseUrl + url
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

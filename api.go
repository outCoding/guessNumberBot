package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var apiBaseUrl = "https://api.sgroup.qq.com"
var websocketUrl = "wss://api.sgroup.qq.com/websocket"

var token = "Bot appId.token"
var botId = "botId"

var heatBeatD int

var ticker *time.Ticker = nil

// 获取websocket连接
func getWebSocketClient(url string) *websocket.Conn {
	client, _, _ := websocket.DefaultDialer.Dial(url, nil)
	return client
}

// 心跳包发送
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

// 鉴权请求
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

// 发送消息
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

package main

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

var gameDataMap = make(map[string]GameData)

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

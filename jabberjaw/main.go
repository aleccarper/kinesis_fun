package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type (
	Data struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}
)

type (
	Payload struct {
		Data []Data `json:"data"`
	}
)

func main() {
	url := "http://localhost:3000/produce"

	for true {
		jsonValue, _ := json.Marshal(buildPayload())

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		fmt.Println(resp)
		fmt.Println(err)
	}

}

func buildPayload() Payload {
	var payload Payload

	payload.Data = append(payload.Data, buildData())
	payload.Data = append(payload.Data, buildData())

	return payload
}

func buildData() Data {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}

	var data Data
	data.Key1 = answers[rand.Intn(len(answers))]
	data.Key2 = answers[rand.Intn(len(answers))]
	return data
}

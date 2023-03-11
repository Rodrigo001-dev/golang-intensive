package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Role string  `json:"role"`
	Content string  `json:"content"`
}

type Request struct {
	Model string `json:"model"`
	Messages []Message `json:"messages"`
	// um token é o resultado de dados que é retornado com resposta da pergunta
	// o ChatGPT não utilizar necessáriamente palavras para computar o token, ele
	// utiliza palavras inteiras ou pedaços de palavras, cada pedaço dessas palavras
	// como token entra no modelo deles para que possa ser possível gerar as repostas
	// e nós vamos ser cobrados pela quantidade de tokens emitido 
	MaxTokens int `json:"max_tokens,omitempty"` // omitempty significa que pode ser opcial
}

type Response struct {
	ID string `json:"id"`
	Object string `json:"object"`
	Created int `json:"created"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index int `json:"index"`
	Message struct {
		Role string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func GenerateGPTText(query string) (string, error) {
	req := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role: "use",
				Content: query,
			},
		},
		MaxTokens: 150,
	}

	reqJson, err := json.Marshal(req) // transformando em JSON

	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJson))

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer YOUR_API_KEY")

	response, err := http.DefaultClient.Do(request)
	
	if err != nil {
		return "", err
	}

	defer response.Body.Close() // roda tudo que ta em baixo do defer e depois roda o defer

	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	// pegar o json e converter para struct
	var resp Response
	err = json.Unmarshal(resBody, &resp)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
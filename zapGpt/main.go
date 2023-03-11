package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func parseBase64RequestData(r string) (string, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return "", err
	}

	data, err := url.ParseQuery(string(dataBytes))
	if data.Has("Body") {
		return data.Get("Body"), nil
	}

	return "", errors.New("Body not found")
}

func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, err := parseBase64RequestData(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: err.Error(),
		}, nil
	}

	text, err := GenerateGPTText(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body: err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: text,
	}, nil
}

func main() {
	lambda.Start(process)
}
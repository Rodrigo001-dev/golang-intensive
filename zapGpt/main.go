package main

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
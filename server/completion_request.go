package main

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"gopkg.in/validator.v2"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
)

var RecipePrompt = `Write a delicious named recipe with the following ingredients:
%s.

Title and recipe:`

func CompletionRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing request...")
	var resp gogpt.CompletionResponse
	error_code := "none"
	type Response struct {
		ErrorCode string `json:"errorcode"`
		Completion string `json:"completion"`
	}
	defer func() {
		var response Response
		if len(resp.Choices) > 0 {
			response = Response{error_code, strings.TrimSpace(resp.Choices[0].Text)}
		} else {
			fmt.Println("Request Failed!:", error_code)
			response = Response{error_code, ""}
		}
		response_json, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		w.Write(response_json)
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error_code = err.Error()
		return
	}
	query_req := struct {
		Ingredients string `validate:"nonzero"`
	}{}
	json.Unmarshal(body, &query_req)
	if errs := validator.Validate(query_req); errs != nil {
		error_code = "Invalid request: " + errs.Error()
		return
	}

	ctx := context.Background()

	prompt := fmt.Sprintf(RecipePrompt, query_req.Ingredients)
	req := gogpt.CompletionRequest{
		Model: "text-davinci-002",
		MaxTokens: 256,
		Temperature: 0.6,
		Prompt: prompt,
	}
	resp, err = GptClient.CreateCompletion(ctx, req)
	if err != nil {
		error_code = "Failed to create completion: " + err.Error()
		return
	}
	fmt.Println("Request Completed Sucessfuly!")
}

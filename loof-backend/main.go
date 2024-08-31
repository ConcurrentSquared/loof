package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/viper"
	"github.com/tmaxmax/go-sse"
)

type OpenRouterRequest struct {
	Model string `json:"model"`

	Prompt string `json:"prompt"`

	Transforms []string `json:"transforms"`

	MaxTokens int  `json:"max_tokens"`
	Streaming bool `json:"stream"`
}

type OpenRouterResponse struct {
	Choices []OpenRouterNonChatData `json:"choices"`
}

// Non-chat streaming data uses the same payloads as regular non-chat data
type OpenRouterNonChatData struct {
	Text string `json:"text"`
}

func ExpandText(array *[]string, currentRecord *models.Record, database *daos.Dao) {
	*array = append(*array, currentRecord.GetString("text"))

	previousNodeId := currentRecord.GetString("previous_node")
	if previousNodeId != "" {
		previous_node, err := database.FindRecordById("nodes", currentRecord.GetString("previous_node"))
		if err != nil {
			panic(err)
		}

		ExpandText(array, previous_node, database)
	}
}

func main() {
	app := pocketbase.New()
	streamingClient := sse.Client{Backoff: sse.Backoff{MaxRetries: -1}}

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("LOOF_SETTINGS_")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	apiKey := viper.GetString("OPENROUTER_API_KEY")

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	app.OnRecordAfterUpdateRequest("nodes").Add(func(e *core.RecordUpdateEvent) error {
		println(e.Record.GetString("state"))
		if e.Record.GetString("state") == "pending" {
			e.Record.Set("state", "editing")

			nodeTextList := []string{}
			ExpandText(&nodeTextList, e.Record, app.Dao())

			prompt := ""
			for index := len(nodeTextList) - 1; index >= 0; index-- {
				prompt += nodeTextList[index]
			}

			databaseRecord, err := app.Dao().FindRecordById("nodes", e.Record.Id)
			if err != nil {
				panic(err)
			}

			author, err := app.Dao().FindRecordById("authors", e.Record.GetString("author"))
			if err != nil {
				panic(err)
			}

			if author.GetString("origin") == "robot" {
				openRouterRequest := OpenRouterRequest{Model: "meta-llama/llama-3.1-405b", Prompt: prompt, Transforms: []string{}, MaxTokens: 128, Streaming: true}

				jsonData, err := json.Marshal(openRouterRequest)
				println(string(jsonData))
				if err != nil {
					return err
				}

				req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
				if err != nil {
					return err
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", ("Bearer " + apiKey))

				req.Header.Set("HTTP-Referer", "https://www.concurrentsquared.com/projects/loof")
				req.Header.Set("X-Title", "Loof")

				responseTokens := ""
				conn := streamingClient.NewConnection(req)
				conn.SubscribeMessages(func(sse sse.Event) {
					println(sse.Data)
					println(sse.Type)
					if sse.Data == "[DONE]" {
						databaseRecord.Set("state", "completed")
						app.Dao().SaveRecord(databaseRecord)
					} else if sse.Data != "" { // OpenRouter stream completion message
						var openRouterResponse OpenRouterResponse
						err = json.Unmarshal([]byte(sse.Data), &openRouterResponse)
						if err != nil {
							panic(err)
						}

						responseTokens += openRouterResponse.Choices[0].Text
						databaseRecord.Set("text", responseTokens)
						app.Dao().SaveRecord(databaseRecord)
					}
				})

				if err = conn.Connect(); !errors.Is(err, io.EOF) {
					return err
				}
			}
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

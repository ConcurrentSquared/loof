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
	"github.com/spf13/viper"
	"github.com/tmaxmax/go-sse"
)

type OpenRouterResponse struct {
	Choices []OpenRouterNonChatData `json:"choices"`
}

// Non-chat streaming data uses the same payloads as regular non-chat data
type OpenRouterNonChatData struct {
	Text string `json:"text"`
}

func main() {
	app := pocketbase.New()
	streamingClient := sse.Client{Backoff: sse.Backoff{MaxRetries: -1}}

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	//viper.SetEnvPrefix("LOOF_SETTINGS_")
	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	apiKey := viper.GetString("OPENROUTER_API_KEY")

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	app.OnRecordAfterCreateRequest("nodes").Add(func(e *core.RecordCreateEvent) error {
		databaseRecord, err := app.Dao().FindRecordById("nodes", e.Record.Id)
		if err != nil {
			panic(err)
		}

		author, err := app.Dao().FindRecordById("authors", e.Record.GetString("author"))

		if true {
			if err != nil {
				panic(err)
			}

			if author.GetString("origin") == "robot" {
				jsonData := []byte(`{"model": "meta-llama/llama-3.1-8b-instruct:free", "prompt": "What is the meaning of bingo?!?!?", "max_tokens": 128, "stream": true}`)
				req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
				if err != nil {
					panic(err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", ("Bearer " + apiKey))

				responseTokens := ""
				conn := streamingClient.NewConnection(req)
				conn.SubscribeMessages(func(sse sse.Event) {
					print(responseTokens)

					if sse.Data == "[DONE]" { // OpenRouter stream completion message
					} else {
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
					panic(err)
				}
			}
		}

		return err
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

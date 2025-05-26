package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bootun/cosmica/agent/common"
	"github.com/bootun/cosmica/config"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	spacemanCfg := cfg.Agents.Spaceman
	spaceman, err := common.NewSpaceMan(context.Background(), openai.ChatModelConfig{
		BaseURL: spacemanCfg.BaseURL,
		APIKey:  spacemanCfg.Token,
		Model:   spacemanCfg.ModelID,
	})
	if err != nil {
		panic(err)
	}

	var history []*schema.Message
	for {
		question := readUserQuestion()
		ctx := context.Background()
		newHis, err := spaceman.HandleQuestion(ctx, question, history)
		if err != nil {
			log.Fatalf("failed to handle question: %v", err)
		}
		history = newHis
	}
}

func readUserQuestion() string {
	fmt.Printf("> ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

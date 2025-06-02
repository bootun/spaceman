package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bootun/cosmica/agent/common"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	spaceman, err := common.NewSpaceMan(ctx)
	if err != nil {
		log.Fatalf("create spaceman: %v", err)
	}
	var history []*schema.Message
	for {
		question := readUserQuestion()
		newHis, err := spaceman.HandleQuestion(ctx, question, history)
		if err != nil {
			log.Fatalf("handle question: %v", err)
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

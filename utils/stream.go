package utils

import (
	"io"
	"log"

	"github.com/cloudwego/eino/schema"
)

func DealStream(sr *schema.StreamReader[*schema.Message], onRecv func(message *schema.Message)) (*schema.Message, error) {
	defer sr.Close()
	msgs := make([]*schema.Message, 0, 100)
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			// 流式输出结束
			return schema.ConcatMessages(msgs)
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		onRecv(message)
		msgs = append(msgs, message)
	}
}

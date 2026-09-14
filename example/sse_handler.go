package main

import (
	"context"
	"fmt"
	"time"
)

type StreamReq struct {
	Room string `query:"room"`
}

type ChatEvent struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func chatHandler(ctx context.Context, req *StreamReq, send func(name string, data ChatEvent) error) error {
	users := []string{"alice", "bob", "charlie"}
	messages := []string{"hello", "how are you?", "goodbye"}

	for i := range 5 {
		event := ChatEvent{
			User:    users[i%len(users)],
			Message: fmt.Sprintf("[%s] %s #%d", req.Room, messages[i%len(messages)], i+1),
		}

		if err := send("message", event); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

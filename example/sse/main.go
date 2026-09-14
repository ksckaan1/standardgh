package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ksckaan1/standardgh"
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

	for i := 0; i < 5; i++ {
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

func main() {
	http.HandleFunc("GET /stream", standardgh.GHforSSE(5*time.Second, chatHandler))

	fmt.Println("SSE server listening on :3040")
	fmt.Println("Try: curl -N 'http://localhost:3040/stream?room=general'")
	log.Fatal(http.ListenAndServe(":3040", nil))
}

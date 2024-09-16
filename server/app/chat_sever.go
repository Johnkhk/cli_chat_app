package app

import (
	"io"
	"sync"
	"time"

	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
	// Store connections to active clients
	clients map[string]chan *chat.MessageResponse
	mu      sync.Mutex
}

func (s *ChatServiceServer) StreamMessages(stream chat.ChatService_StreamMessagesServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Process incoming encrypted messages here
		// (decrypt, store, etc.)

		// Send response back to the sender
		response := &chat.MessageResponse{
			MessageId: req.MessageId,
			Status:    "delivered",
			Timestamp: time.Now().String(),
		}
		stream.Send(response)
	}
}

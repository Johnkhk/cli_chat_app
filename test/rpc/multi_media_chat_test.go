package rpc

import (
	"context"
	"crypto/sha256"
	"os"
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/lib"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
	utils "github.com/johnkhk/cli_chat_app/test"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

func ComputeSHA256(data []byte) [32]byte {
	return sha256.Sum256(data)
}
func sendAndVerifyMultiMediaMessage(t *testing.T, sender *app.RpcClient, receiver *app.RpcClient, message []byte, expectedType chat.EncryptionType, fileOpts *lib.SendMessageOptions) {
	err := sender.ChatClient.SendMessage(context.Background(), receiver.CurrentUserID, 0, message, fileOpts)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	select {
	case msg := <-receiver.ChatClient.MessageChannel:
		if msg.EncryptionType != expectedType {
			t.Fatalf("Expected %v message, but received: %v", expectedType, msg.EncryptionType)
		}
		decryptedBytes, err := receiver.ChatClient.DecryptMessage(context.Background(), msg)
		if err != nil {
			t.Fatalf("Failed to decrypt message: %v", err)
		}

		// Compare SHA-256 hashes instead of raw bytes
		if ComputeSHA256(decryptedBytes) != ComputeSHA256(message) {
			t.Fatalf("Decrypted media does not match the original message")
		}

		// Write the decrypted bytes to a file
		// decryptedFilePath := fmt.Sprintf("decrypted_%s", fileOpts.FileName)
		// err = os.WriteFile(decryptedFilePath, decryptedBytes, 0644)
		// if err != nil {
		// 	t.Fatalf("Failed to write decrypted file: %v", err)
		// }

		// Compare the file type, size, and name
		if fileOpts.FileType != msg.FileType {
			t.Fatalf("Expected file type %s, but got: %s", fileOpts.FileType, msg.FileType)
		}
		if fileOpts.FileSize != msg.FileSize {
			t.Fatalf("Expected file size %d, but got: %d", fileOpts.FileSize, msg.FileSize)
		}
		if fileOpts.FileName != msg.FileName {
			t.Fatalf("Expected file name %s, but got: %s", fileOpts.FileName, msg.FileName)
		}

	}
}

func TestMultiMediaChat(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0]
	client2 := rpcClients[1]

	utils.RegisterAndLoginUser(t, client1, "user1")
	utils.RegisterAndLoginUser(t, client2, "user2")

	utils.WaitForWelcomeMessage(t, client1, "user1")
	utils.WaitForWelcomeMessage(t, client2, "user2")

	// Read a .jpeg file from the specified path
	jpegFilePath := "multi_media_assets/cat.jpeg"
	jpegFileContent, err := os.ReadFile(jpegFilePath)
	if err != nil {
		t.Fatalf("Failed to read .jpeg file: %v", err)
	}
	sendAndVerifyMultiMediaMessage(t, client1, client2, jpegFileContent, chat.EncryptionType_PREKEY, &lib.SendMessageOptions{
		FileType: "image",
		FileSize: uint64(len(jpegFileContent)),
		FileName: "cat.jpeg",
	})

	time.Sleep(2 * time.Second) // Allow processing time

	// Check chat history for User1
	chatHistory1, err := client1.ChatClient.Store.GetChatHistory(client1.CurrentUserID, client2.CurrentUserID)
	if err != nil {
		t.Fatalf("Error getting chat history for User1: %v", err)
	}

	// Check if the message is saved in the chat history of User1
	message := chatHistory1[0]
	if string(message.Media) != string(jpegFileContent) {
		t.Fatalf("Expected message %s, but got: %s", string(jpegFileContent), message.Message)
	}
	if message.FileType != "image" {
		t.Fatalf("Expected file type %s, but got: %s", "image", message.FileType)
	}
	if message.FileSize != uint64(len(jpegFileContent)) {
		t.Fatalf("Expected file size %d, but got: %d", len(jpegFileContent), message.FileSize)
	}

	// Check chat history for User2
	chatHistory2, err := client2.ChatClient.Store.GetChatHistory(client1.CurrentUserID, client2.CurrentUserID)
	if err != nil {
		t.Fatalf("Error getting chat history for User2: %v", err)
	}

	// Check if the message is saved in the chat history of User2
	message = chatHistory2[0]
	if string(message.Media) != string(jpegFileContent) {
		t.Fatalf("Expected message %s, but got: %s", string(jpegFileContent), message.Message)
	}
	if message.FileType != "image" {
		t.Fatalf("Expected file type %s, but got: %s", "image", message.FileType)
	}
	if message.FileSize != uint64(len(jpegFileContent)) {
		t.Fatalf("Expected file size %d, but got: %d", len(jpegFileContent), message.FileSize)
	}
}

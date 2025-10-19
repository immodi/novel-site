package updater

import (
	"context"
	"fmt"
	pb "immodi/novel-site/internal/app/grpc"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h *UpdaterHandler) Updater(w http.ResponseWriter, r *http.Request) {
	err := h.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	grpcConn, err := grpc.NewClient(h.GrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to connect to gRPC server"))
		return
	}
	defer grpcConn.Close()

	client := pb.NewUpdaterServiceClient(grpcConn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopped := false
	var mu sync.Mutex

	writeCh := make(chan string, 32)
	done := make(chan struct{})
	var once sync.Once

	safeClose := func() {
		once.Do(func() {
			close(writeCh)
			close(done)
		})
	}

	go WriteMessageToWs(conn, writeCh, &mu, done, safeClose)
	go HandleWsIncomingMessage(conn, writeCh, &mu, &stopped, client, done, h.messagesQueue)

	// broadcast the history messages to the user
	for _, message := range h.messagesQueue.GetAll() {
		SafeSend(writeCh, done, message, h.messagesQueue)
	}

	stream, err := client.StreamUpdates(ctx, &pb.Empty{})
	if err != nil {
		SafeSend(writeCh, done, fmt.Sprintf("Error opening stream: %v", err), h.messagesQueue)
		safeClose()
		return
	}

	for {
		select {
		case <-done:
			return // connection closed, stop streaming
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				SafeSend(writeCh, done, "Stream ended.", h.messagesQueue)
				safeClose()
				return
			}
			if err != nil {
				SafeSend(writeCh, done, fmt.Sprintf("Stream error: %v", err), h.messagesQueue)
				safeClose()
				return
			}

			mu.Lock()
			s := stopped
			mu.Unlock()
			if !s {
				SafeSend(writeCh, done, msg.Message, h.messagesQueue)
			}
		}
	}
}

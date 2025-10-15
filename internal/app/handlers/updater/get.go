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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *UpdaterHandler) Updater(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
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

	// 🔒 Create a single goroutine that handles all WebSocket writes
	writeCh := make(chan string, 32)
	go func() {
		for msg := range writeCh {
			mu.Lock()
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			mu.Unlock()
			if err != nil {
				break
			}
		}
	}()

	// 🧠 Handle incoming WebSocket messages
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				close(writeCh)
				return
			}

			switch string(msg) {
			case "START":
				mu.Lock()
				stopped = false
				mu.Unlock()

				resp, err := client.StartUpdate(context.Background(), &pb.UpdateRequest{IntervalHours: 1})
				if err != nil {
					writeCh <- fmt.Sprintf("Error starting updater: %v", err)
					continue
				}
				writeCh <- resp.Message

			case "STOP":
				mu.Lock()
				stopped = true
				mu.Unlock()

				resp, err := client.StopUpdate(context.Background(), &pb.Empty{})
				if err != nil {
					writeCh <- fmt.Sprintf("Error stopping updater: %v", err)
					continue
				}
				writeCh <- resp.Message

			default:
				writeCh <- "Unknown command"
			}
		}
	}()

	// 🛰️ Stream updates from gRPC
	stream, err := client.StreamUpdates(ctx, &pb.Empty{})
	if err != nil {
		writeCh <- fmt.Sprintf("Error opening stream: %v", err)
		close(writeCh)
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			writeCh <- "Stream ended."
			break
		}
		if err != nil {
			writeCh <- fmt.Sprintf("Stream error: %v", err)
			break
		}

		mu.Lock()
		s := stopped
		mu.Unlock()
		if !s {
			writeCh <- msg.Message
		}
	}

	close(writeCh)
}

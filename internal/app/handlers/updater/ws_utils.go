package updater

import (
	"context"
	"fmt"
	pb "immodi/novel-site/internal/app/grpc"
	"immodi/novel-site/pkg"
	"sync"

	"github.com/gorilla/websocket"
)

func WriteMessageToWs(conn *websocket.Conn, writeCh chan string, mu *sync.Mutex, done chan struct{}, safeClose func()) {
	defer safeClose()
	for {
		select {
		case <-done:
			return
		case msg, ok := <-writeCh:
			if !ok {
				return
			}
			mu.Lock()
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			mu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

func HandleWsIncomingMessage(
	conn *websocket.Conn,
	writeCh chan string,
	mu *sync.Mutex,
	stopped *bool,
	client pb.UpdaterServiceClient,
	done <-chan struct{},
	messageQueue *pkg.MessageQueue,
) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		switch string(msg) {
		case "START":
			mu.Lock()
			*stopped = false
			mu.Unlock()

			resp, err := client.StartUpdate(context.Background(), &pb.UpdateRequest{IntervalHours: 3})
			if err != nil {
				SafeSend(writeCh, done, fmt.Sprintf("Error starting updater: %v", err), messageQueue)
				continue
			}
			SafeSend(writeCh, done, resp.Message, messageQueue)

		case "STOP":
			mu.Lock()
			*stopped = true
			mu.Unlock()

			resp, err := client.StopUpdate(context.Background(), &pb.Empty{})
			if err != nil {
				SafeSend(writeCh, done, fmt.Sprintf("Error stopping updater: %v", err), messageQueue)
				continue
			}
			SafeSend(writeCh, done, resp.Message, messageQueue)

		default:
			SafeSend(writeCh, done, "Unknown command", messageQueue)
		}
	}
}

func SafeSend(writeCh chan string, done <-chan struct{}, msg string, messageQueue *pkg.MessageQueue) {
	messageQueue.Add(msg)

	// Try to send non-blocking
	select {
	case writeCh <- msg:
		// sent successfully
	case <-done:
		// channel closed
	default:
		// channel full, skip to prevent blocking
	}
}

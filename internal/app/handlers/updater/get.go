package updater

import (
	"context"
	"fmt"
	pb "immodi/novel-site/internal/app/grpc"
	"immodi/novel-site/internal/config"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		adminOrigin := strings.TrimSuffix(config.AdminSiteURL, "/")
		requestOrigin := strings.TrimSuffix(r.Header.Get("Origin"), "/")

		return requestOrigin == adminOrigin
	},
}

func (h *UpdaterHandler) validate(r *http.Request) error {
	cookie, err := r.Cookie("admin_auth_cookie")
	if err != nil {
		return fmt.Errorf("missing auth cookie")
	}

	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return fmt.Errorf("failed to decode cookie")
	}

	token := strings.Trim(decodedValue, `"`)
	if token == "" {
		return fmt.Errorf("missing token in cookie")
	}

	userID, err := pkg.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("could not get the user from the token")
	}

	if user.Role != string(sql.UserRoleAdmin) {
		return fmt.Errorf("user is not an admin")
	}

	return nil
}

func (h *UpdaterHandler) Updater(w http.ResponseWriter, r *http.Request) {
	err := h.validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
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

	go writeMessageToWs(conn, writeCh, &mu, done, safeClose)
	go handleWsIncomingMessage(conn, writeCh, &mu, &stopped, client, done)

	stream, err := client.StreamUpdates(ctx, &pb.Empty{})
	if err != nil {
		safeSend(writeCh, done, fmt.Sprintf("Error opening stream: %v", err))
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
				safeSend(writeCh, done, "Stream ended.")
				safeClose()
				return
			}
			if err != nil {
				safeSend(writeCh, done, fmt.Sprintf("Stream error: %v", err))
				safeClose()
				return
			}

			mu.Lock()
			s := stopped
			mu.Unlock()
			if !s {
				safeSend(writeCh, done, msg.Message)
			}
		}
	}
}

func writeMessageToWs(conn *websocket.Conn, writeCh chan string, mu *sync.Mutex, done chan struct{}, safeClose func()) {
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

func handleWsIncomingMessage(
	conn *websocket.Conn,
	writeCh chan string,
	mu *sync.Mutex,
	stopped *bool,
	client pb.UpdaterServiceClient,
	done <-chan struct{},
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
				safeSend(writeCh, done, fmt.Sprintf("Error starting updater: %v", err))
				continue
			}
			safeSend(writeCh, done, resp.Message)

		case "STOP":
			mu.Lock()
			*stopped = true
			mu.Unlock()

			resp, err := client.StopUpdate(context.Background(), &pb.Empty{})
			if err != nil {
				safeSend(writeCh, done, fmt.Sprintf("Error stopping updater: %v", err))
				continue
			}
			safeSend(writeCh, done, resp.Message)

		default:
			safeSend(writeCh, done, "Unknown command")
		}
	}
}
func safeSend(writeCh chan string, done <-chan struct{}, msg string) {
	select {
	case writeCh <- msg:
	case <-done:
		// channel is closed, skip sending
	}
}

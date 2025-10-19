package updater

import (
	"immodi/novel-site/internal/app/services/profile"
	"immodi/novel-site/pkg"
)

type UpdaterHandler struct {
	GrpcURL        string
	profileService profile.ProfileService
	messagesQueue  *pkg.MessageQueue
}

func NewUpdaterHandler(grpcURL string, profileService profile.ProfileService) *UpdaterHandler {
	return &UpdaterHandler{GrpcURL: grpcURL, profileService: profileService, messagesQueue: pkg.NewMessageQueue(100)}
}

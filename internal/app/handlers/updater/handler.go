package updater

type UpdaterHandler struct {
	GrpcURL string
}

func NewUpdaterHandler(GrpcURL string) *UpdaterHandler {
	return &UpdaterHandler{GrpcURL: GrpcURL}
}

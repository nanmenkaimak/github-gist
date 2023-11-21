package http

import (
	"io/fs"
	"mime"
	"net/http"

	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"github.com/nanmenkaimak/github-gist/swagger"
	"go.uber.org/zap"
)

type EndpointHandler struct {
	gistService gist.UseCase
	logger      *zap.SugaredLogger
}

func NewEndpointHandler(gistService gist.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		gistService: gistService,
		logger:      logger,
	}
}

type swaggerServer struct {
	openApi http.Handler
}

func (h *EndpointHandler) Swagger() http.Handler {
	if err := mime.AddExtensionType(".svg", "image/svg+xml"); err != nil {
		h.logger.Error("AddExtensionType mimetype error: %w", zap.Error(err))
	}

	openApi, err := fs.Sub(swagger.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}

	return &swaggerServer{
		openApi: http.StripPrefix("/swagger/", http.FileServer(http.FS(openApi))),
	}
}

func (sws *swaggerServer) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	sws.openApi.ServeHTTP(w, rq)
}

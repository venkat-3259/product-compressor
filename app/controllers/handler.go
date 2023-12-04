package controllers

import (
	"os"

	"zocket/app/queries"
	"zocket/pkg/configs"
	imageprocessor "zocket/pkg/image_processor"

	"github.com/rs/zerolog"
)

type Handler struct {
	Queries *queries.Queries
	config  *configs.Config
	Process *imageprocessor.ImageProcessor
	logger  zerolog.Logger
}

// OpenDBConnection func for opening database connection.
func NewHandler(config *configs.Config, imageProcess *imageprocessor.ImageProcessor) *Handler {

	logger := zerolog.New(os.Stderr).With().Timestamp().Str("Service", "API Handler").Logger()

	return &Handler{imageProcess.Queries, config, imageProcess, logger}
}

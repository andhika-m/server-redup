package rest

import (
	"redup/internal/usecase/redup"
)

type handler struct {
	redupUsecase redup.Usecase
}

func RedupHandler(redupUsecase redup.Usecase) *handler {
	return &handler{redupUsecase: redupUsecase}
}

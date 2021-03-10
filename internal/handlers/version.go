package handlers

import (
	"github.com/AlexandrGurkin/http_template/internal/ver"
	"github.com/AlexandrGurkin/http_template/models"
	"github.com/AlexandrGurkin/http_template/restapi/operations/version"
	"github.com/go-openapi/runtime/middleware"
)

type VersionHandler struct {
}

func (vh VersionHandler) Handle(_ version.GetVersionParams) middleware.Responder {
	return version.NewGetVersionOK().WithPayload(&models.ResponseVersion{
		Version:   ver.GetVersion(),
		Branch:    ver.GetBranch(),
		Commit:    ver.GetCommit(),
		BuildTime: ver.GetBuildTime(),
	})
}

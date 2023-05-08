package api

import (
	"budgetting/api/http"
	"budgetting/config"
)

func ListenAndServe(config *config.Config) error {
	return http.ListenAndServe(config)
}

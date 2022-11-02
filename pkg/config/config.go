package config

import (
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"os"
)

var (
	DefaultDocumentStage = []string{
		"Sent",
		"Verified",
		"Approved",
	}

	DefaultUser = &entity.User{
		Username: "admin",
		Password: "admin",
		Role:     3,
	}
)

func LoadConfig() map[string]string {
	env := make(map[string]string)

	env["DB_HOST"] = os.Getenv("DB_HOST")
	env["DB_PORT"] = os.Getenv("DB_PORT")
	env["DB_USER"] = os.Getenv("DB_USER")
	env["DB_PASS"] = os.Getenv("DB_PASS")
	env["DB_NAME"] = os.Getenv("DB_NAME")
	env["PORT"] = os.Getenv("PORT")
	env["JWT_SECRET"] = os.Getenv("JWT_SECRET")
	env["QR_PATH"] = os.Getenv("QR_PATH")

	return env
}

package biz

import "github.com/google/uuid"

func CreateUuid() string {
	uuid := uuid.New()
	return uuid.String()
}

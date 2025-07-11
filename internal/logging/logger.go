package logging

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LogError(err error, fields map[string]interface{}) string {
	id := uuid.New().String()
	fields["id"] = id
	logrus.WithFields(fields).WithError(err).Error("Error occurred")

	return id
}

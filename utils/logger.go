package utils

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// Уровень логирования
	Logger.SetLevel(logrus.DebugLevel)

	// Формат вывода: JSON или текст
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

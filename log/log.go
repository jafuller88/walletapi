package log

import (
	"fmt"
	"time"
	"walletapi/config"
)

// Msgf logs
func Msgf(logLevel int, logMessage string, a ...interface{}) {
	if logLevel <= config.LogLevel {
		t := time.Now()
		logM := fmt.Sprintf("%v: %v", t.Format(time.RFC3339), logMessage)
		fmt.Printf(logM, a...)
	}
}

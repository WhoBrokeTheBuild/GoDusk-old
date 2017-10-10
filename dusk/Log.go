package dusk

import (
	"fmt"
	"os"
)

func LogVerbose(format string, a ...interface{}) (int, error) {
	return fmt.Printf("[VERB] "+format+"\n", a...)
}

func LogLoad(format string, a ...interface{}) (int, error) {
	return fmt.Printf("[LOAD] "+format+"\n", a...)
}

func LogInfo(format string, a ...interface{}) (int, error) {
	return fmt.Printf("[INFO] "+format+"\n", a...)
}

func LogWarn(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, "[WARN] "+format+"\n", a...)
}

func LogError(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, "[ERRO] "+format+"\n", a...)
}

package utils

import (
	"fmt"
	"time"
)

func Log(logObject any) {
	cTime := time.Now().Local().Format("2006/01/02 15:04:05")
	fmt.Printf("[%s] %v\n", cTime, logObject)
}

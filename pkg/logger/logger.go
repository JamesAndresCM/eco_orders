package logger

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "ECO_ORDERS: ", log.LstdFlags|log.Lshortfile)


package public

import (
	"log"
	"os"
)

var (
	Warning *log.Logger
	Info  *log.Logger
	Error *log.Logger
	Success *log.Logger
	Progress *log.Logger
)

func init() {
	Progress = log.New(os.Stdout, "[Progress]::| ",0)
	Warning = log.New(os.Stdout, "[Warning]:: | ", log.Ldate|log.Ltime)
	Success = log.New(os.Stdout, "[Success]:: | ", log.Ldate|log.Ltime)
	Info = log.New(os.Stdout, "[Info]::    | ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[Error]::   | ", log.Ldate|log.Ltime|log.Lshortfile)
}
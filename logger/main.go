package logger

import (
	"log"
	"net"
	"time"
)


type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	remote bool
	conn *net.Conn
}

func NewLogger(remote bool) *Logger {
	
		log.Print("connecting to log server...")
		conn, err := net.DialTimeout("tcp", "127.0.0.1:8081", time.Second * 5)
		if err != nil {
			log.Fatalf("error connecting to remote log server: %v", err)
		}
		log.Printf("connected to log server")


	return &Logger{
		infoLogger:   log.New(conn, "INFO: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		errorLogger: log.New(conn, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		remote: remote,
		conn: &conn,
	}
}

func (logger *Logger) Info(format string, v ...any) {
	logger.infoLogger.Printf(format, v...)
}

func (logger *Logger) Error(format string, v ...any) {
	logger.errorLogger.Printf(format, v...)
}
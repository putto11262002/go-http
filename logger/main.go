package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"time"
)


type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	rw io.WriteCloser
}


type RemoteBufferedWriter struct {
	conn *net.Conn
	buffereWriter *bufio.Writer
}

func NewRemoteBufferedWriter() *RemoteBufferedWriter{
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8081", time.Second * 5)
	if err != nil {
		log.Fatalf("error connecting to remote log server: %v", err)
	}
	return &RemoteBufferedWriter {
		conn: &conn,
		buffereWriter: bufio.NewWriter(conn),
	}
}

func (w *RemoteBufferedWriter) Write(p []byte) (n int, err error) {
	n, err = w.buffereWriter.Write(p)
	if err != nil {
		log.Printf("error writing to connection: %v", err)
	}
	if err := w.buffereWriter.Flush(); err != nil {
		log.Printf("error flushing to connection: %v", err)
	}
	return n, err
}

func (w *RemoteBufferedWriter) Close() (err error) {
	err = w.buffereWriter.Flush()
	if err != nil {
		log.Printf("error flushing to connection: %v", err)
	}

	if err := (*w.conn).Close(); err != nil {
		log.Printf("error closing connection: %v", err)
	}
	
	return err
	
}

func NewLogger(remote bool) *Logger {

	remoteWriter := NewRemoteBufferedWriter()
	w := io.MultiWriter(remoteWriter, os.Stdout)

	return &Logger{
		infoLogger:   log.New(w, "INFO: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		errorLogger: log.New(w, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		rw: remoteWriter,
	}
}

func (logger *Logger) Close() error {
	err := logger.rw.Close()
	if err != nil {
		return fmt.Errorf("closing log writer: %w", err)
	}
	return nil
}



func (logger *Logger) Info(format string, v ...any) {
	logger.infoLogger.Printf(format, v...)
}

func (logger *Logger) Error(format string, v ...any) {
	logger.errorLogger.Printf(format, v...)
}
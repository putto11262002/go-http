package logserver

import (
	"fmt"
	"log"
	"net"
	"time"
)

type LogServer struct {
	addr string
	listener net.Listener
}

func NewLogServer(addr string) *LogServer {
	return &LogServer{
		addr: addr,
	}
}

func (logServer *LogServer) Run(){
	l, err :=  net.Listen("tcp", logServer.addr)

	if err != nil {
		log.Fatalf("error starting tcp listenser: %v", err)
	}

	logServer.listener = l
	log.Printf("log server listening on %v", logServer.addr)

	defer l.Close()

	for {
		con, err := logServer.listener.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		go func (conn net.Conn)  {
			defer func ()  {
				log.Printf("%s disconnected", conn.RemoteAddr())
				 conn.Close()
			}()
			conn.SetReadDeadline(time.Now().Add(time.Hour * 24))
			
			log.Printf("%s connected", conn.RemoteAddr())
			buffer := make([]byte, 1024)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					log.Printf("Error reading from %s", conn.RemoteAddr())
					return
				}

				logMessage := string(buffer[:n])
				fmt.Printf("%s: %s", conn.RemoteAddr(), logMessage)
			}
			
		}(con)
		
	}

}


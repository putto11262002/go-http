package logserver

import (
	"bufio"
	"log"
	"net"
	"os"
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
			// Set read deadline to 24 hours
			conn.SetReadDeadline(time.Now().Add(time.Hour * 24))

			file, err := os.OpenFile("app.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("error opening log file: %v", err)
				return 
			}

			writer := bufio.NewWriter(file)
			reader := bufio.NewReader(conn)
			

			defer func ()  {
				
				writer.Flush()
				 file.Close()
				 conn.Close()
				 log.Printf("%s disconnected", conn.RemoteAddr())
			}()

		
			
			for {
				var conClose bool
				incoming, err := reader.ReadBytes('\n')
				
				if err != nil {
					// if err == io.EOF {
					// 	conClose = true
					// }
					log.Printf("Error reading from %s", conn.RemoteAddr())
					conClose = true
				}

				

				_, err = writer.Write(incoming)
				if err != nil {
					log.Fatalf("error writing to log file: %v\n", err)
				}

				if conClose {
					break
				}

			}
			
		}(con)
		
	}

}


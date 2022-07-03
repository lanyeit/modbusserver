package mbserver

import (
	"io"
	"log"

	"github.com/goburrow/serial"
)

// ListenRTU starts the Modbus server listening to a serial device.
// For example:  err := s.ListenRTU(&serial.Config{Address: "/dev/ttyUSB0"})
func (s *Server) ListenRTU(serialConfig *serial.Config) (err error) {
	port, err := serial.Open(serialConfig)

	//log.Fatalf("ccc %s: %v\n", serialConfig.Address, err)
	if err != nil {
		log.Fatalf("failed to open %s: %v\n", serialConfig.Address, err)
	}
	s.ports = append(s.ports, port)

	s.portsWG.Add(1)
	go func() {
		defer s.portsWG.Done()
		s.acceptSerialRequests(port)
	}()
	//log.Fatal(s.portsWG)
	return err
}

func (s *Server) acceptSerialRequests(port serial.Port) {

SkipFrameError:
	for {
		select {
		case <-s.portsCloseChan:
			return
		default:
		}

		buffer := make([]byte, 1024)

		bytesRead, err := port.Read(buffer)
		//bytesRead = append(bytesRead,1)
		//log.Printf("error: %v\n", buffer)
		log.Printf("SDFDS222: %v\n", bytesRead)
		if err != nil {
			if err != io.EOF {
				log.Printf("serial read error %v\n", err)
			}
			return
		}

		if bytesRead != 0 {
			//bytesRead = append(bytesRead,1)
			// Set the length of the packet to the number of read bytes.
			//buffer = append([]byte{1},buffer...)
			//buffer = append(buffer,1)
			packet := buffer[:bytesRead]
			//pLen := len(packet)
			log.Printf("SDFDS2: %v\n", packet)
			frame, err := NewRTUFrame(packet)
			if err != nil {
				log.Printf("bad serial frame error %v\n", err)
				log.Printf("bad serial frame error %v\n", bytesRead)
				//The next line prevents RTU server from exiting when it receives a bad frame. Simply discard the erroneous
				//frame and wait for next frame by jumping back to the beginning of the 'for' loop.
				log.Printf("Keep the RTU server running!!\n")
				continue SkipFrameError
				//return
			}

			request := &Request{port, frame}

			s.requestChan <- request
		}
	}
}

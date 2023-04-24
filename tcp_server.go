package main

import (
	"github.com/Nicoletapa/is105sem03/mycrypt"
	"io"
	"log"
	"net"
	"sync"
	"fmt"
	
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:9800")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for lokke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					kryptertMelding := mycrypt.Krypter([]rune(string(dekryptertMelding)),mycrypt.ALF_SEM03, len( mycrypt.ALF_SEM03)-4)
        				log.Println("Kryptert melding: ", string(kryptertMelding))
       					 _, err = conn.Write([]byte(string(kryptertMelding)))

					switch msg := string(dekryptertMelding); msg {
					case "ping":
						_, err = c.Write([]byte("pong"))
					case "Kjevik":
					var celsius float64 = 6
					
					fahr :=  (celsius * 1.8) +32
					response := fmt.Sprintf("Kjevik;SN39040;18.03.2022 01:50;%.2f", fahr)
					_, err = c.Write([]byte(response))
					default:
						_, err = c.Write([]byte(string(kryptertMelding)))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for lkke
					}

				}
			}(conn)
		}
	}()
	wg.Wait()

}

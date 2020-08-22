package main

import (
	"fmt"
	"lib/grinds"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	port := ":" + arguments[1]
	writersServer, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer writersServer.Close()

	storageManager := grinds.NewStorageManager()
	rand.Seed(time.Now().Unix())
	for {
		conn, err := writersServer.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go storageManager.handleWriterConnection(conn)
	}
}

// func handleConnection(conn net.Conn) {
// 	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
// 	for {
// 		netData, err := bufio.NewReader(conn).ReadString('\n')
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		temp := strings.TrimSpace(string(netData))
// 		if temp == "STOP" {
// 			break
// 		}

// 		fmt.Println("Received: " + temp)

// 		result := "Reply to '" + temp + "' is: " + strconv.Itoa(rand.Int()) + "\n"
// 		conn.Write([]byte(string(result)))

// 	}
// 	conn.Close()
// }

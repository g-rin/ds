// Package ds provides primitives for writing into and read from
// a record storage.
package ds

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Record struct {
	Timestamp time.Time
	Author    string
	Type      int32
	Size      int64
	Data      []byte
}

type StorageWriter struct {
	Address  string
	Channel  chan Record
	NextSize int64
}

type StorageManager struct {
	mutex   sync.RWMutex
	Writers []string
	Channel chan Record
	buffer  []Record
}

func NewStorageManager() *StorageManager {
	var manager StorageManager
	manager.Channel = make(chan Record)
	return &manager
}

func (manager *StorageManager) handleWriterConnection(conn net.Conn) {
	var writer StorageWriter
	writer.Address = conn.RemoteAddr().String()
	fmt.Printf("A writer connected from %s\n", writer.Address)
	writer.Channel = manager.Channel
	manager.addWriter(writer.Address)
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		data := strings.TrimSpace(string(netData))
		fmt.Println("Received: " + data)
		result := "Reply to '" + data + "' \n"
		conn.Write([]byte(string(result)))
	}
	manager.removeWriter(writer.Address)
	conn.Close()
}

func (manager *StorageManager) addWriter(address string) {
	manager.mutex.Lock()
	manager.Writers = append(manager.Writers, address)
	manager.mutex.Unlock()
}

func (manager *StorageManager) removeWriter(address string) {
	manager.mutex.Lock()
	index := 0
	manager.Writers = append(manager.Writers[:index],
		manager.Writers[index+1:]...)
	manager.mutex.Unlock()
}

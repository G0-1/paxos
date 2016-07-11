package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/n0_1/paxos/protocol"
)

var wg sync.WaitGroup

func main() {
	for i := 1; i <= 9999; i++ {
		conn, err := net.Dial("tcp", ":80")
		if err != nil {
			log.Println("dial ser error: ", err)
			continue
		}
		wg.Add(1)
		go handleConn(conn, uint64(i))
	}
	wg.Wait()
}

func handleConn(conn net.Conn, id uint64) {
	id = 1111111111111111110 + id
	conn.Write(protocol.AsmCmd("cre", strconv.Itoa(int(id))))

	for {
		start := time.Now()

		conn.Write(protocol.AsmCmd("inc", strconv.Itoa(int(id))))

		r := bufio.NewReader(conn)
		recv, rerr := r.ReadSlice('\n')
		if rerr != nil {
			log.Println("read error: ", rerr)
			break
		}
		recv = recv[:len(recv)-1]
		if len(recv) > 0 && recv[len(recv)-1] == '\r' {
			recv = recv[:len(recv)-1]
		}
		log.Println("recv content: ", string(recv))

		d := time.Now().Sub(start)
		log.Println(strconv.Itoa(int(id))+" inc  spend: ", d)
	}

	wg.Done()
}

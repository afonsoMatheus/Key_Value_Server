// Implementation of a KeyValueServer. Students should write their code in this file.
//change this

package p1

import (
	"net"
	"bufio"
	"strings"
	"strconv"
)

type Cliente struct{
	conn net.Conn
	channelReader chan []byte
}

type keyValueServer struct {

	listClient[] *Cliente
	listen net.Listener
	channelList chan Cliente
	channelReader chan []string
	channelBroad chan []byte
	channelCount chan int
	channelQuit chan Cliente

	err error

}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {

	serv := &keyValueServer{
		listClient:make([]*Cliente,0),
		listen:nil,
		channelList:make(chan Cliente),
		channelReader:make(chan []string),
		channelBroad:make(chan []byte),
		channelCount:make(chan int),
		channelQuit:make(chan Cliente),
	    err: nil}

    return serv
}

func (kvs *keyValueServer) Start(port int) error {

	kvs.listen, kvs.err = net.Listen("tcp", "localhost:" + strconv.Itoa(port))

	if kvs.err != nil {
		return kvs.err
	}

	go runner(kvs)
	go managerClient(kvs)
	go managerServerBank(kvs)



	return nil
}

func (kvs *keyValueServer) Close() {

	kvs.listen.Close()
}

func (kvs *keyValueServer) Count() int {

	kvs.channelCount <- 1

	return <-kvs.channelCount
}

func runner(kvs *keyValueServer){

	for{

		conn, err := kvs.listen.Accept()

		if err != nil {
			return
		}

		client := Cliente{conn, make(chan []byte ,500)}

		kvs.channelList <- client

	}

}

func managerClient (kvs *keyValueServer){

	for{
		select {

		case ch1 := <-kvs.channelList:
			kvs.listClient = append(kvs.listClient, &ch1)
			go handleConn(ch1, kvs)
			go mywriter(ch1)

		case ch1 := <-kvs.channelBroad:
			for _, c := range kvs.listClient{
				if len(c.channelReader) != 500 {
					c.channelReader <- ch1
				}
			}

		case ch1 := <-kvs.channelQuit:
			for i, c := range kvs.listClient{
				if c.conn == ch1.conn{
					kvs.listClient = append(kvs.listClient[:i],kvs.listClient[i+1:]...)
					c.conn.Close()
				}
			}

		case <-kvs.channelCount:
			kvs.channelCount <-len(kvs.listClient)

		}
	}
}

func handleConn(client Cliente, kvs *keyValueServer){

	fin := bufio.NewReader(client.conn)

	for{
		finn, err := fin.ReadBytes('\n')

		if err != nil {
			kvs.channelQuit <- client
			return
		}

		str := strings.Split(string(finn), ",")

		if str[0] == "put"{

			msg := make([]string,2)

			msg[0] = str[1]
			msg[1] = str[2]

			kvs.channelReader <- msg

		}

		if str[0] == "get"{

			msg := str[1]

			msg2 := strings.Split(msg, "\n")

			msg3 := strings.Split(msg2[0]," ")

			kvs.channelReader <- msg3

		}

	}

}

func managerServerBank(kvs *keyValueServer){

	for{
		select {
		case ch1 := <-kvs.channelReader:

			if len(ch1) == 1{

				val := get(ch1[0])

				msg := ch1[0]+","+string(val)

				kvs.channelBroad <- []byte(msg)

			}

			if len(ch1) == 2{

				put(ch1[0], []byte(ch1[1]))

			}

		}

	}
}


func mywriter(client Cliente){

	for{
		select {
		case ch1 := <-client.channelReader:
			client.conn.Write(ch1)
		}
	}
}
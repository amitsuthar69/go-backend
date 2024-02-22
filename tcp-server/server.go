/*
CREATING A CUSTOM TCP SERVER:

1. create a TCP listener and reserve a port.
2. Accept client connection requests on this port to establish a connection, (blocking => wait and block the execution until the client connects.)
3. Once the connection is established, we can do the following things with the connection:
  -> read from the request, (blocking)
	-> write back a response, (blocking)
	-> close the connection.
==> In this code, we spin off a new goroutine for every new incoming requests.

NOTE:
- This code below demonstrates basic network programming principles but is not recommended for building actual web applications.
- It is not suitable for building production-grade web applications.
- While the code below can be useful for educational and specific scenarios, for most practical web development needs,
  using the net/http package is a safer and more efficient approach.
*/

package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var start = time.Now()

func do(conn net.Conn) {
	buffer := make([]byte, 1024) // this buffer is a temporary storage of 1kb in memory to hold the data being read.

	_, err := conn.Read(buffer) // conn.Read() returns number of bytes read and error.
	if err != nil {
		log.Fatal("error reading from connection: ", err)
	}

	time.Sleep(time.Second * 8) // fake delay

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHey Client!\r\n")) // responding with a HTTP Status code 200 OK
	defer conn.Close()
}

func main() {
	l, err := net.Listen("tcp", ":4221") // creating a TCP listener which listens on port 4221
	if err != nil {
		log.Fatal("Failed binding to port 4221", err.Error())
	}

	for {
		fmt.Println("waiting for a client to connect...")

		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting connection: ", err.Error())
		}

		fmt.Println("client connected at: ", time.Since(start))

		go do(conn) // remove go keyword to make this function call single threaded
	}
}

/*
In above code, we fork a new thread to handle the request so if there are 'n' requests,
we would have 'n' threads handling them.

This looks awesome so what's the problem then?

Well, what happens when 'n' shoots up?
  -> we would have large no of threads running.
  -> resource consuming.
  -> overwhelms the hardware.
  -> It can hang and crash the machine.

Honce we can't just have threads spinning up every now and then.
- We need to limit maximum numbers of thread we create.
- This is exactly what thread pool solves.
*/

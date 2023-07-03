package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	unixAddr, err := net.ResolveUnixAddr("unix", "temp/demo.sock")

	if err != nil {

		fmt.Println(err)

		return
	}

	conn, cErr := net.DialUnix("unix", nil, unixAddr)

	if cErr != nil {

		fmt.Println(cErr)

		return
	}

	//conn.SetDeadline()

	defer func() {

		conn.Close()
		//conn.Close()

	}()

	conn.Write([]byte("hello\n"))

	reader := bufio.NewReader(conn)

	message, rErr := reader.ReadString('\n')

	if rErr != nil {

		fmt.Println(rErr)

		return
	}

	fmt.Println(message)

}

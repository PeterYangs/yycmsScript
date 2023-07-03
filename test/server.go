package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	os.Remove("temp/demo.sock")

	unixAddr, err := net.ResolveUnixAddr("unix", "temp/demo.sock")

	if err != nil {

		fmt.Println(err)

		return
	}

	unixListener, lErr := net.ListenUnix("unix", unixAddr)

	if lErr != nil {

		fmt.Println(lErr)

		return
	}

	defer unixListener.Close()

	for {

		unixConn, aErr := unixListener.AcceptUnix()

		if aErr != nil {

			fmt.Println(aErr)

			continue
		}

		go func(u *net.UnixConn) {

			defer u.Close()

			reader := bufio.NewReader(u)

			for {

				message, rErr := reader.ReadString('\n')

				if rErr != nil && rErr != io.EOF {

					fmt.Println(rErr)

					return
				}

				fmt.Println(message)

				u.Write([]byte("nice\n"))

			}

		}(unixConn)

	}

}

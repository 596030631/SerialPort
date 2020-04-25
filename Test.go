package main

import "fmt"

func main() {
	heads := []byte{0x76, 0x7D}

	for i, a := range heads {
		fmt.Println(i, " -> ", a)
	}
}

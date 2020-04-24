package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"time"
)

var (
	/* start byte */
	VMC_SIGN               = "78"
	VMC_POLL               = "76"
	VMC_OUT_GOOD           = "7c"
	VMC_CHANNEL_RUN_INFO   = "79"
	VMC_MACHINE_RUN_INFO   = "7d"
	VMC_SYSTEM_STATE       = "73"
	VMC_UUID               = "71"
	VMC_SYSTEM_CONFIG      = "7a"
	VMC_CHANNEL_GOOD_INFO  = "7e"
	VMC_CHANNEL_PRICE_INFO = "7f"
	VMC_SUMMARY_OF_SALES   = "7b"
	VMC_PICK_UP_CODE       = "74"
	VMC_CHANNEL_SALE_COUNT = "75"
	VMC_CHANNEL_STATE      = "72"

	/* special words */
	VMC_POLL_SUCCESS    = "760076"
	VMC_POLL_OUTGOOG    = "7603"
	VMC_OUT_GOOD_REFUSE = "76158b"
)

func main() {

	// 接收线程为主
	chReader := make(chan string)
	go Reader(chReader)

	// 发送线程被动回复
	chWriter := make(chan string)
	go Writer(chWriter)

	c := &serial.Config{Name: "COM2", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Print(err)
	}

	for true {

		time.Sleep(1 * time.Second)
		n, err := s.Write([]byte("test"))
		if err != nil {
			fmt.Print(err)
		}

		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%q", buf[:n])
	}

}

/* read */
func Reader(ch chan string) {

}

/* write */
func Writer(ch chan string) {
	ch <- "eeee"
}

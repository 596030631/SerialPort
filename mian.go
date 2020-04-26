package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"regexp"
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

	heads = []byte{0x76, 0x79, 0x7D, 0x7C}
	HWL   = make(map[byte]int) // header with length

)

func main() {

	//HEAD := []byte{0xEF, 0xEE, 0xFE}
	//HEADSTR := "efeefe"

	HWL[0x71] = 86
	HWL[0x7e] = 108
	HWL[0x7f] = 204
	HWL[0x76] = 20
	HWL[0x7b] = 162
	HWL[0x7d] = 40
	HWL[0x79] = 120
	HWL[0x7a] = 44

	/* 打开串口 */
	ch := make(chan io.ReadWriteCloser)
	go openSerial(ch)
	snake := <-ch

	/* 读取数据 */
	stream := make([]byte, 0)
	for true {

		fmt.Println("----------------------------------------------------------------------------------------")
		//chReader := make(chan []byte)
		//Reader(chReader, snake)
		buf := make([]byte, 256)
		n, err := snake.Read(buf[:])
		if err != nil {
			fmt.Println(err)
		}
		bs := buf[:n] // 抛弃未写入部分
		//bs := <- chReader
		for _, n := range bs {
			stream = append(stream, n)
		}

		s := hex.EncodeToString(stream)
		rule, _ := regexp.Compile(".{2}efeefe([0-9]|[a-z])*?efeefe") // 非贪婪模式
		var orders []string
		for true {
			index := rule.FindStringIndex(s)
			if index == nil {
				break
			}
			s = s[:index[1]-8]
			stream = stream[(index[1]-8)/2-1:]
			order := s[(index[0]):(index[1] - 8)]
			orders = append(orders, order)
		}

		if orders == nil {
			time.Sleep(50)
			continue
		}

		fmt.Println("orders -> ", orders)

		for _, value := range orders {
			/* 数据分支 */
			by, _ := hex.DecodeString(value[:2])
			if by == nil {
				continue
			}
			switch by[0] { // "76" -> 0x76
			case 0x76:
				/* write */
				if len(value) == 4 {
					chWriter := make(chan string)
					food := []byte{0x76, 0x00, 0x76}
					Writer(chWriter, snake, food)
					break
				}
			}
		}

		//temp := heads
		//var location []int

		//for i, n := range stream {
		//	for _, m := range temp {
		//		if n == m {
		//			if len(location) > 0 {
		//				ord := stream[location[len(location)-1]] // 最新一个位置对应的开头
		//				//fmt.Println("ORD -> ",ord)
		//				L := HWL[ord] // 获得最新开头的长度
		//				//fmt.Println("HWL -> ",HWL[ord])
		//
		//				if L <= i - location[len(location)-1]{ // 不在指令长度内，是开头
		//					location = append(location, i)
		//				}
		//			} else {
		//				location = append(location, i)
		//			}
		//			break
		//		}
		//	}
		//}

		//fmt.Println(location)
		// 提示  根据指令长度来判断，不能只根据开头 错误示范：
		// 7defeefe2000000000000000000078000000000000000000000000000000000000000000000000000000000000000000000000000

		//num := len(location)
		//var orders [][]byte
		//
		//if num > 1 {
		//	for s := 1; s < num; s++ {
		//		orders = append(orders, stream[location[s-1]:location[s]+1])
		//	}
		//	stream = stream[location[num-1]:]
		//}

		time.Sleep(200 * time.Millisecond) // 休眠
	}
}

func openSerial(ch chan io.ReadWriteCloser) {
	/* 打开串口 */
	c := &serial.Config{Name: "COM3", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("Open serial failed err -> ", err)
		return
	}
	ch <- s
}

/* read */
func Reader(ch chan []byte, reader io.ReadWriteCloser) {
	buf := make([]byte, 128)

	_, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println(err)
	}

	//index := bytes.IndexByte(buf, 0)

	ch <- buf[:]
}

/* write */
func Writer(ch chan string, writer io.ReadWriteCloser, data []byte) {
	_, err := writer.Write(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(data))
}

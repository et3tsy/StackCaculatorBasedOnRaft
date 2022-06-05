package main

import (
	"bufio"
	"client/network"
	"client/settings"
	"client/validate"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// load viper
	if err := settings.Init(); err != nil {
		log.Panicf("init settings failed, err:%v\n", err)
	}

	network.Init()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		// get commands from stdin, line by line
		input, _ := inputReader.ReadString('\n')
		input = strings.Replace(input, "\r\n", "", -1)

		if input == "q" {
			return
		}
		arrStr := strings.Split(input, " ")

		cmd := arrStr[0]
		args := make([]int64, 0)

		for i := 1; i < len(arrStr); i++ {
			v, err := strconv.ParseInt(arrStr[i], 10, 64)
			if err != nil {
				break
			}
			args = append(args, v)
		}

		// to ensure the format of the arguments pass to
		// server is valid
		if !validate.Check(cmd, args) {
			fmt.Println("Format error")
			continue
		}

		// to excute the command
		resp, err := network.Excute(cmd, args)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		fmt.Println(resp.Message)
	}
}

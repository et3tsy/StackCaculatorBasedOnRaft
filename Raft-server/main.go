package main

import (
	"fmt"
	"log"
	"os"
	"server/calculator"
	"server/models"
	"server/network"
	"server/raft"
	"server/settings"
	"strconv"

	"github.com/spf13/viper"
)

func main() {
	// load viper
	if err := settings.Init(); err != nil {
		log.Panicf("init settings failed, err:%v\n", err)
	}

	// need 1 argument, to identify the id of the server
	if len(os.Args) < 2 {
		log.Panic("need cmd argument\n")
	}

	if v, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Panicf("arguments error,err: %v", err)
	} else {
		network.Me = v
	}

	// to get the servers' address
	network.Network = viper.GetString("rpc.Network")
	network.Peers = viper.GetStringSlice("rpc.addr")

	applyCh := make(chan models.ApplyMsg)

	// create raft
	rf := raft.Make(network.Peers, network.Network, network.Me, applyCh)

	// set RPC listener
	if err := network.InitRPC(rf); err != nil {
		log.Panicf("init RPC settings failed, err: %v\n", err)
	}

	// create stack calculator
	c := calculator.Make(rf, applyCh)

	// set client listener
	err := network.InitClient(c)
	if err != nil {
		log.Panicf("init client settings failed, err: %v\n", err)
	}

	for {
		s := ""
		fmt.Scan(&s)
		if s == "q" {
			break
		}
		if s == "a" {
			r := models.Request{
				Instruction: "create",
				Params:      []int64{},
			}
			c.Raft.Start(r)
		}
		fmt.Println("====================================================")
		term, _ := c.Raft.GetState()
		leader := c.Raft.GetLeader()
		fmt.Printf("curTerm: %v\n", term)
		fmt.Printf("curLeader: %v\n", leader)
		fmt.Printf("stackNum: %v\n", c.GetNum())
		fmt.Println("====================================================")
	}
}

package network

import (
	"bufio"
	"client/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/spf13/viper"
)

var (
	Addr []string
	Net  string
	Size int
)

const (
	requestDelay = time.Millisecond * 200
)

func Init() {
	rand.Seed(time.Now().Unix())
	Size = viper.GetInt("client.size")
	Addr = viper.GetStringSlice("client.addr")
	Net = viper.GetString("client.network")
}

// Use TCP to send requests.
func postCommand(req []byte, hostID int) ([]byte, error) {
	var buf [128]byte

	if hostID >= len(Addr) {
		return nil, fmt.Errorf("host ID out of bound")
	}

	conn, err := net.Dial(Net, Addr[hostID])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	n, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println("read from server failed, err:", err)
		return nil, err
	}
	return buf[:n], nil
}

// Excute the command.
func Excute(cmd string, args []int64) (models.Response, error) {
	req := models.Request{
		Instruction: cmd,
		Params:      args,
	}
	resp := models.Response{}

	// use JSON format
	msg, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	hostID := -1
	for {
		if hostID == -1 {
			hostID = rand.Intn(Size)
		}
		reply, err := postCommand(msg, hostID)
		if err != nil {
			return resp, err
		}
		if err = json.Unmarshal([]byte(reply), &resp); err != nil {
			return resp, err
		}
		if resp.Message == "no leader" {
			hostID = resp.Value
		} else {
			break
		}
		time.Sleep(requestDelay)
	}
	return resp, nil
}

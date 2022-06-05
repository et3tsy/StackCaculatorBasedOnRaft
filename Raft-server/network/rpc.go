package network

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"server/raft"

	"github.com/spf13/viper"
)

var (
	addr    string
	Network string
	Me      int
	Peers   []string
)

// Initialize servers, listening RPCs from peers.
func InitRPC(rf *raft.Raft) error {
	// load size, to validate Me
	size := viper.GetInt("rpc.size")
	if size <= Me {
		return fmt.Errorf("cmd arguments out of bound")
	}

	addr = Peers[Me]

	// register RPC service
	rpc.RegisterName("Raft", rf)

	// listen at addr
	listener, err := net.Listen(Network, addr)
	if err != nil {
		log.Fatalf("TCP listen error:%v", err)
	}

	// connect to peers
	for i := 0; i < size*10; i++ {
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Fatal("Accept error:", err)
				}
				rpc.ServeConn(conn)
				conn.Close()
			}
		}()
	}

	return nil
}

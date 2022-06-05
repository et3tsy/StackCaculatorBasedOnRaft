package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"server/calculator"
	"server/models"

	"github.com/spf13/viper"
)

const (
	workerNum = 3
)

var (
	stackCal *calculator.Calculator
)

// 定义协程池类型
type Pool struct {
	worker_num  int           // 协程池最大worker数量,限定Goroutine的个数
	JobsChannel chan net.Conn // 协程池内部的任务就绪队列
}

// 创建一个协程池
func NewPool(cap int) *Pool {
	p := Pool{
		worker_num:  cap,
		JobsChannel: make(chan net.Conn),
	}
	return &p
}

// worker处理函数
func process(conn net.Conn) {
	defer conn.Close()

	// read the command from client in JSON format
	reader := bufio.NewReader(conn)
	var buf [128]byte
	n, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println("read from client failed, err:", err)
		return
	}

	msg := buf[:n]

	request := models.Request{}
	response := models.Response{
		Message: "",
		Success: false,
	}

	// unmarshal the request in JSON format and push it to manager
	if err := json.Unmarshal(msg, &request); err != nil {
		response.Message = "data format error"
	} else {
		fmt.Println(request)
		response = stackCal.Manage(request)
	}

	msg, _ = json.Marshal(response)
	conn.Write(msg)
}

// 协程池中每个worker的功能
func (p *Pool) worker() error {
	//worker不断的从JobsChannel内部任务队列中拿Conn
	for conn := range p.JobsChannel {
		//如果拿到Conn,则执行对应处理
		process(conn)
	}
	return nil
}

// 协程池Pool开始工作
func (p *Pool) Run(connNet, listenAddr string) {
	// 设置监听端口
	listen, err := net.Listen(connNet, listenAddr)
	if err != nil {
		fmt.Printf("server fail to listen")
		return
	}

	// 首先根据协程池的worker数量限定,开启固定数量的Worker,
	// 每一个Worker用一个Goroutine承载
	for i := 0; i < p.worker_num; i++ {
		go p.worker()
	}

	// 将新申请的连接加入到就绪队列
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		p.JobsChannel <- conn
	}
}

// Initialize servers, listening requests from clients.
func InitClient(c *calculator.Calculator) error {
	if c == nil {
		return fmt.Errorf("cannot create stack calculator")
	}
	stackCal = c

	c.NotiftyMap = make(map[int64]chan calculator.ApplyResp)

	// to get the server address
	connNet := viper.GetString("client.network")
	peers := viper.GetStringSlice("client.addr")
	address := peers[Me]

	// 创建一个协程池,最大开启3个协程worker
	p := NewPool(workerNum)

	// 设定监听，启动协程池p
	go p.Run(connNet, address)

	return nil
}

package calculator

import (
	"fmt"
	"server/models"
	"server/raft"
)

// Create stack calculator.
func Make(rf *raft.Raft, applyCh <-chan models.ApplyMsg) *Calculator {
	c := Calculator{
		InstanceMap: make(map[int64]*Instance, 0),
		stackID:     0,
		Raft:        rf,
		ApplyCh:     applyCh,
	}
	go c.listenAndApply(applyCh)
	return &c
}

// Return how many instances.
func (c *Calculator) GetNum() int64 {
	return c.stackID
}

// When Raft applys some logs, nofities manager by applyCh.
func (c *Calculator) listenAndApply(applyCh <-chan models.ApplyMsg) {
	for {
		req := <-applyCh
		if req.Command.Instruction == "" {
			fmt.Println("empty command")
			continue
		}

		result, err := c.Excution(req.Command)
		msg := fmt.Sprintf("result: %v, error: %v", result, err)

		// notify the manager that write request has excuted
		c.mu.Lock()
		ch, ok := c.NotiftyMap[int64(req.CommandIndex)]
		if ok {
			ch <- ApplyResp{
				Message:  msg,
				IsLeader: req.IsLeader,
				Success:  true,
			}
		}
		c.mu.Unlock()
	}
}

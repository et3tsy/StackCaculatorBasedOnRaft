package calculator

import (
	"fmt"
	"server/models"
)

// Handle with all requests from clients.
// Response to clients may have dalay.
func (c *Calculator) Manage(req models.Request) (resp models.Response) {
	// read-only request
	if req.Instruction == "get" || req.Instruction == "Get" {
		if c.Raft.CheckBeforeRead() {
			v, err := c.Get(req.Params[0])
			resp.Message = fmt.Sprintf("[%v]Reply: %v, err: %v", req.Params[0], v, err)
			resp.Success = true
		} else {
			resp.Message = "no leader"
			resp.Value = int(c.Raft.GetLeader())
		}
	} else {
		// write request, start agreement
		index, _, leader := c.Raft.Start(req)
		if !leader {
			resp.Message = "no leader"
			resp.Value = int(c.Raft.GetLeader())
			return
		}

		// register for notification, and notify the one
		// before listening to this index's channel to return
		ar := make(chan ApplyResp, 0)
		c.mu.Lock()
		ch, ok := c.NotiftyMap[int64(index)]
		if ok {
			ch <- ApplyResp{
				Message: "fail to write",
				Success: false,
			}
		}
		c.NotiftyMap[int64(index)] = ar
		c.mu.Unlock()

		// the to check term and index
		info := <-ar
		if info.Success && info.IsLeader {
			resp.Success = true
			resp.Message = info.Message
		} else {
			resp.Message = "fail to write"
		}
	}
	return
}

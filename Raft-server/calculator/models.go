package calculator

import (
	"server/models"
	"server/raft"
	"sync"
)

type Instance struct {
	data []int64
}

type Calculator struct {
	InstanceMap map[int64]*Instance
	stackID     int64
	Raft        *raft.Raft
	ApplyCh     <-chan models.ApplyMsg
	NotiftyMap  map[int64]chan ApplyResp
	mu          sync.Mutex
}

// After apply locally, response to the write request.
type ApplyResp struct {
	Message  string
	IsLeader bool
	Success  bool
}

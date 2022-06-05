package models

// As each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service
// on the same server, via the applyCh passed to Make().
// Set CommandValid to true to indicate that the ApplyMsg contains
// a newly committed log entry.
type ApplyMsg struct {
	IsLeader     bool
	CommandValid bool
	CommandIndex int
	Command      Request
}

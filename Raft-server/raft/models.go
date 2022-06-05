package raft

import (
	"server/models"
	"sync"
	"time"
)

// Implement a single entry.
type Entry struct {
	Index   int64
	Term    int64
	Command models.Request
}

// Implement a single Raft peer.
type Raft struct {
	mu      sync.Mutex           // lock to protect shared access between RPCs
	mu2     sync.Mutex           // lock when operating log
	peers   []string             // RPC end points of all peers
	me      int                  // this peer's index into peers[]
	leader  int64                // leader ID, renew when AppendEntries RPC
	dead    int32                // set by Kill()
	applyCh chan models.ApplyMsg // to pass messages that apply entries

	currentTerm int64       // lastest term server has seen
	actAs       int64       // ActAsLeader, ActAsFollower or ActAsCandidate
	clock       *time.Timer // to record timeout

	log           []Entry // the log entries
	commitIndex   int64   // highest log entry known to be commited
	lastApplied   int64   // the last index applied. valid entries start from 1
	voteFor       int64   // candidateId that received vote in current term
	readAvailable int64   // whether leader's data is lastest
	readCommit    int64   // once leader comes into power, commit an empty entry

	nextIndex  []int64 // index of the next log entry to send to i's peer
	matchIndex []int64 // index of highest log entry known to be replicated on server

	network string // the network RPC use
}

type RaftRPC Raft

type RequestVoteArgs struct {
	Term         int64
	CandidateId  int64
	LastLogIndex int64
	LastLogTerm  int64
}

type RequestVoteReply struct {
	Term        int64
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int64
	LeaderId     int64
	PrevLogIndex int64
	PrevLogTerm  int64
	LeaderCommit int64
	Entries      []Entry
}

type AppendEntriesReply struct {
	Term    int64
	Success bool
}

package raft

import (
	"fmt"
	"net/rpc"
	"sync/atomic"
)

// RequestVote RPC handler.
// Instead of Raft, RaftRPC is used to be exposed to rpc.Register.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	// avoid concurrency problems between RPCs
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// to maintain the CurrentTerm
	rf.updateCurrentTerm(args.Term)

	currentTerm := rf.getCurrentTerm()

	reply.VoteGranted = false
	reply.Term = currentTerm

	// Receiver implementation 1
	if args.Term < currentTerm {
		return nil
	}

	index, term := rf.getLastEntry()

	// Receiver implementation 2
	vf := atomic.LoadInt64(&rf.voteFor)
	if !(vf == -1 || vf == args.CandidateId) ||
		(term == args.LastLogTerm && index > args.LastLogIndex) ||
		term > args.LastLogTerm {
		return nil
	}

	// grant vote
	rf.clock.Reset(getTimeout()) // to reset the clock
	atomic.StoreInt64(&rf.voteFor, args.CandidateId)
	reply.VoteGranted = true
	return nil
}

// Call RequestVote RPC handler. If the server is
// down, may return error. Renew term when returning.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	c, err := rpc.Dial(rf.network, rf.peers[server])
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	defer c.Close()

	err = c.Call("Raft.RequestVote", args, reply)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}

	// to maintain the CurrentTerm
	rf.updateCurrentTerm(reply.Term)

	return true
}

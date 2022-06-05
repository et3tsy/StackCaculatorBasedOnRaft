package raft

import (
	"fmt"
	"net/rpc"
	"sync/atomic"
)

// RequestVote RPC handler.
// Instead of Raft, RaftRPC is used to
// be exposed to rpc.Register.
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

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
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

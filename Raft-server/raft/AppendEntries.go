package raft

import (
	"fmt"
	"net/rpc"
	"sync/atomic"
)

// Implement AppendEntries RPC handler
// Instead of Raft, RaftRPC is used to
// be exposed to rpc.Register.
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	// avoid concurrency problems between RPCs
	rf.mu.Lock()
	defer rf.mu.Unlock()

	for i := range args.Entries {
		fmt.Printf("[%v]Receive new entries: %v\n", rf.me, args.Entries[i].Index)
	}

	// to maintain the CurrentTerm
	rf.updateCurrentTerm(args.Term)
	rf.clock.Reset(getTimeout())

	// the rule 3 in figure 2 (Rules for Servers $5.2)
	if rf.getActAs() != ActAsFollower {
		fmt.Printf("[%v]Reject, prevIndex: %v, prevLog: %v\n", rf.me, args.PrevLogIndex, args.PrevLogTerm)
		rf.changeTo(ActAsFollower)
	}

	currentTerm := rf.getCurrentTerm()
	reply.Success = false
	reply.Term = currentTerm

	// Receiver implementation 1
	if args.Term < currentTerm {
		return nil
	}

	// Receiver implementation 2
	if rf.getEntry(args.PrevLogIndex) != args.PrevLogTerm {
		return nil
	}

	atomic.StoreInt64(&rf.leader, args.LeaderId)

	// append entries, and apply log[lastApplied] to state machine
	rf.appendEntries(args.Entries)

	// to renew commitIndex
	rf.renewCommitIndex(args.LeaderCommit)

	reply.Success = true
	return nil
}

// Call AppendEntries RPC handler.
func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	c, err := rpc.Dial(rf.network, rf.peers[server])
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	defer c.Close()

	err = c.Call("Raft.AppendEntries", args, reply)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}

	// to maintain the CurrentTerm
	rf.updateCurrentTerm(reply.Term)
	return true
}

// Append entries, and apply log[lastApplied] to state machine.
func (rf *Raft) appendEntries(entries []Entry) {
	// lock when operate log[]
	rf.mu2.Lock()
	defer rf.mu2.Unlock()

	// iterate the entries received, and check conflicts
	for i := range entries {
		if entries[i].Index > int64(len(rf.log)-1) {
			rf.log = append(rf.log, entries[i])
		} else {
			index := entries[i].Index
			if rf.log[index].Term != entries[i].Term {
				rf.log = rf.log[:index]
				rf.log = append(rf.log, entries[i])
			}
		}
	}
}

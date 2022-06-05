package raft

import (
	"fmt"
	"server/models"
	"sync/atomic"
)

func min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

// Apply an entry with the given index.
func (rf *Raft) applyEntry(index int64) {
	fmt.Printf("[%v]Apply index: %v\n", rf.me, index)
	// lock when operate log[]
	rf.mu2.Lock()
	defer rf.mu2.Unlock()

	// increment lastApplied
	atomic.AddInt64(&rf.lastApplied, 1)

	// apply the commands in log[index]
	rf.applyCh <- models.ApplyMsg{
		CommandValid: true,
		IsLeader:     rf.getActAs() == ActAsLeader,
		Command:      rf.log[index].Command,
		CommandIndex: int(index),
	}
	fmt.Printf("[%v]Apply index: %v Success!\n", rf.me, index)
}

// Change into new identity.
func (rf *Raft) changeTo(nextActAs int) {
	switch nextActAs {
	case ActAsCandidate:
		{
			atomic.StoreInt64(&rf.actAs, ActAsCandidate)
			rf.clock.Reset(getTimeout())
			go rf.startElection()
		}
	case ActAsLeader:
		{
			atomic.StoreInt64(&rf.actAs, ActAsLeader)
			rf.clock.Reset(0)
		}
	default:
		{
			atomic.StoreInt64(&rf.actAs, ActAsFollower)
			rf.clock.Reset(getTimeout())
		}
	}
}

// Update CurrentTerm.
func (rf *Raft) updateCurrentTerm(term int64) {
	// to maintain the CurrentTerm
	if term > rf.getCurrentTerm() {
		atomic.StoreInt64(&rf.currentTerm, term)
		atomic.StoreInt64(&rf.voteFor, -1)
		if rf.getActAs() != ActAsFollower {
			rf.changeTo(ActAsFollower)
		}
	}
}

// Renew commitIndex passed by AppendEntries RPC,
// and apply entries.
func (rf *Raft) renewCommitIndex(commited int64) {
	if commited > rf.getCommitIndex() {
		idx, _ := rf.getLastEntry()
		atomic.StoreInt64(&rf.commitIndex, min(commited, idx))
	}
}

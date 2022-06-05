package raft

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	electionTimeoutLower = 900 // *time.Millisecond
	electionTimeoutDelta = 400 // *time.Millisecond
	ActAsLeader          = 1
	ActAsFollower        = 2
	ActAsCandidate       = 3
)

// To get the role which the server is acting as.
func (rf *Raft) getActAs() int64 {
	return atomic.LoadInt64(&rf.actAs)
}

// To get the curTerm.
func (rf *Raft) getCurrentTerm() int64 {
	return atomic.LoadInt64(&rf.currentTerm)
}

// To get the last log entry.
func (rf *Raft) getLastEntry() (int64, int64) {
	// lock when operate log[]
	rf.mu2.Lock()
	defer rf.mu2.Unlock()

	idx := int64(len(rf.log) - 1)
	if idx == 0 {
		return 0, 0
	}
	return idx, atomic.LoadInt64(&rf.log[idx].Term)
}

// To get the term of the log entry by the given index.
func (rf *Raft) getEntry(idx int64) int64 {
	// lock when operate log[]
	rf.mu2.Lock()
	defer rf.mu2.Unlock()

	// be careful when idx == 0, its term should be 0
	if idx == 0 {
		return 0
	} else if idx < 0 || idx >= int64(len(rf.log)) {
		return -1
	}
	return atomic.LoadInt64(&rf.log[idx].Term)
}

// To get the commitIndex of the leader.
func (rf *Raft) getCommitIndex() int64 {
	return atomic.LoadInt64(&rf.commitIndex)
}

// To get the leader ID.
func (rf *Raft) GetLeader() int64 {
	return atomic.LoadInt64(&rf.leader)
}

// To get the lastApplied of the leader.
func (rf *Raft) getLastApplied() int64 {
	return atomic.LoadInt64(&rf.lastApplied)
}

// To get the nextIndex[i] from the leader.
func (rf *Raft) getNextIndex(i int) int64 {
	return atomic.LoadInt64(&rf.nextIndex[i])
}

// To get the matchIndex[i] from the leader.
func (rf *Raft) getMatchIndex(i int) int64 {
	return atomic.LoadInt64(&rf.matchIndex[i])
}

// To get readCommit.
func (rf *Raft) getReadCommit() int64 {
	return atomic.LoadInt64(&rf.readCommit)
}

// To fill the AppendEntriesArgs.
func (rf *Raft) getAppendEntriesArgs(index int64, e []Entry) AppendEntriesArgs {
	return AppendEntriesArgs{
		Term:         rf.getCurrentTerm(),
		LeaderId:     int64(rf.me),
		PrevLogIndex: index - 1,
		PrevLogTerm:  rf.getEntry(index - 1),
		LeaderCommit: rf.getCommitIndex(),
		Entries:      e,
	}
}

// To get random time.
func getTimeout() time.Duration {
	return time.Millisecond * time.Duration(rand.Intn(electionTimeoutDelta)+electionTimeoutLower)
}

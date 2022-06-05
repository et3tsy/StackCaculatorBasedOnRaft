package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"fmt"
	"server/models"
	"sync/atomic"
	"time"
)

const (
	checkReadTimeout = time.Millisecond * 800
	applyDelay       = time.Millisecond * 30
)

// Return currentTerm and whether this server believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	term := rf.getCurrentTerm()
	isleader := false
	if rf.getActAs() == ActAsLeader {
		isleader = true
	}
	return int(term), isleader
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// The first return value is the index that the command will appear at
// if it's ever committed. The second return value is the current
// term. The third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command models.Request) (index, term int, isLeader bool) {
	index = -1
	term = -1
	isLeader = false

	if rf.getActAs() != ActAsLeader {
		return
	}

	fmt.Printf("[%v]leader append new cmd from client\n", rf.me)

	rf.mu2.Lock()
	defer rf.mu2.Unlock()

	index = len(rf.log)
	term = int(rf.getCurrentTerm())
	isLeader = true
	rf.log = append(rf.log, Entry{
		Index:   int64(index),
		Term:    int64(term),
		Command: command,
	})
	return
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []string, network string, me int, applyCh chan models.ApplyMsg) *Raft {

	// Raft initialization
	rf := &Raft{
		peers:       peers,
		network:     network,
		me:          me,
		dead:        0,
		currentTerm: 0,
		commitIndex: 0,
		lastApplied: 0,
		actAs:       ActAsFollower,
		clock:       time.NewTimer(0),
		log:         make([]Entry, 1),
		nextIndex:   make([]int64, len(peers)),
		matchIndex:  make([]int64, len(peers)),
		applyCh:     applyCh,
	}

	// act as follower at first
	rf.changeTo(ActAsFollower)

	// start ticker goroutine to start elections
	go rf.ticker()

	// apply entries
	go func() {
		for {
			for rf.getCommitIndex() > rf.getLastApplied() {
				rf.applyEntry(rf.getLastApplied() + 1)
			}
			time.Sleep(applyDelay)
		}
	}()

	return rf
}

// The ticker go routine starts a new election if
// this peer hasn't received heartsbeats recently.
func (rf *Raft) ticker() {
	for {
		<-rf.clock.C
		if rf.getActAs() == ActAsLeader {
			continue
		} else {
			rf.changeTo(ActAsCandidate)
		}
	}
}

// To implement linearizable semantics, we should exchange heartbeat
// messages with a majority of the cluster before responding to
// read-only requests.
func (rf *Raft) CheckBeforeRead() bool {
	if rf.getActAs() != ActAsLeader {
		return false
	}

	count := int64(1)
	t := time.NewTimer(checkReadTimeout)
	for i := range rf.peers {
		if i != rf.me {
			go func(id int) {
				args := rf.getAppendEntriesArgs(rf.getNextIndex(id), []Entry{})
				reply := AppendEntriesReply{
					Term:    0,
					Success: false,
				}

				// send heartbeats
				if rf.sendAppendEntries(id, &args, &reply) {
					atomic.AddInt64(&count, 1)
					if atomic.LoadInt64(&count) > int64(len(rf.peers)/2) {
						t.Reset(0)
					}
				}
			}(i)
		}
	}
	<-t.C
	if atomic.LoadInt64(&count) > int64(len(rf.peers)/2) && rf.getActAs() == ActAsLeader {
		return true
	} else {
		return false
	}
}

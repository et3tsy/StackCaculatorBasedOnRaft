package raft

import (
	"sync/atomic"
	"time"
)

const (
	VoteAgree          = 1
	VoteDisagree       = 0
	TimeoutRequestVote = time.Millisecond * 800
)

// Change to candidate, and start an election.
func (rf *Raft) startElection() {
	// increment currentTerm
	voteTerm := rf.getCurrentTerm() + 1
	atomic.StoreInt64(&rf.currentTerm, voteTerm)
	atomic.StoreInt64(&rf.voteFor, int64(rf.me))

	// to record the vote from other Raft servers
	voteCount := int64(1)
	voteTarget := int64(len(rf.peers)/2 + 1)

	// wake up exact one goroutine
	flag := int64(0)

	for i := range rf.peers {
		if i == rf.me {
			continue
		}
		go func(id int) {
			for {
				// if the server doesn't act as candidate, then return
				if voteTerm < rf.getCurrentTerm() || rf.getActAs() != ActAsCandidate {
					return
				}

				// get the infomation for last entry
				index, term := rf.getLastEntry()

				args := RequestVoteArgs{
					Term:         voteTerm,
					CandidateId:  int64(rf.me),
					LastLogIndex: index,
					LastLogTerm:  term,
				}

				replys := RequestVoteReply{
					Term:        0,
					VoteGranted: false,
				}

				// call RequestVote RPC
				if rf.sendRequestVote(id, &args, &replys) {
					if replys.VoteGranted {
						// fmt.Printf("[%v]Granted by %v\n", rf.me, id)
						atomic.AddInt64(&voteCount, 1)
					}
				}

				// if the server is elected as leader, wake up a goroutine to do relate work
				if atomic.LoadInt64(&voteCount) >= voteTarget {
					if atomic.LoadInt64(&flag) == 0 {
						atomic.StoreInt64(&flag, 1)
						atomic.StoreInt64(&rf.readAvailable, 0)
						rf.changeTo(ActAsLeader)
						go rf.actAsLeader()
					}
					return
				}

				if replys.VoteGranted {
					return
				}

				time.Sleep(TimeoutRequestVote)
			}
		}(i)
	}
}

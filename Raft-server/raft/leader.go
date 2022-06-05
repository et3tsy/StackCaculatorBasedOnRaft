package raft

import (
	"fmt"
	"server/models"
	"sync/atomic"
	"time"
)

const (
	DelayHeartbeat = time.Millisecond * 100
	CommitDelay    = time.Millisecond * 100
	LeaderFailed   = time.Millisecond * 1000
)

// Act as Leader, and do related work.
func (rf *Raft) actAsLeader() {
	atomic.StoreInt64(&rf.leader, int64(rf.me))
	fmt.Printf("[%v]Be a leader, currentTerm: %v\n", rf.me, rf.getCurrentTerm())

	// reinitialize after election
	for i := 0; i < len(rf.peers); i++ {
		idx, _ := rf.getLastEntry()
		atomic.StoreInt64(&rf.nextIndex[i], idx+1)
		atomic.StoreInt64(&rf.matchIndex[i], 0)
	}

	t := time.NewTimer(LeaderFailed)

	// if leader fails, turn into follower
	go func() {
		<-t.C
		rf.changeTo(ActAsFollower)
	}()

	// renew commitIndex
	go func() {
		for rf.getActAs() == ActAsLeader {
			count := 1
			commited := rf.getCommitIndex()
			for i := range rf.peers {
				// to count how many peers have appended entries
				if commited+1 <= rf.getMatchIndex(i) {
					count++
				}

				// if majority have appended, then commit
				if count > len(rf.peers)/2 {
					atomic.AddInt64(&rf.commitIndex, 1)
					rf.renewCommitIndex(commited + 1)
					if commited+1 >= rf.getReadCommit() {
						atomic.StoreInt64(&rf.readAvailable, 1)
					}
					break
				}
			}
			time.Sleep(CommitDelay)
		}
	}()

	// upon election, send initial empty AppendEntries RPCs,
	// and sending heartbeats periodically
	for i := range rf.peers {
		if i != rf.me {
			go func(id int) {
				e := []Entry{}
				for rf.getActAs() == ActAsLeader {
					// to check if need to append entries
					rf.mu2.Lock()
					if len(e) == 0 && rf.getNextIndex(id) <= int64(len(rf.log)-1) {
						e = append(e, rf.log[rf.getNextIndex(id)])
					}
					rf.mu2.Unlock()

					args := rf.getAppendEntriesArgs(rf.getNextIndex(id), e)
					reply := AppendEntriesReply{
						Term:    0,
						Success: false,
					}

					// send AppendEntries RPC
					if rf.sendAppendEntries(id, &args, &reply) {
						t.Reset(LeaderFailed) // leader may fail. once failed, change to follower

						// if len(e) == 0, we actually send out heartbeats
						if reply.Success {
							atomic.AddInt64(&rf.nextIndex[id], int64(len(e)))
							atomic.StoreInt64(&rf.matchIndex[id], rf.getNextIndex(id)-1)
							rf.mu2.Lock()
							index := rf.getNextIndex(id)
							e = e[:0]
							for i := index; i <= int64(len(rf.log)-1); i++ {
								e = append(e, rf.log[i])
							}
							rf.mu2.Unlock()
						} else {
							atomic.AddInt64(&rf.nextIndex[id], -1)
							e = e[:0]
						}
					} else {
						e = e[:0]
					}

					// the delay of heartbeats
					time.Sleep(DelayHeartbeat)
				}
			}(i)
		}
	}

	// implement linearizable reads, once leader comes into power,
	// commit an empty entry.
	go func() {
		r, _, _ := rf.Start(models.Request{
			Instruction: "",
			Params:      []int64{},
		})
		atomic.StoreInt64(&rf.readCommit, int64(r))
	}()
}

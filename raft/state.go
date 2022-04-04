package raft

import (
	"sync"

	"github.com/justin0u0/raft/pb"
)

type RaftState uint32

const (
	Follower RaftState = iota
	Candidate
	Leader
)

func (s RaftState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unkwown"
	}
}

type raftState struct {
	// raft state
	state RaftState

	// persistent state on all servers

	currentTerm uint64
	votedFor    uint32
	logs        []*pb.Entry

	// volatile state on all servers

	commitIndex uint64
	lastApplied uint64

	// volatile state on leader

	nextIndex  map[uint32]int64
	matchIndex map[uint32]int64

	mu sync.Mutex
}

// getLastLog get last log id and last log term,
// returns 0, 0 if none
func (rs *raftState) getLastLog() (id, term uint64) {
	if len(rs.logs) == 0 {
		return 0, 0
	}

	log := rs.logs[len(rs.logs)-1]

	return log.GetId(), log.GetTerm()
}

func (rs *raftState) toFollower(term uint64) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.state = Follower

	if rs.currentTerm < term {
		rs.currentTerm = term
		rs.votedFor = 0
	}
}

func (rs *raftState) toCandidate() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.state = Candidate
}

func (rs *raftState) toLeader() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.state = Leader
}

func (rs *raftState) voteFor(id uint32, voteForSelf bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// if vote for self, increase current term
	if voteForSelf {
		rs.currentTerm++
	}

	rs.votedFor = id
}

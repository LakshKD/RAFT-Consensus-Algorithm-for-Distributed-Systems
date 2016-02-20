package main

import (
	"math"
	//"fmt"
)

type error struct {
	err string
}

type Log struct {
	logindex int
	term     int
	entry    []byte
}

var log [201]Log

type send struct {
	DestID int
	Event  interface{}
}

type LogStore struct {
	index int

	entry []byte
}

type Commit struct {
	index int
	data  []byte
	error string
}

type Alarm struct {
	t int
}

type AppendEv struct {
	data []byte
}
type TimeoutEv struct {
	Time int64
}

type AppendEntriesRequestEv struct {
	Term         int
	PrevLogIndex int
	PrevLogTerm  int
	CommitIndex  int
	LeaderId     int
	Entry        []byte
}

type AppendEntriesResponseEv struct {
	Id      int
	Term    int
	Success bool
}

type VoteRequestEv struct {
	Term         int
	LastLogIndex int
	LastLogTerm  int
	CandidateId  int
}

type VoteResponseEv struct {
	Id          int
	Term        int
	VoteGranted bool
}

type StateMachine struct {
	id           int   // server id
	peers        []int // other server ids
	term         int
	PrevLogIndex int
	PrevLogTerm  int
	CommitIndex  int
	LastLogIndex int
	LastLogTerm  int
	votedfor     int
	votetracker  [6]int //It tracks from which machine the candidate received the vote
	nextIndex    [6]int
	matchIndex   [6]int
	status       string
	votecounter  int
	Entry        []byte
	LeaderID     int
	count        int
	falsevote    int
}

//AppendEntriesRequest Functions

func (msg *AppendEntriesRequestEv) AppendEntriesRequest_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)

	if sm.term <= msg.Term {
		//sm.LeaderID=msg.LeaderId
		sm.term = msg.Term //Update follower's term and set it to leaders term
		sm.votedfor = 0
		if len(msg.Entry) == 0 {
			//Its a heartbeat message so reset election timer of the follower

			actions = append(actions, Alarm{10})
			return (actions)

		} else {

			if sm.LastLogIndex == msg.PrevLogIndex && sm.LastLogTerm == msg.PrevLogTerm {
				index := sm.LastLogIndex

				if len(log[index+1].entry) != 0 && log[index+1].term != 0 {

					for i := index + 1; len(log[i].entry) != 0 && log[i].term != 0 && i <= 200; i++ {
						log[i].logindex = 0
						log[i].term = 0

					}

				}

				sm.LastLogIndex++
				//log[sm.LastLogIndex].logindex={sm.LastLogIndex,msg.Term,msg.Entry}

				if msg.CommitIndex > sm.CommitIndex {

					sm.CommitIndex = int(math.Min(float64(msg.CommitIndex), float64(sm.LastLogIndex)))

				}
				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, true}})
				actions = append(actions, LogStore{sm.LastLogIndex, msg.Entry})
				return (actions)

				if msg.CommitIndex > sm.CommitIndex {

					sm.CommitIndex = int(math.Min(float64(msg.CommitIndex), float64(sm.LastLogIndex)))

				}

			} else if msg.PrevLogIndex != -1 && len(log[msg.PrevLogIndex].entry) == 0 { //no entry at PrevlogIndex(sent by the leader)
				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
				return (actions)

			} else if msg.PrevLogIndex == -1 {

				if len(log[0].entry) != 0 && log[0].term != 0 {

					for i := 0; len(log[i].entry) != 0 && log[i].term != 0 && i <= 200; i++ {
						log[i].logindex = 0
						log[i].term = 0
						//log[i].entry=nil

					}
					sm.LastLogIndex = 0
					sm.LastLogTerm = msg.Term
				}

				actions = append(actions, LogStore{0, msg.Entry})
				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, true}})
				return (actions)
			} else if sm.LastLogIndex == msg.PrevLogIndex && sm.LastLogTerm != sm.LastLogTerm {
				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
				return (actions)
			}

		}

	} else {

		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
		return (actions)
	}

	return (make([]interface{}, 1))

}

func (msg *AppendEntriesRequestEv) AppendEntriesRequest_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)

	if sm.term <= msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
		if len(msg.Entry) == 0 {

			actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, true}})
			return (actions)
		} else {

			actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
			return (actions)

		}
	} else {

		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
		return (actions)

	}
}

func (msg *AppendEntriesRequestEv) AppendEntriesRequest_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)

	if sm.term >= msg.Term {

		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})

	} else {
		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
	}
	return (actions)
}

//AppendEntriesResponse Functions

func (msg *AppendEntriesResponseEv) AppendEntriesReponse_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
	}

	actions = append(actions, error{"Error"})

	return (actions)

}

func (msg *AppendEntriesResponseEv) AppendEntriesReponse_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
	}

	actions = append(actions, error{"Error"})

	return (actions)

}

func (msg *AppendEntriesResponseEv) AppendEntriesReponse_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	if msg.Success == false {

		if sm.term < msg.Term {
			sm.term = msg.Term
			sm.votedfor = 0
			sm.status = "Follower"
		} else {

			sm.nextIndex[msg.Id]--
			j := sm.nextIndex[msg.Id]
			j = j - 1
			actions = append(actions, send{msg.Id, AppendEntriesRequestEv{sm.term, j, log[j].term, sm.CommitIndex, sm.id, log[j].entry}})
			return (actions)
		}

	} else if msg.Success == true {

		sm.matchIndex[msg.Id]++
		j := sm.nextIndex[msg.Id]
		sm.nextIndex[msg.Id]++
		for i := 1; i <= 5; i++ {
			if sm.matchIndex[i] >= j {
				sm.count++
			}
		}
		if sm.count >= 3 {

			sm.count = 0
			sm.CommitIndex = sm.LastLogIndex
			actions = append(actions, Commit{index: j, data: log[sm.LastLogIndex].entry})
			return (actions)
		}

	}
	return (actions)

}

//VoteRequest Functions

func (msg *VoteRequestEv) VoteRequest_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)
	if msg.Term < sm.term {

		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})

		return (actions)

	} else if (sm.votedfor == 0 || sm.votedfor == msg.CandidateId) && msg.LastLogTerm >= sm.LastLogTerm && msg.LastLogIndex >= sm.LastLogIndex {

		sm.votedfor = msg.CandidateId

		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: true}})
		return (actions)
	} else {
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})
		return (actions)
	}

}

func (msg *VoteRequestEv) VoteRequest_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)
	if msg.Term < sm.term {
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})
		return (actions)
	} else if (sm.votedfor == 0 || sm.votedfor == msg.CandidateId) && msg.LastLogIndex >= sm.LastLogIndex && msg.LastLogTerm >= sm.LastLogTerm {
		if sm.votedfor == 0 {
			sm.term = msg.Term
			sm.votedfor = msg.CandidateId

		}
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: true}})
		return (actions)

	} else {

		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})
		return (actions)
	}

}
func (msg *VoteRequestEv) VoteRequest_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
	}
	actions = append(actions, error{"Error"})

	return (actions)
}

//VoteResponse functions

func (msg *VoteResponseEv) VoteResponse_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
	}
	actions = append(actions, error{"Error"})

	return (actions)

}

func (msg *VoteResponseEv) VoteResponse_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)

	if msg.VoteGranted == true {

		sm.votecounter++
		sm.votetracker[msg.Id] = 1

		if sm.votecounter >= 3 {
			sm.status = "Leader"

			for i := 0; i <= 3; i++ {
				actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term, sm.PrevLogIndex, sm.PrevLogTerm, sm.CommitIndex, sm.id, sm.Entry}})

			}
			return (actions)
		}
	} else if msg.VoteGranted == false {

		if sm.term < msg.Term {
			sm.term = msg.Term
			sm.votedfor = 0
			sm.status = "Follower"
		}
		sm.falsevote++
		if sm.falsevote >= 3 {
			sm.votedfor = 0
			sm.status = "Follower"
		}

	}
	return (make([]interface{}, 1))

}

func (msg *VoteResponseEv) VoteResponse_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
	}
	actions = append(actions, error{"Error"})
	return (actions)

}

//Append functions

func (msg *AppendEv) Append_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)

	actions = append(actions, Commit{data: msg.data, error: "NOT a Leader"})

	return (actions)

}

func (msg *AppendEv) Append_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)

	actions = append(actions, Commit{data: msg.data, error: "NOT a Leader"})

	return (actions)

}

func (msg *AppendEv) Append_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	sm.LastLogIndex++
	sm.PrevLogIndex++
	sm.matchIndex[sm.id]++
	for i := 0; i <= 3; i++ {
		sm.nextIndex[sm.peers[i]] = sm.LastLogIndex
	}
	actions = append(actions, LogStore{sm.LastLogIndex, msg.data})
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term, sm.PrevLogIndex, log[sm.PrevLogIndex].term, sm.CommitIndex, sm.id, msg.data}})

	}
	return (actions)
}

//Timeout funtcions

func (msg *TimeoutEv) Timeout_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)
	s := [6]int{0, 0, 0, 0, 0, 0}
	sm.votetracker = s
	sm.status = "Candidate"
	sm.votedfor = 0
	sm.term++
	sm.votetracker[sm.id] = 1
	sm.votecounter++
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{sm.term, sm.LastLogIndex, sm.LastLogTerm, sm.id}})

	}

	return (actions)
}

func (msg *TimeoutEv) Timeout_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)
	s := [6]int{0, 0, 0, 0, 0, 0}
	sm.votetracker = s
	sm.votecounter = 0
	sm.status = "Candidate"
	sm.votedfor = 0
	sm.term++
	sm.votetracker[sm.id] = 1
	sm.votecounter++
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{sm.term, sm.LastLogIndex, sm.LastLogTerm, sm.id}})

	}
	return (actions)

}

func (msg *TimeoutEv) Timeout_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	//Send HeartBeat message to all the followers
	Entry := []byte{}
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term, sm.PrevLogIndex, sm.LastLogIndex, sm.CommitIndex, sm.id, Entry}})
	}

	return (actions)

}

func (sm *StateMachine) ProcessEvent(ev interface{}) []interface{} {

	switch ev.(type) {

	case AppendEntriesRequestEv:
		AEReqEv := ev.(AppendEntriesRequestEv)
		//fmt.Println(AEReqEv)
		if sm.status == "Follower" {
			return (AEReqEv.AppendEntriesRequest_Handler_F())
		} else if sm.status == "Candidate" {
			return (AEReqEv.AppendEntriesRequest_Handler_C())
		} else if sm.status == "Leader" {
			return (AEReqEv.AppendEntriesRequest_Handler_L())
		}

	case AppendEntriesResponseEv:
		AEResEv := ev.(AppendEntriesResponseEv)
		if sm.status == "Follower" {
			return (AEResEv.AppendEntriesReponse_Handler_F())
		} else if sm.status == "Candidate" {
			return (AEResEv.AppendEntriesReponse_Handler_C())
		} else if sm.status == "Leader" {
			return (AEResEv.AppendEntriesReponse_Handler_L())
		}
	case VoteRequestEv:
		VReqEv := ev.(VoteRequestEv)
		if sm.status == "Follower" {
			return (VReqEv.VoteRequest_Handler_F())
		} else if sm.status == "Candidate" {
			return (VReqEv.VoteRequest_Handler_C())
		} else if sm.status == "Leader" {
			return (VReqEv.VoteRequest_Handler_L())
		}

	case VoteResponseEv:
		VResEv := ev.(VoteResponseEv)
		if sm.status == "Follower" {
			return (VResEv.VoteResponse_Handler_F())
		} else if sm.status == "Candidate" {
			return (VResEv.VoteResponse_Handler_C())
		} else if sm.status == "Leader" {
			return (VResEv.VoteResponse_Handler_L())
		}

	case AppendEv:
		ApndEv := ev.(AppendEv)
		if sm.status == "Follower" {
			return (ApndEv.Append_Handler_F())
		} else if sm.status == "Candidate" {
			return (ApndEv.Append_Handler_C())
		} else if sm.status == "Leader" {
			return (ApndEv.Append_Handler_L())
		}

	case TimeoutEv:
		TimeEv := ev.(TimeoutEv)
		if sm.status == "Follower" {
			return (TimeEv.Timeout_Handler_F())
		} else if sm.status == "Candidate" {
			return (TimeEv.Timeout_Handler_C())
		} else if sm.status == "Leader" {
			return (TimeEv.Timeout_Handler_L())
		}

		//default: ("Unrecognized")
	}

	return (make([]interface{}, 1))

}

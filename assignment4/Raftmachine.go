package main

import (
	_"fmt"
	"math"
	"time"
	log2 "github.com/cs733-iitb/log"
	"math/rand"
)

var sm *StateMachine

var num_of_resp int
type Log struct {
	Logindex int
	Term     int
	Entry    []byte
}

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
	//CommitIndex int
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
	id           int    // server id
	peers        [4]int // other server ids
	term         int
	PrevLogIndex int
	PrevLogTerm  int
	CommitIndex  int
	LastLogIndex int
	LastLogTerm  int
	votedfor     int
	votetracker  [6]int //It tracks from which machine the candidate received the vote
	nextIndex    []int
	matchIndex   [6]int
	status       string
	votecounter  int
	Entry        []byte
	LeaderID     int
	count        int
	falsevote    int
	log          *log2.Log
	ElectionTimeout int
	HeartBeatTimeout int
}

//AppendEntriesRequest Functions
func (msg *AppendEntriesRequestEv) AppendEntriesRequest_Handler_F() []interface{} {
	var actions = make([]interface{}, 1)
	if sm.term <= msg.Term {

		sm.LeaderID = msg.LeaderId
		sm.term = msg.Term //Update follower's term and set it to leaders term
		sm.votedfor = 0
		if len(msg.Entry) == 0 {
			if sm.CommitIndex < msg.CommitIndex {
				sm.CommitIndex = msg.CommitIndex
				//fmt.Println("Commitindex",sm.id,sm.CommitIndex)
				d,_:=sm.log.Get(int64(sm.CommitIndex))
				//fmt.Println(d)
				command:= d.(Log)
				actions = append(actions, Commit{index: sm.CommitIndex, data:command.Entry })
			}

			actions = append(actions, Alarm{sm.ElectionTimeout})
			//time.Sleep(1 * time.Second)
			//return (actions)

		} else {
			ld,_:=sm.log.Get(int64(msg.PrevLogIndex))
			lm:=ld.(Log)
			sm.LastLogIndex = lm.Logindex
			sm.LastLogTerm = lm.Term

			if sm.LastLogIndex == msg.PrevLogIndex && sm.LastLogTerm == msg.PrevLogTerm {
				//index := sm.LastLogIndex

				sm.LastLogIndex++

				if msg.CommitIndex > sm.CommitIndex {

					sm.CommitIndex = int(math.Min(float64(msg.CommitIndex), float64(sm.LastLogIndex)))

				}
				actions = append(actions, LogStore{sm.LastLogIndex, msg.Entry})
				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, true}})
				return (actions)

			}else{ 
                 if sm.LastLogIndex == msg.PrevLogIndex && sm.LastLogTerm != msg.PrevLogIndex {

				actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
				return (actions)
			  }
           }
		}

	} else {

		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
		
		
	}
	

	return (actions)

}

func (msg *AppendEntriesRequestEv) AppendEntriesRequest_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)

	if sm.term <= msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
		actions = append(actions,Alarm{sm.ElectionTimeout})
		//Call Follower's func
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

	}else {
		actions = append(actions, send{msg.LeaderId, AppendEntriesResponseEv{sm.id, sm.term, false}})
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
		actions = append(actions,Alarm{sm.HeartBeatTimeout})
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

	actions = append(actions, "Error")

	return (actions)

}

func (msg *AppendEntriesResponseEv) AppendEntriesReponse_Handler_C() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
		actions = append(actions,Alarm{sm.ElectionTimeout})
	}

	//actions = append(actions, error{"Error"})

	return (actions)

}

func (msg *AppendEntriesResponseEv) AppendEntriesReponse_Handler_L() []interface{} {

    //fmt.Println("Append Entry Response from ",msg)
    var prev int
	var actions = make([]interface{}, 1)
	if msg.Success == false {
		//fmt.Println("Hello its false")
		if sm.term < msg.Term {
			sm.term = msg.Term
			sm.votedfor = 0
			sm.status = "Follower"
			actions = append(actions,Alarm{sm.ElectionTimeout})
		} else {

			if sm.nextIndex[msg.Id] > 0 {
				sm.nextIndex[msg.Id]--
			}
			j := sm.nextIndex[msg.Id]
			if j > 0 {
				prev = j - 1
			}
			//fmt.Println("Its j",j)
			ld,_:=sm.log.Get(int64(j))
			ld1:=ld.(Log)
			actions = append(actions, send{msg.Id, AppendEntriesRequestEv{sm.term, prev, ld1.Term, sm.CommitIndex, sm.id, ld1.Entry}})
			//time.Sleep(1 * time.Second)
			
		}

	} else if msg.Success == true {
		//fmt.Println("Hello its true")
	
		
		sm.matchIndex[msg.Id]++
		j := sm.nextIndex[msg.Id]
		sm.nextIndex[msg.Id]++
		sm.count=0;
		for i := 1; i <= 5; i++ {
			if sm.matchIndex[i] >= j {
				sm.count++
			}
		}
		//fmt.Println("Value of count",sm.count)
		if(num_of_resp==0){
			sm.count=0
		}
		
		if sm.count >=3 && num_of_resp!=0 {
            num_of_resp=0
			
			sm.count = 0
			sm.CommitIndex = sm.LastLogIndex
			ld,_:=sm.log.Get(int64(sm.LastLogIndex))
			//fmt.Println("append resp_L",ld)
			ld1:=ld.(Log)
			//fmt.Println("Last Log index",sm.LastLogIndex)
			actions = append(actions, Commit{index: sm.LastLogIndex, data: ld1.Entry})
		}
		
	}
	

	return (actions)

}

//VoteRequest Functions

func (msg *VoteRequestEv) VoteRequest_Handler_F() []interface{} {
    //fmt.Println("Inside Follower vote Request",sm.id)
	var actions = make([]interface{}, 1)
	if msg.Term < sm.term {

		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})

		return (actions)

	} else if (sm.votedfor == 0 || sm.votedfor == msg.CandidateId) && msg.LastLogTerm >= sm.LastLogTerm && msg.LastLogIndex >= sm.LastLogIndex {

		sm.votedfor = msg.CandidateId
		//fmt.Println(sm.id,"Votedfor",sm.votedfor)
		sm.term = msg.Term
		actions = append(actions,Alarm{sm.ElectionTimeout})
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: true}})
		return (actions)
	} else {
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})
		return (actions)
	}

}

func (msg *VoteRequestEv) VoteRequest_Handler_C() []interface{} {
    //fmt.Println("Inside Vote_request_candidate\n")
	var actions = make([]interface{}, 1)
	
		actions = append(actions, send{msg.CandidateId, VoteResponseEv{Id: sm.id, Term: sm.term, VoteGranted: false}})
		return (actions)
	

}
func (msg *VoteRequestEv) VoteRequest_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term
		sm.votedfor = 0
		sm.status = "Follower"
		actions = append(actions,Alarm{sm.ElectionTimeout})
		
	}
	//actions = append(actions, error{"Error"})

	return (actions)
}

//VoteResponse functions

func (msg *VoteResponseEv) VoteResponse_Handler_F() []interface{} {

	var actions = make([]interface{}, 1)
	if sm.term < msg.Term {
		sm.term = msg.Term

		sm.votedfor = 0
	}
	return (actions)

}

func (msg *VoteResponseEv) VoteResponse_Handler_C() []interface{} {
	//fmt.Println("VOTE RESPONSE FROM",msg.Id,"to",sm.id)
	var actions = make([]interface{}, 1)

	if msg.VoteGranted == true {
         //fmt.Println("True VOTE RESPONSE FROM",msg.Id,"to",sm.id)
		sm.votecounter++
		sm.votetracker[msg.Id] = 1

		if sm.votecounter >= 3 {
			//fmt.Println("Hello am the leader",sm.id)
			//reset nextindex and matchindex
			sm.status = "Leader"
			actions = append(actions,Alarm{sm.HeartBeatTimeout})
			sm.LeaderID = sm.id
			for i := 0; i <= 3; i++ {
				actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term, sm.PrevLogIndex, sm.PrevLogTerm, sm.CommitIndex, sm.id, sm.Entry}})

			}
			return (actions)
		}
	} else if msg.VoteGranted == false {
        //fmt.Println("False VOTE RESPONSE FROM",msg.Id,"to",sm.id)
		if sm.term < msg.Term {
			sm.term = msg.Term
		}
		sm.falsevote++
		if sm.falsevote >= 3 {
			sm.votedfor = 0
			sm.status = "Follower"
			actions = append(actions,Alarm{sm.ElectionTimeout})
             return(actions)
			//Alarm reset
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
		actions = append(actions,Alarm{sm.ElectionTimeout})
		//Alarm reset
	}
	//actions = append(actions, error{"Error"})
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
    var j int
    num_of_resp=4
	var actions = make([]interface{}, 1)
	//fmt.Println("Next index array",sm.nextIndex)
	//fmt.Println("Inside Append of Leader")
	sm.PrevLogIndex++
	//fmt.Println("Leader's Prevlogindex",sm.PrevLogIndex)
	sm.LastLogIndex++
	//fmt.Println("Leader's LastLogIndex",sm.LastLogIndex)
	sm.matchIndex[sm.id]++ //Incrementing own match index
	//sm.count++;
	actions = append(actions, LogStore{sm.LastLogIndex, msg.data})
	for i := 0; i <= 3; i++ {
          p:=sm.nextIndex[sm.peers[i]]
            //fmt.Println("nextindex:",p,"sm.LastLogIndex:",sm.LastLogIndex)
          j=p
         for j=p;j<sm.LastLogIndex;j++{
         	//fmt.Println("its j",j)
         	ld,_:=sm.log.Get(int64(j))
			ldc:=ld.(Log)
			m:=j
			if(m>0){ 
             	m=m-1
              }  
             ldd,_:=sm.log.Get(int64(m))
             ldp:=ldd.(Log)
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term,ldp.Logindex,ldp.Term, sm.CommitIndex, sm.id, ldc.Entry}})
         
       }
       //fmt.Println("value of j",j)
       lm,_:=sm.log.Get(int64(j-1))
	   lmc:=lm.(Log)
	   //fmt.Println(lmc)
      actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term,lmc.Logindex,lmc.Term, sm.CommitIndex, sm.id, msg.data}}) 

    }
    for i := 0; i <= 3; i++ {
		sm.nextIndex[sm.peers[i]] = sm.LastLogIndex
	}
	//fmt.Println(actions)
	return(actions)
	
}

//Timeout funtcions

func (msg *TimeoutEv) Timeout_Handler_F() []interface{} {
    //fmt.Println("Inside Timeout_Handler_F(Machine timing out:)",sm.id)
	var actions = make([]interface{}, 1)
	s := [6]int{0, 0, 0, 0, 0, 0}
	sm.votetracker = s
	sm.status = "Candidate"
	actions = append(actions,Alarm{sm.ElectionTimeout})
	sm.votedfor = sm.id
	sm.term++
	sm.votetracker[sm.id] = 1
	sm.votecounter++
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{sm.term, sm.LastLogIndex, sm.LastLogTerm, sm.id}})

	}

	return (actions)
}

func (msg *TimeoutEv) Timeout_Handler_C() []interface{} {
  //fmt.Println("Inside Timeout_Handler_C")
	var actions = make([]interface{}, 1)
	s := [6]int{0, 0, 0, 0, 0, 0}
	sm.votetracker = s
	sm.votecounter = 0
	sm.status = "Follower"///Change to follower mode and set  Alarm between election timeout and 2*election timeout
	sm.ElectionTimeout=randInt(sm.ElectionTimeout,2*sm.ElectionTimeout)
	actions = append(actions,Alarm{sm.ElectionTimeout})
	sm.votedfor = 0
	sm.term++
	sm.votetracker[sm.id] = 1
	sm.votecounter++
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{sm.term, sm.LastLogIndex, sm.LastLogTerm, sm.id}})

	}
	return (actions)

}
func randInt(min int, max int) int {
	//fmt.Println("Hello")
    rand.Seed(time.Now().UTC().UnixNano())
    return min + rand.Intn(max-min)
}


func (msg *TimeoutEv) Timeout_Handler_L() []interface{} {

	var actions = make([]interface{}, 1)
	//When timeout occurs at Leader it will Send HeartBeat message to all the followers
	Entry := []byte{}
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{sm.term, sm.PrevLogIndex, sm.LastLogIndex, sm.CommitIndex, sm.id, Entry}})
	}

	return (actions)

}

func (SM *StateMachine) ProcessEvent(ev interface{}) []interface{} {
	sm = SM

	switch ev.(type) {

	case AppendEntriesRequestEv:
		AEReqEv := ev.(AppendEntriesRequestEv)
		if SM.status == "Follower" {
			
			return (AEReqEv.AppendEntriesRequest_Handler_F())
		} else if SM.status == "Candidate" {
			return (AEReqEv.AppendEntriesRequest_Handler_C())
		} else if SM.status == "Leader" {
			return (AEReqEv.AppendEntriesRequest_Handler_L())
		}

	case AppendEntriesResponseEv:
		AEResEv := ev.(AppendEntriesResponseEv)
		if SM.status == "Follower" {
			return (AEResEv.AppendEntriesReponse_Handler_F())
		} else if SM.status == "Candidate" {
			return (AEResEv.AppendEntriesReponse_Handler_C())
		} else if SM.status == "Leader" {
			return (AEResEv.AppendEntriesReponse_Handler_L())
		}
	case VoteRequestEv:
		VReqEv := ev.(VoteRequestEv)
		if SM.status == "Follower" {
			return (VReqEv.VoteRequest_Handler_F())
		} else if SM.status == "Candidate" {
			return (VReqEv.VoteRequest_Handler_C())
		} else if SM.status == "Leader" {
			return (VReqEv.VoteRequest_Handler_L())
		}

	case VoteResponseEv:
		VResEv := ev.(VoteResponseEv)
		if SM.status == "Follower" {
			return (VResEv.VoteResponse_Handler_F())
		} else if SM.status == "Candidate" {
			return (VResEv.VoteResponse_Handler_C())
		} else if SM.status == "Leader" {
			return (VResEv.VoteResponse_Handler_L())
		}

	case AppendEv:
		ApndEv := ev.(AppendEv)
		if SM.status == "Follower" {
			return (ApndEv.Append_Handler_F())
		} else if SM.status == "Candidate" {
			return (ApndEv.Append_Handler_C())
		} else if SM.status == "Leader" {
			
			return (ApndEv.Append_Handler_L())
		}

	case TimeoutEv:
		TimeEv := ev.(TimeoutEv)
		if SM.status == "Follower" {
              return (TimeEv.Timeout_Handler_F())
		} else if SM.status == "Candidate" {
			return (TimeEv.Timeout_Handler_C())
		} else if SM.status == "Leader" {
			return (TimeEv.Timeout_Handler_L())
		}

		
	}

	return (make([]interface{}, 1))

}

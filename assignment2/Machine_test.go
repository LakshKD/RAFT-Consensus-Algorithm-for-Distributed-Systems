package main

import (
	//"bufio"
	//"log"
	//"net/rpc"
	//"os"
	"fmt"
	"reflect"
	"testing"
)

/*
type Actions struct{

    Name string
    var m map[string]Value

}*/

var sm *StateMachine

//log *Log
//fmt.Println("*****************************************Follower's Testing Zone**********************************")

func TestFollowerAppendEntryreq(t *testing.T) {

	fmt.Println("-----------------------------------TestFollowerAppendEntryreq-------------------------------------------")
	/*id 			int // server id
		 peers 			[]int // other server ids
		 term 			int
		 PrevLogIndex  	uint64
	     PrevLogTerm   	uint64
	     CommitIndex   	uint64
	     LastLogIndex  	uint64
		 LastLogTerm   	uint64
	     votedfor      	[]int
	     votecounter   	[]int
	     nextIndex      []int
	     matchIndex     []int
	     status         string*/

	//(Test1) Empty entries array heartbeat message(return Timeout t)

	peers := []int{2, 3, 4, 5}
	Entry := []byte{}

	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions1 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry})

	//fmt.Println((actions1[1:]))
	//fmt.Println()
	var actions = make([]interface{}, 1)
	actions = append(actions, Alarm{t: 10}) //expected
	expect1(t, actions[1:], actions1[1:])
	//Alarm{t:10}

	//(Test2) Everything is matching between statemachine and Leader(return Log store)

	Entry1 := []byte{123}
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions2 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry1})
	//fmt.Println(actions2[1:])

	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, true}}) //expected
	actions = append(actions, LogStore{2, Entry1})                          //expected
	expect1(t, actions[1:], actions2[1:])

	//fmt.Println()
	//(Test3)the term of the append entries request message is less than the term of the statemachine(return Append entries response with false)
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 0, LastLogTerm: 1, status: "Follower", CommitIndex: 0}
	actions3 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 0, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry1})

	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 1, false}})
	expect1(t, actions[1:], actions3[1:])
	//fmt.Println(actions3[1:])
	//fmt.Println()
	//(Test4) Its the first message send by the Leader so PrevLogIndex is -1(return LogStore)
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 0, LastLogTerm: 1, status: "Follower"}
	actions4 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: -1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry1})

	actions = make([]interface{}, 1)
	actions = append(actions, LogStore{0, Entry1})
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, true}}) //expected actions
	expect1(t, actions[1:], actions4[1:])
	//fmt.Println(actions4[1:])
	//fmt.Println()
	//(Test5) LastLogIndex of stateMachine matches with the PrevLogIndex but the LastLogTerm is not matching with the Previous Log Term(Return Append entries response with false)
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions5 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 2, CommitIndex: 4, LeaderId: 4, Entry: Entry1})

	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, false}})
	expect1(t, actions[1:], actions5[1:])
	//fmt.Println(actions5[1:])
	//fmt.Println()

}

func TestFollowerAppendEntryResp(t *testing.T) {

	fmt.Println("----------------------------TestFollowerAppendEntryResp-----------------------------------------------------")

	peers := []int{2, 3, 4, 5}
	//(Test1) statemachine's term is less than the term in AppendEntriesResponseEv so return event is an error and term of SM gets changed
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions1 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 2, Term: 2, Success: false})
	actions := make([]interface{}, 1)
	actions = append(actions, error{"Error"}) //expected output actions
	expect2(t, actions[1:], actions1[1:], 2, sm.term)
	//expect2(t,2,sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2) statemachine's term is more than the term in AppendEntriesResponseEv so return event is an error and term of SM remains same
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions2 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 2, Term: 1, Success: true})
	actions = make([]interface{}, 1)
	actions = append(actions, error{"Error"}) //expected output actions
	expect2(t, actions[1:], actions2[1:], 2, sm.term)
	//expect2(t,1,sm.term)

	//fmt.Println(sm.term)
	//fmt.Println(actions2[1:])

}

func TestFollowerVoteReq(t *testing.T) {

	fmt.Println("------------------------------------TestFollowerVoteReq-------------------------------------------------")
	/*id 			int // server id
		 peers 			[]int // other server ids
		 term 			int
		 PrevLogIndex  	uint64
	     PrevLogTerm   	uint64
	     CommitIndex   	uint64
	     LastLogIndex  	uint64
		 LastLogTerm   	uint64
	     votedfor      	[]int
	     votecounter   	[]int
	     nextIndex      []int
	     matchIndex     []int
	     status         string*/

	peers := []int{2, 3, 4, 5}
	//(Test1) StateMachine's term is > than candidate term so return is voteresponse with false
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Follower"}
	actions1 := sm.ProcessEvent(VoteRequestEv{Term: 1, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions := make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 2, false}}) //expected
	expect1(t, actions[1:], actions1[1:])
	//fmt.Println(actions1[1:])
	//fmt.Println()

	//(Test2) follower will vote for the candidate and return event is vote reponse with true

	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Follower"}
	actions2 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 2, true}}) //expected

	expect1(t, actions[1:], actions2[1:])
	//fmt.Println(actions2[1:])

	//fmt.Println()

	//(Test3) LastLogindex of candidate is less so return event is vote response with false

	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 2, LastLogTerm: 1, votedfor: 0, status: "Follower"}
	actions3 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, false}}) //expected

	expect1(t, actions[1:], actions3[1:])
	//fmt.Println(actions3[1:])
	//fmt.Println()
	//(Test4) Last log index is same but the term of last log entry is less than the follower's entry so the return event is vote response with false
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 3, votedfor: 0, status: "Follower"}
	actions4 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: -1, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, false}}) //expected
	expect1(t, actions[1:], actions4[1:])
	//fmt.Println(actions4[1:])
	//fmt.Println()
	//(Test5)  Follower already voted for candidate 3 so the return event will be the vote response with false
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, votedfor: 3, status: "Follower"}
	actions5 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 1, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, false}}) //expected
	expect1(t, actions[1:], actions5[1:])

	//fmt.Println(actions5[1:])
	//fmt.Println()

}

func TestFollowerVoteResp(t *testing.T) {

	fmt.Println("----------------------------TestFollowerVoteResp-----------------------------------------------------")
	peers := []int{2, 3, 4, 5}
	actions := make([]interface{}, 1)
	//(Test1) statemachine's term is less than the term in  VoteResponseEv so return event is an error and term of SM gets changed
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions1 := sm.ProcessEvent(VoteResponseEv{Term: 2, VoteGranted: false})
	actions = append(actions, error{"Error"}) //expected output actions
	expect2(t, actions[1:], actions1[1:], 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2) statemachine's term is more than the term in VoteResponseEv so return event is an error and term of SM remains same
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions2 := sm.ProcessEvent(VoteResponseEv{Term: 1, VoteGranted: true})
	actions = make([]interface{}, 1)
	actions = append(actions, error{"Error"}) //expected output actions
	expect2(t, actions[1:], actions2[1:], 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(actions2[1:])

}
func TestFollowerAppendReq(t *testing.T) {

	peers := []int{2, 3, 4, 5}
	actions := make([]interface{}, 1)
	fmt.Println("----------------------------TestFollowerAppendReq-----------------------------------------------------")
	//Test1 In this case the state machine will simply return an error

	Entry1 := []byte{123}
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower"}
	actions1 := sm.ProcessEvent(AppendEv{data: Entry1})
	actions = append(actions, Commit{data: Entry1, error: "NOT a Leader"}) //expected output
	expect1(t, actions[1:], actions1[1:])
	//fmt.Println(actions1[1:])

}
func TestFollowerTimeout(t *testing.T) {

	fmt.Println("--------------------------------TestFollowerTimeout-----------------------------------------------------")

	peers := []int{2, 3, 4, 5}
	Entry1 := [6]int{0, 0, 0, 0, 0, 0}
	v := [6]int{0, 1, 0, 0, 0, 0}
	//Test1  send timeout event to follower so its term should increment and its status should also get changed
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Follower", votecounter: 0, votetracker: Entry1}
	actions := make([]interface{}, 1)
	actions1 := sm.ProcessEvent(TimeoutEv{Time: 10})
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{2, 1, 1, 1}})

	}
	expect3(t, 2, sm.term, "Candidate", sm.status, 1, sm.votecounter, v, sm.votetracker, actions, actions1)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(sm.votecounter)
	//fmt.Println(sm.votetracker)

	//fmt.Println(actions1[1:])
}

//fmt.Println("*****************************************Candidate's Testing Zone*********************************************")

func TestCandidateAppendEntryreq(t *testing.T) {

	fmt.Println("-----------------------------------TestCandidateAppendEntryreq-------------------------------------------")
	/*id 			int // server id
		 peers 			[]int // other server ids
		 term 			int
		 PrevLogIndex  	uint64
	     PrevLogTerm   	uint64
	     CommitIndex   	uint64
	     LastLogIndex  	uint64
		 LastLogTerm   	uint64
	     votedfor      	[]int
	     votecounter   	[]int
	     nextIndex      []int
	     matchIndex     []int
	     status         string*/

	actions := make([]interface{}, 1)
	peers := []int{2, 3, 4, 5}
	Entry := []byte{}
	//(Test1) Term of the message is less than the state Machine's Term so the return is an append entry response  with a false
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate"}
	actions1 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 1, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry})
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, false}}) //expected

	expect1(t, actions[1:], actions1[1:])
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2) Term of the message is >= the StateMachine's Term and its a heartbeat message so response is state gets changes to follower and Append entry response with true
	Entry1 := []byte{}
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate"}
	actions2 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry1})
	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, true}}) //expected

	expect4(t, actions[1:], actions2[1:], "Follower", sm.status)
	//fmt.Println(sm.status)
	//fmt.Println(actions2[1:])

	//fmt.Println()
	//(Test3)Term of the message is >= sm.term and its not heartbeat message so the response is append entry respnse with false
	Entry1 = []byte{123}
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 0, LastLogTerm: 1, status: "Candidate", CommitIndex: 0}
	actions3 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry1})
	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, false}}) //expected

	expect4(t, actions[1:], actions3[1:], "Follower", sm.status)

	//fmt.Println(sm.status)
	//fmt.Println(actions3[1:])

}

func TestCandidateAppendEntryResp(t *testing.T) {

	fmt.Println("----------------------------TestCandidateAppendEntryResp-----------------------------------------------------")
	peers := []int{2, 3, 4, 5}
	actions := make([]interface{}, 1)
	//(Test1) statemachine's term is less than the term in AppendEntriesResponseEv so return event is an error and term of SM gets changed
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate"}
	actions1 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 4, Term: 2, Success: false})
	actions = append(actions, error{"Error"}) //expected output actions
	expect2(t, actions[1:], actions1[1:], 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2) statemachine's term is more than the term in AppendEntriesResponseEv so return event is an error and term of SM remains same
	sm = &StateMachine{id: 1, peers: peers, term: 3, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate"}
	actions2 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 3, Term: 1, Success: true})
	//actions=append(actions,error{"Error"})//expected output actions
	expect2(t, actions[1:], actions2[1:], 3, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(actions2[1:])

}

func TestCandidateVoteReq(t *testing.T) {

	fmt.Println("------------------------------------TestCandidateVoteReq-------------------------------------------------")

	peers := []int{2, 3, 4, 5}
	actions := make([]interface{}, 1)
	//(Test1) Messgae term is < than StateMachine's Term so return event is Voteresponse with false
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Candidate"}
	actions1 := sm.ProcessEvent(VoteRequestEv{Term: 1, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions = append(actions, send{2, VoteResponseEv{1, 2, false}})
	expect1(t, actions[1:], actions1[1:])

	//fmt.Println(actions1[1:])
	//fmt.Println()

	//(Test2) Statemachine's term is <= Message term and msg.LastLogindex and msg.LastLogTerm is as updated as that of the Candidate
	//Expected Return Event is VoteResonse with votegranted=true
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Candidate"}
	actions2 := sm.ProcessEvent(VoteRequestEv{Term: 3, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 3, true}}) //expected
	expect5(t, actions[1:], actions2[1:], 3, sm.term, 2, sm.votedfor)
	//fmt.Println(sm.term)
	//fmt.Println(sm.votedfor)

	//fmt.Println(actions2[1:])
	//fmt.Println()

	//(Test3) StateMachine's term <=msg.term and LastLogindex Matches but LastLog Term not matches so
	//Return Event in this case VoteresponseEv with votegranted=false

	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 3, LastLogTerm: 3, votedfor: 0, status: "Candidate"}
	actions3 := sm.ProcessEvent(VoteRequestEv{Term: 3, LastLogIndex: 3, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 2, false}})
	expect1(t, actions[1:], actions3[1:])

	//fmt.Println(actions3[1:])
	//fmt.Println()
	//(Test4) StateMachine's term <=msg.term and LastLogindex does not matches
	//Return Event in this case VoteresponseEv with votegranted=false
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 3, LastLogTerm: 3, votedfor: 0, status: "Candidate"}
	actions4 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 2, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, false}})
	expect1(t, actions[1:], actions4[1:])

	//fmt.Println(actions4[1:])
	//fmt.Println()

	//(Test5)  Voted for candidate id but if the again request comes(may be becoz the Candidate does'nt get the message)
	//Return Event is VoteResponse with votegranted=true

	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, votedfor: 2, status: "Candidate"}
	actions5 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 2, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, true}})
	expect1(t, actions[1:], actions5[1:])
	//fmt.Println(actions5[1:])
	//fmt.Println()

	//(Test6) Voted for some other candidate id so now if some other candidate will try voterequest then it will
	//Return Vote Response with votegranted=false
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, votedfor: 3, status: "Candidate"}
	actions6 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 2, LastLogTerm: 2, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, send{2, VoteResponseEv{1, 1, false}})
	expect1(t, actions[1:], actions6[1:])
	//fmt.Println(actions6[1:])
	//fmt.Println()

}

func TestCandidateVoteResp(t *testing.T) {

	fmt.Println("----------------------------TestCandidateVoteResp-----------------------------------------------------")
	peers := []int{2, 3, 4, 5}
	a := [6]int{0, 1, 0, 0, 0, 0}
	v := [6]int{0, 1, 1, 0, 0, 0}
	actions := make([]interface{}, 1)
	//sm.votetracker=a //vote for self
	//sm.votecounter=1
	//(Test1) VoteGranted=true  from ID 2
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate", votetracker: a, votecounter: 1}
	sm.ProcessEvent(VoteResponseEv{Id: 2, Term: 1, VoteGranted: true})

	expect6(t, "Candidate", sm.status, 2, sm.votecounter, v, sm.votetracker)
	//fmt.Println(sm.status)
	//fmt.Println(sm.votecounter)
	//fmt.Println(sm.votetracker)
	//fmt.Println(actions1[1:])
	//fmt.Println()

	//(Test2) VoteGranted=true   from ID 3 at this point the candidate should become the leader
	//sm  = &StateMachine{id: 1,peers:peers,term:2,LastLogIndex:1,LastLogTerm:1,status:"Candidate"}
	sm.Entry = []byte{}
	actions2 := sm.ProcessEvent(VoteResponseEv{Id: 4, Term: 1, VoteGranted: true})
	actions = make([]interface{}, 1)
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{2, 0, 0, 0, 1, []byte{}}})

	}
	v1 := [6]int{0, 1, 1, 0, 1, 0}

	expect7(t, "Leader", sm.status, 3, sm.votecounter, v1, sm.votetracker, actions, actions2)
	//fmt.Println(sm.status)
	//fmt.Println(sm.votecounter)
	//fmt.Println(sm.votetracker)
	//fmt.Println(actions2[1:])
	//fmt.Println()

	//(Test3) VoteGranted=false   from ID 3
	//Term of the statemachine is less than the term of the voteresponse from ID 3 so state gets changed to follower and votedfor=0
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate", votetracker: a, votecounter: 1, votedfor: 2}
	sm.ProcessEvent(VoteResponseEv{Id: 3, Term: 2, VoteGranted: false})

	expect8(t, "Follower", sm.status, 0, sm.votedfor)

	//fmt.Println(sm.status)
	//fmt.Println(sm.votedfor)
	//fmt.Println()

	//(Test4) VoteGranted=false   from ID 3
	//Falsevote counter gets incremented to 1
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, falsevote: 0, LastLogTerm: 1, status: "Candidate", votetracker: a, votecounter: 1, votedfor: 2}

	sm.ProcessEvent(VoteResponseEv{Id: 3, Term: 1, VoteGranted: false})

	expect9(t, "Candidate", sm.status, 1, sm.falsevote, 2, sm.votedfor)
	//fmt.Println(sm.status)
	//fmt.Println(sm.falsevote)
	//fmt.Println(sm.votedfor)
	//fmt.Println()

	//falsevote incremented to 2
	sm.ProcessEvent(VoteResponseEv{Id: 4, Term: 1, VoteGranted: false})
	expect9(t, "Candidate", sm.status, 2, sm.falsevote, 2, sm.votedfor)
	//fmt.Println(sm.status)
	//fmt.Println(sm.falsevote)
	//fmt.Println(sm.votedfor)
	//fmt.Println()
	//falsevote incremented to 3 so majority gave false vote , state changed to candidate
	sm.ProcessEvent(VoteResponseEv{Id: 5, Term: 1, VoteGranted: false})
	expect9(t, "Follower", sm.status, 3, sm.falsevote, 0, sm.votedfor)
	//fmt.Println(sm.status)
	//fmt.Println(sm.falsevote)
	//fmt.Println(sm.votedfor)
	//fmt.Println()

}

func TestCandidateAppendReq(t *testing.T) {

	peers := []int{2, 3, 4, 5}
	fmt.Println("----------------------------TestCandidateAppendReq-----------------------------------------------------")
	//Test1 In this case the state machine will simply return an error

	Entry1 := []byte{225}
	actions := make([]interface{}, 1)
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate"}
	actions1 := sm.ProcessEvent(AppendEv{data: Entry1})

	actions = append(actions, Commit{0, Entry1, "NOT a Leader"})
	expect1(t, actions, actions1)

	//fmt.Println(actions1[1:])

}
func TestCandidateTimeout(t *testing.T) {

	fmt.Println("--------------------------------TestCandidateTimeout-----------------------------------------------------")

	peers := []int{1, 3, 4, 5}
	Entry1 := [6]int{0, 1, 0, 1, 1, 0} //Initial state of the voteTracker array
	v := [6]int{0, 0, 1, 0, 0, 0}
	actions := make([]interface{}, 1)
	//Test1  send timeout event to follower so its term should increment and its status should also get changed
	sm = &StateMachine{id: 2, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Candidate", votecounter: 4, votetracker: Entry1}

	actions1 := sm.ProcessEvent(TimeoutEv{Time: 10})
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], VoteRequestEv{3, 1, 1, 2}})

	}
	expect10(t, 3, sm.term, "Candidate", sm.status, 1, sm.votecounter, v, sm.votetracker, actions, actions1)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(sm.votecounter)
	//fmt.Println(sm.votetracker)
	//fmt.Println(actions1[1:])
}

//******************************************************Leader's Testing Zone****************************************************

func TestLeaderAppendEntryreq(t *testing.T) {

	fmt.Println("-----------------------------------TestLeaderAppendEntryreq-------------------------------------------")

	actions := make([]interface{}, 1)
	peers := []int{2, 3, 4, 5}
	Entry := []byte{}

	//(Test1)The Leader will simply rply with AppendEntry Response having false and status remains Leader
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Leader"}
	actions1 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 1, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry})
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 2, false}})

	expect4(t, actions, actions1, "Leader", sm.status)
	//fmt.Println(sm.status)
	//fmt.Println(actions1[1:])

	//(Test2)The Leader will simply rply with AppendEntry Response having false and the state of the leader becomes follower
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Leader"}
	actions2 := sm.ProcessEvent(AppendEntriesRequestEv{Term: 2, PrevLogIndex: 1, PrevLogTerm: 1, CommitIndex: 4, LeaderId: 4, Entry: Entry})
	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesResponseEv{1, 1, false}})
	expect4(t, actions, actions2, "Follower", sm.status)
	//fmt.Println(sm.status)
	//fmt.Println(actions2[1:])

}

func TestLeaderAppendEntryResp(t *testing.T) {

	fmt.Println("----------------------------TestLeaderAppendEntryResp-----------------------------------------------------")
	peers := []int{2, 3, 4, 5}

	nextindex := [6]int{0, 0, 5, 5, 5, 5}

	matchindex := [6]int{0, 5, 4, 4, 3, 2}
	/*
	   type Log struct{

	   logindex int
	   term     int
	   entry    []byte

	   }
	*/entry1 := []byte{12}
	entry2 := []byte{13}
	entry3 := []byte{14}
	entry4 := []byte{15}
	entry5 := []byte{16}
	entry6 := []byte{17}

	log[0] = Log{0, 1, entry1}
	log[1] = Log{1, 1, entry2}
	log[2] = Log{2, 1, entry3}
	log[3] = Log{3, 2, entry4}
	log[4] = Log{4, 2, entry5}
	log[5] = Log{5, 2, entry6}
	nxtindx := [6]int{0, 0, 6, 5, 5, 5}
	mindex := [6]int{0, 5, 5, 4, 3, 2}
	//(Test1) id of the leader is 1 and nextindex to send to peers is 4 and 4 is replicated on majority so entry will be committed
	sm = &StateMachine{id: 1, peers: peers, term: 2, count: 0, PrevLogIndex: 4, LastLogIndex: 5, LastLogTerm: 2, CommitIndex: 4, nextIndex: nextindex, matchIndex: matchindex, status: "Leader"}
	actions1 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 2, Term: 2, Success: true})
	actions := make([]interface{}, 1)
	expect11(t, 2, sm.count, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 4, sm.CommitIndex, actions, actions1)
	//fmt.Println(sm.count)
	//fmt.Println(sm.nextIndex)
	//fmt.Println(sm.matchIndex)
	//fmt.Println(sm.CommitIndex)
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2)
	nxtindx = [6]int{0, 0, 6, 6, 5, 5}
	mindex = [6]int{0, 5, 5, 5, 3, 2}
	actions2 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 3, Term: 1, Success: true})
	actions = append(actions, Commit{index: 5, data: entry6})
	expect11(t, 0, sm.count, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 5, sm.CommitIndex, actions, actions2)
	//fmt.Println(sm.count)
	//fmt.Println(sm.nextIndex)
	//fmt.Println(sm.matchIndex)
	//fmt.Println(sm.CommitIndex)
	//fmt.Println(actions2[1:])

	//Test3
	nxtindx = [6]int{0, 0, 6, 6, 4, 5}
	mindex = [6]int{0, 5, 5, 5, 3, 2}
	actions3 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 4, Term: 1, Success: false})
	actions = make([]interface{}, 1)
	actions = append(actions, send{4, AppendEntriesRequestEv{2, 3, 2, 5, 1, entry4}})

	expect12(t, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 5, sm.CommitIndex, actions, actions3)

	//fmt.Println(sm.nextIndex)
	//fmt.Println(sm.matchIndex)
	//fmt.Println(sm.CommitIndex)
	//fmt.Println(actions3[1:])

	//Test4

	nxtindx = [6]int{0, 0, 6, 6, 4, 4}
	mindex = [6]int{0, 5, 5, 5, 3, 2}
	actions = make([]interface{}, 1)
	actions4 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 5, Term: 1, Success: false})
	actions = append(actions, send{5, AppendEntriesRequestEv{2, 3, 2, 5, 1, entry4}})

	expect12(t, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 5, sm.CommitIndex, actions, actions4)
	//fmt.Println(sm.nextIndex)
	//fmt.Println(sm.matchIndex)
	//fmt.Println(sm.CommitIndex)
	//fmt.Println(actions4[1:])

	//Test5
	nxtindx = [6]int{0, 0, 6, 6, 4, 3}
	mindex = [6]int{0, 5, 5, 5, 3, 2}
	actions = make([]interface{}, 1)
	actions5 := sm.ProcessEvent(AppendEntriesResponseEv{Id: 5, Term: 1, Success: false})
	actions = append(actions, send{5, AppendEntriesRequestEv{2, 2, 1, 5, 1, entry3}})
	expect12(t, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 5, sm.CommitIndex, actions, actions5)
	//fmt.Println(sm.nextIndex)
	//fmt.Println(sm.matchIndex)
	//fmt.Println(sm.CommitIndex)
	//fmt.Println(actions5[1:])

	//Test6(Statemachine's term is less than the term in the Appendentry response)
	sm = &StateMachine{id: 1, peers: peers, term: 1, count: 0, PrevLogIndex: 4, LastLogIndex: 5, LastLogTerm: 2, CommitIndex: 4, nextIndex: nextindex, matchIndex: matchindex, status: "Leader"}
	sm.ProcessEvent(AppendEntriesResponseEv{Id: 5, Term: 2, Success: false})
	expect8(t, "Follower", sm.status, 2, sm.term)

	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(actions6[1:])
}

func TestLeaderVoteReq(t *testing.T) {

	fmt.Println("------------------------------------TestLeaderVoteReq-------------------------------------------------")
	actions := make([]interface{}, 1)
	peers := []int{2, 3, 4, 5}
	//(Test1) StateMachine's term is > than candidate term so return is error and status remains same
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Leader"}
	actions1 := sm.ProcessEvent(VoteRequestEv{Term: 1, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})

	actions = append(actions, error{"Error"})
	expect13(t, actions, actions1, "Leader", sm.status, 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(actions1[1:])
	//fmt.Println()

	//(Test2) StateMachine's term is < than candidate term so return is error and the term of SM gets changed to Follower
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, votedfor: 0, status: "Leader"}
	actions2 := sm.ProcessEvent(VoteRequestEv{Term: 2, LastLogIndex: 1, LastLogTerm: 1, CandidateId: 2})
	actions = make([]interface{}, 1)
	actions = append(actions, error{"Error"})
	expect13(t, actions, actions2, "Follower", sm.status, 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(actions2[1:])

	//fmt.Println()

}

func TestLeaderVoteResp(t *testing.T) {

	fmt.Println("----------------------------TestLeaderVoteResp-----------------------------------------------------")
	peers := []int{2, 3, 4, 5}
	actions := make([]interface{}, 1)
	//(Test1) statemachine's term is less than the term in  VoteResponseEv so return event is an error and term of SM gets changed
	//and it will become follower
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Leader"}
	actions1 := sm.ProcessEvent(VoteResponseEv{Term: 2, VoteGranted: false})
	actions = append(actions, error{"Error"})
	expect13(t, actions, actions1, "Follower", sm.status, 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(actions1[1:])
	//fmt.Println()
	//(Test2) statemachine's term is more than the term in VoteResponseEv so return event is an error and term of SM remains same
	sm = &StateMachine{id: 1, peers: peers, term: 2, LastLogIndex: 1, LastLogTerm: 1, status: "Leader"}
	actions2 := sm.ProcessEvent(VoteResponseEv{Term: 1, VoteGranted: true})
	actions = make([]interface{}, 1)
	actions = append(actions, error{"Error"})
	expect13(t, actions, actions2, "Leader", sm.status, 2, sm.term)
	//fmt.Println(sm.term)
	//fmt.Println(sm.status)
	//fmt.Println(actions2[1:])

}

func TestLeaderAppendReq(t *testing.T) {

	fmt.Println("-------------------------------TestLeaderAppendReq-----------------------------------------------------")

	actions := make([]interface{}, 1)
	nextindex := [6]int{0, 0, 5, 5, 5, 5}
	matchindex := [6]int{0, 5, 4, 4, 3, 2}

	entry1 := []byte{12}
	entry2 := []byte{13}
	entry3 := []byte{14}
	entry4 := []byte{15}
	entry5 := []byte{16}
	entry6 := []byte{17}

	log[0] = Log{0, 1, entry1}
	log[1] = Log{1, 1, entry2}
	log[2] = Log{2, 1, entry3}
	log[3] = Log{3, 2, entry4}
	log[4] = Log{4, 2, entry5}
	log[5] = Log{5, 2, entry6}

	Entry1 := []byte{123}
	peers := []int{2, 3, 4, 5}
	nxtindx := [6]int{0, 0, 6, 6, 6, 6}
	mindex := [6]int{0, 6, 4, 4, 3, 2}
	//Test1  response is appendentry request for all the peers with new data and
	sm = &StateMachine{id: 1, peers: peers, term: 2, PrevLogIndex: 4, PrevLogTerm: 2, LastLogIndex: 5, LastLogTerm: 2, status: "Leader", nextIndex: nextindex, matchIndex: matchindex, CommitIndex: 5}
	actions1 := sm.ProcessEvent(AppendEv{data: Entry1})
	actions = append(actions, LogStore{6, Entry1})
	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{2, 5, 2, 5, 1, Entry1}})

	}
	expect14(t, nxtindx, sm.nextIndex, mindex, sm.matchIndex, 6, sm.LastLogIndex, 5, sm.PrevLogIndex, actions, actions1)
	/*fmt.Println(sm.nextIndex)
	  fmt.Println(sm.matchIndex)
	  fmt.Println(sm.LastLogIndex)
	  fmt.Println(sm.PrevLogIndex)
	  fmt.Println(actions1[1:])*/

}
func TestLeadrTimeout(t *testing.T) {

	fmt.Println("-----------------------------------TestLeaderTimeout-----------------------------------------------------")

	peers := []int{2, 3, 4, 5}
	Entry1 := [6]int{0, 0, 0, 0, 0, 0}
	actions := make([]interface{}, 1)
	entry := []byte{}
	//Test1  send timeout event to follower so its term should increment and its status should also get changed
	sm = &StateMachine{id: 1, peers: peers, term: 1, LastLogIndex: 1, LastLogTerm: 1, status: "Leader", votecounter: 0, votetracker: Entry1}
	//It will send the heartbeat message to all the peers
	actions1 := sm.ProcessEvent(TimeoutEv{Time: 10})

	for i := 0; i <= 3; i++ {
		actions = append(actions, send{sm.peers[i], AppendEntriesRequestEv{1, 0, 1, 0, 1, entry}})
	}

	expect1(t, actions, actions1)
	//fmt.Println(actions1[1:])
}

func expect1(t *testing.T, x, y []interface{}) {

	z := reflect.DeepEqual(x, y)
	//fmt.Println(x)
	if z != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		fmt.Println("Test Failed")
	}
	fmt.Println("Test Passed")
}
func expect2(t *testing.T, x []interface{}, y []interface{}, i int, j int) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	if z != true || a != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}
}
func expect3(t *testing.T, x int, y int, a string, b string, p int, q int, c [6]int, d [6]int, actions []interface{}, actions1 []interface{}) {

	l := x == y
	g := a == b
	h := p == q
	m := reflect.DeepEqual(c, d)
	f := reflect.DeepEqual(actions, actions1)

	//fmt.Println(x)

	if l != true || g != true || h != true || m != true || f != true {
		//t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", a, b))
		t.Error(fmt.Sprintf("Expected %v, found %v", p, q))
		t.Error(fmt.Sprintf("Expected %v, found %v", c, d))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions1))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}
}

func expect4(t *testing.T, x []interface{}, y []interface{}, i string, j string) {

	a := reflect.DeepEqual(x, y)
	b := i == j
	if a != true || b != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}

}

func expect5(t *testing.T, x []interface{}, y []interface{}, i int, j int, m int, n int) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	b := reflect.DeepEqual(m, n)
	if z != true || a != true || b != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}
}

func expect6(t *testing.T, x string, y string, i int, j int, m [6]int, n [6]int) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	b := reflect.DeepEqual(m, n)
	if z != true || a != true || b != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect7(t *testing.T, x string, y string, i int, j int, m [6]int, n [6]int, actions []interface{}, actions2 []interface{}) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	b := reflect.DeepEqual(m, n)
	f := reflect.DeepEqual(actions, actions2)
	if z != true || a != true || b != true || f != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions2))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect8(t *testing.T, x string, y string, i int, j int) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)

	if z != true || a != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))

		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect9(t *testing.T, x string, y string, i int, j int, m int, n int) {

	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	b := reflect.DeepEqual(m, n)
	if z != true || a != true || b != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect10(t *testing.T, e int, l int, x string, y string, i int, j int, m [6]int, n [6]int, actions []interface{}, actions2 []interface{}) {

	g := reflect.DeepEqual(e, l)
	z := reflect.DeepEqual(x, y)
	a := reflect.DeepEqual(i, j)
	b := reflect.DeepEqual(m, n)
	f := reflect.DeepEqual(actions, actions2)
	if z != true || a != true || b != true || f != true || g != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", e, l))
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions2))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect11(t *testing.T, i int, j int, k [6]int, l [6]int, m [6]int, n [6]int, h int, p int, actions []interface{}, actions1 []interface{}) {

	g := reflect.DeepEqual(i, j)
	z := reflect.DeepEqual(k, l)
	a := reflect.DeepEqual(m, n)
	a1 := reflect.DeepEqual(h, p)
	a2 := reflect.DeepEqual(actions, actions1)
	if z != true || a != true || a1 != true || a2 != true || g != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", k, l))
		t.Error(fmt.Sprintf("Expected %v, found %v", h, p))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions1))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}
func expect12(t *testing.T, k [6]int, l [6]int, m [6]int, n [6]int, h int, p int, actions []interface{}, actions1 []interface{}) {

	z := reflect.DeepEqual(k, l)
	a := reflect.DeepEqual(m, n)
	a1 := reflect.DeepEqual(h, p)
	a2 := reflect.DeepEqual(actions, actions1)
	if z != true || a != true || a1 != true || a2 != true {

		t.Error(fmt.Sprintf("Expected %v, found %v", k, l))
		t.Error(fmt.Sprintf("Expected %v, found %v", h, p))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions1))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")

	}
}

func expect13(t *testing.T, x []interface{}, y []interface{}, i string, j string, m int, n int) {

	a := reflect.DeepEqual(x, y)
	b := i == j
	c := m == n
	if a != true || b != true || c != true {
		t.Error(fmt.Sprintf("Expected %v, found %v", x, y))
		t.Error(fmt.Sprintf("Expected %v, found %v", i, j))
		t.Error(fmt.Sprintf("Expected %v, found %v", m, n))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}

}
func expect14(t *testing.T, k [6]int, l [6]int, p1 [6]int, p2 [6]int, h int, p int, j int, l1 int, actions []interface{}, actions1 []interface{}) {

	z := reflect.DeepEqual(k, l)
	a1 := reflect.DeepEqual(h, p)
	m1 := reflect.DeepEqual(p1, p2)
	a2 := reflect.DeepEqual(actions, actions1)
	a3 := j == l1
	if z != true || a1 != true || a2 != true || a3 != true || m1 != true {

		t.Error(fmt.Sprintf("Expected %v, found %v", k, l))
		t.Error(fmt.Sprintf("Expected %v, found %v", h, p))
		t.Error(fmt.Sprintf("Expected %v, found %v", j, l1))
		t.Error(fmt.Sprintf("Expected %v, found %v", actions, actions1))
		t.Error(fmt.Sprintf("Expected %v, found %v", p1, p2))
		fmt.Println("Test Failed")
	} else {
		fmt.Println("Test Passed")
	}
}

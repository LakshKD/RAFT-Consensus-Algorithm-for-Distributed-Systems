package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	cluster "github.com/cs733-iitb/cluster"
	log1 "github.com/cs733-iitb/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
)

type Event interface{}

type Node interface {

	// Client's message to Raft node
	Append([]byte)

	// A channel for client to listen on. What goes into Append must come out of here at some point.
	CommitChannel() <-chan CommitInfo

	// Last known committed index in the log.
	CommittedIndex() int //This could be -1 until the system stabilizes.

	// Returns the data at a log index, or an error.
	Get(index int64) (interface{}, error)

	// Node's id
	Id() int

	// Id of leader. -1 if unknown
	LeaderId() int

	// Signal to shut down all goroutines, stop sockets, flush log and close it, cancel timers.
	Shutdown()
}

// Index is valid only if err == nil
type CommitInfo struct {
	Data  []byte
	Index int
	error string
}

// This is an example structure for Config .. change it to your convenience.
type Config struct {
	Cluster          []NetConfig // Information about all servers, including this.
	Id               int         // this node's id. One of the cluster's entries should match.
	LogDir           string      // Log file directory for this node
	ElectionTimeout  int
	HeartbeatTimeout int
}

func ConfigRaft(configuration interface{}) (config *Config, err error) {
	var cfg Config
	var ok bool
	var configFile string
	//var err error
	if configFile, ok = configuration.(string); ok {
		var f *os.File
		if f, err = os.Open(configFile); err != nil {
			return nil, err
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		if err = dec.Decode(&cfg); err != nil {
			return nil, err
		}
	} else if cfg, ok = configuration.(Config); !ok {
		return nil, errors.New("Expected a configuration.json file or a Config structure")
	}
	return &cfg, nil
}

type NetConfig struct {
	Id   int
	Host string
	Port int
}

type RaftNode struct { // implements Node interface

	config    *Config
	sm        *StateMachine
	srv       cluster.Server
	lg        log1.Log
	eventCh   chan Event
	Heart_Beat_Timer *time.Ticker
	Timer     *time.Timer
	timeoutCh chan interface{}
	shutdown  bool
	CommitCh  chan CommitInfo
}

// It inits the cluster (using the cluster package discussed above), inits the log, reads/entries from the log,
//creates a state machine in follower mode with the log, reads the node-specific file that stores lastVotedFor and term,

func New1(myid int, configuration interface{}) (n Node, err error) {

	Raft := new(RaftNode)
	n = Raft
	var config *Config
	var pdata []byte
	if config, err = ConfigRaft(configuration); err != nil {
		return nil, err
	}
	config.Id = myid

    
    pdata, err = ioutil.ReadFile(fmt.Sprintf("persistent_store_%d", myid))
    if err != nil {
			panic(err)
		}
	p := strings.Split(string(pdata), " ")
	p[0] = strings.TrimSpace(p[0])
	p[1] = strings.TrimSpace(p[1])

	smterm, _ := strconv.Atoi(p[1])
	smvotedfor, _ := strconv.Atoi(p[0])
	

	server, err := cluster.New(myid, "cluster_test_config.json")

	if err != nil {
		panic(err)
	}

	Raft.srv = server

	Raft.config = config
	nodesid := [4]int{}
	j := 0
	for i := 0; i <= 4; i++ {

		if config.Cluster[i].Id != myid {
			nodesid[j] = config.Cluster[i].Id
			j = j + 1
		}
	}

    next:=[]int{1,1,1,1,1,1}
	sm = &StateMachine{id: myid, term: smterm, LastLogIndex: 0,CommitIndex:0,PrevLogIndex: -1,nextIndex:next,peers: nodesid, votedfor: smvotedfor, status: "Follower", LeaderID: 0,ElectionTimeout:config.ElectionTimeout,HeartBeatTimeout:config.HeartbeatTimeout}
    
	Raft.sm = sm
    rand.Seed(time.Now().UnixNano())
    RandomNo := rand.Intn(500)
    Raft.Timer = time.AfterFunc(time.Duration(config.ElectionTimeout+RandomNo)*time.Millisecond, func() { Raft.timeoutCh <- TimeoutEv{}})

    Raft.Heart_Beat_Timer = time.NewTicker(time.Duration(config.HeartbeatTimeout) * time.Millisecond)
	go Raft.processEvents()
	go Raft.sm.ProcessEvent(make([]interface{}, 1))
	
	Events := make(chan Event, 10)
	commitchannel := make(chan CommitInfo, 1)
	timeout := make(chan interface{})
	Raft.eventCh = Events
	Raft.CommitCh = commitchannel
	Raft.timeoutCh = timeout
    config.LogDir = fmt.Sprintf("log_File_%d", myid)
	lg, _ := log1.Open(config.LogDir)
	
	if lg == nil {
		fmt.Println("Error opening Log File")
	}
     
    
	Raft.lg = *lg
	Raft.sm.log=&Raft.lg
	//Raft.lg.TruncateToEnd(0)
	Raft.lg. RegisterSampleEntry(Log{})
	er:=Raft.lg.Append(Log{Logindex:0,Term:0,Entry:[]byte("foo")})
	 if(er!=nil){
	 	fmt.Println(er)
	 }
   
	
	return
}

func (rn *RaftNode) Append(data []byte) {
   //fmt.Println("Inside Append")
	go func() {

		rn.eventCh <- AppendEv{data: data}

	}()
}

func (rn *RaftNode) Get(index int64)(interface{}, error) {
	return rn.lg.Get(index)

}

func (rn *RaftNode) Id() int {

	return rn.config.Id

}

func (rn *RaftNode) CommittedIndex() int {
	return rn.sm.CommitIndex
}

func (rn *RaftNode) LeaderId() int {

	return rn.sm.LeaderID

}

func (rn *RaftNode) Shutdown() {
	rn.shutdown = true
	rn.srv.Close()
}

func (rn *RaftNode) CommitChannel() <-chan CommitInfo {

	return rn.CommitCh
}


func (rn *RaftNode) doActions(ev []interface{}) {

	gob.Register(AppendEntriesRequestEv{})
	gob.Register(VoteRequestEv{})
	gob.Register(VoteResponseEv{})
	gob.Register(AppendEntriesResponseEv{})
	length := len(ev)
	var idd int64
	idd = 1
	for i := 1; i < length; i++ {

		switch ev[i].(type) {

		case send:
			l := ev[i].(send)
			rn.srv.Outbox() <- &cluster.Envelope{Pid: l.DestID, MsgId: idd, Msg: l.Event}
			//time.Sleep(2 * time.Second)

		case LogStore:
			l := ev[i].(LogStore)
			rn.lg.Append(Log{Logindex:l.index,Term:rn.sm.term,Entry:l.entry})

		case Commit:
			
			c := ev[i].(Commit)
			//fmt.Println("Inside Commit",rn.sm.id,c)
			rn.CommitCh <- CommitInfo{c.data, c.index, c.error}

		case Alarm:
			A := ev[i].(Alarm)
           rn.Timer.Stop()
           rand.Seed(time.Now().UnixNano())
           RandomNo := rand.Intn(500)
		   rn.Timer = time.AfterFunc(time.Duration(A.t+RandomNo)*time.Millisecond, func() { rn.timeoutCh <- TimeoutEv{} })
    	
    }

	}
	idd++
}

func (rn *RaftNode) processEvents() {

	go func() {

		if rn.shutdown == true {
			//fmt.Println("ITs way true")
			return

		}

	}()
	go func() {
		for {
			d1 := <-rn.srv.Inbox()
			actions1 := rn.sm.ProcessEvent(d1.Msg)
			rn.doActions(actions1)
			
		}
	}()

	for {
		var ev Event
		select {
		case ev = <-rn.eventCh:

		case ev = <-rn.timeoutCh:

		case <-rn.Heart_Beat_Timer.C:
			
			if rn.sm.status == "Leader" {
				rn.eventCh <- TimeoutEv{}

			}
	     

		}
		
		actions2 := rn.sm.ProcessEvent(ev)
		rn.doActions(actions2)
		
		
	}

}

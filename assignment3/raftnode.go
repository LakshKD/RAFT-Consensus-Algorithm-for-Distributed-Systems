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
	Get(index int64) ([]byte, error)

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
	timeoutCh chan int
	shutdown  bool
	CommitCh  chan CommitInfo
}

// It inits the cluster (using the cluster package discussed above), inits the log, reads/entries from the log,
//creates a state machine in follower mode with the log, reads the node-specific file that stores lastVotedFor and term,

func New1(myid int, configuration interface{}) (n Node, err error) {
	//var Raft RaftNode
	//fmt.Println(myid)
	Raft := new(RaftNode)
	n = Raft
	var config *Config
	// var err1 error
	//var p [2]string
	var pdata []byte
	if config, err = ConfigRaft(configuration); err != nil {
		return nil, err
	}
	config.Id = myid
	//fmt.Println(config.Id)
	if config.Id == 1 {

		config.LogDir = "log_file_1"

	} else if config.Id == 2 {

		config.LogDir = "log_file_2"

	} else if config.Id == 3 {

		config.LogDir = "log_file_3"

	} else if config.Id == 4 {

		config.LogDir = "log_file_4"

	} else if config.Id == 5 {

		config.LogDir = "log_file_5"
	}
	//  Reading the LastvotedFor and Term from the Persistent Store

	if myid == 1 {

		pdata, err = ioutil.ReadFile("persistent_store_1")
		if err != nil {
			panic(err)
		}

	} else if myid == 2 {
		pdata, err = ioutil.ReadFile("persistent_store_2")
		if err != nil {
			panic(err)
		}

	} else if myid == 3 {

		pdata, err = ioutil.ReadFile("persistent_store_3")
		if err != nil {
			panic(err)
		}

	} else if myid == 4 {

		pdata, err = ioutil.ReadFile("persistent_store_4")
		if err != nil {
			panic(err)
		}

	} else if myid == 5 {

		pdata, err = ioutil.ReadFile("persistent_store_5")
		if err != nil {
			panic(err)
		}

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

	sm = &StateMachine{id: myid, term: smterm, LastLogIndex: 0, PrevLogIndex: 0, peers: nodesid, votedfor: smvotedfor, status: "Follower", LeaderID: 0}

	Raft.sm = sm

	go Raft.processEvents()
	go Raft.sm.ProcessEvent(make([]interface{}, 1))
	//time.Sleep(2*time.Second)
	Events := make(chan Event, 10)
	commitchannel := make(chan CommitInfo, 1)
	timeout := make(chan int)
	Raft.eventCh = Events
	Raft.CommitCh = commitchannel
	Raft.timeoutCh = timeout

	lg, _ := log1.Open(config.LogDir)
	//fmt.Println(config.LogDir)
	if lg == nil {
		fmt.Println("Error opening Log File")
	}

	Raft.lg = *lg
	var k int64
	for k = 0; k <= Raft.lg.GetLastIndex(); k++ {

		Raft.sm.log[int(k)].logindex = int(k)
		Raft.sm.log[int(k)].entry, _ = Raft.lg.Get(k)
		//fmt.Println("abc")

	}
	//Raft.lg.TruncateToEnd(0)

	// fmt.Println("bb")
	return
}

func (rn *RaftNode) Append(data []byte) {

	go func() {

		rn.eventCh <- AppendEv{data: data}

	}()
}

func (rn *RaftNode) Get(index int64) ([]byte, error) {
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
	//fmt.Println("Inside doActions")
	length := len(ev)
	//fmt.Println(ev)
	var idd int64
	idd = 1
	for i := 1; i < length; i++ {

		switch ev[i].(type) {

		case send:
			l := ev[i].(send)
			rn.srv.Outbox() <- &cluster.Envelope{Pid: l.DestID, MsgId: idd, Msg: l.Event}
			time.Sleep(1 * time.Second)

		case LogStore:
			l := ev[i].(LogStore)
			rn.lg.Append([]byte(l.entry))

		case Commit:
			c := ev[i].(Commit)
			rn.CommitCh <- CommitInfo{c.data, c.index, c.error}

		case Alarm:
			A := ev[i].(Alarm)
			rn.timeoutCh <- A.t

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
			//fmt.Println(rn.sm.id,d1)
			//fmt.Println(d1)
			actions1 := rn.sm.ProcessEvent(d1.Msg)

			//fmt.Println("Actions along with ID of different state machines",rn.sm.id,actions1)
			rn.doActions(actions1)
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		var ev Event
		select {
		case ev = <-rn.eventCh:

		case ev = <-rn.timeoutCh:

		}
		//fmt.Println(ev)
		actions2 := rn.sm.ProcessEvent(ev)

		//fmt.Println("Heloo",actions2)
		rn.doActions(actions2)
	}

}

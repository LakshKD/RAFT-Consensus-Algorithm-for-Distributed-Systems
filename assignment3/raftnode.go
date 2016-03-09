package main

import (
	"io/ioutil"
	"os"
	"strings"
)

var sm *StateMachine

// Index is valid only if err == nil
type CommitInfo struct {
	Data  []byte
	Index int64 // or int .. whatever you have in your code
	Err   error // Err can be errred
}

// This is an example structure for Config .. change it to your convenience.
type Config struct {
	cluster          []NetConfig // Information about all servers, including this.
	Id               int         // this node's id. One of the cluster's entries should match.
	LogDir           string      // Log file directory for this node
	ElectionTimeout  int
	HeartbeatTimeout int
}

type NetConfig struct {
	Id   int
	Host string
	Port int
}

type RaftNode struct { // implements Node interface
	eventCh   chan Event
	timeoutCh <-chan Time
}

type Node interface {

	// Client's message to Raft node
	Append([]byte)

	// A channel for client to listen on. What goes into Append must come out of here at some point.
	CommitChannel() <-chan CommitInfo

	// Last known committed index in the log.
	CommittedIndex() int //This could be -1 until the system stabilizes.

	// Returns the data at a log index, or an error.
	Get(index int) (err, []byte)

	// Node's id
	Id()

	// Id of leader. -1 if unknown
	LeaderId() int

	// Signal to shut down all goroutines, stop sockets, flush log and close it, cancel timers.
	Shutdown()
}

func (rn *RaftNode) Append(data) {
	rn.eventCh <- Append{data: data}

}

func (rn *RaftNode) Get(index int) (err, []byte) {
     return lg.ldb.Get(toBytes(index), nil)
}

func (rn *RaftNode) Id() int {

}

func (rn *RaftNode) LeaderId() int {

}

func (rn *RaftNode) Shutdown() {

}

// It inits the cluster (using the cluster package discussed above), inits the log, reads/entries from the log, creates a state machine in
//follower mode with the log, reads the node-specific file that stores lastVotedFor and term,
func New(c Config) Node {

	fo, err := os.Create(c.LogDir)
	if err != nil {
		panic(err)
	}

	if c.Id == 1 {
		dat, err := ioutil.ReadFile("persistent_store_1")
		if err != nil {
			panic(err)
		}

	} else if c.Id == 2 {
		dat, err := ioutil.ReadFile("persistent_store_2")
		if err != nil {
			panic(err)
		}

	} else if c.Id == 3 {

		dat, err := ioutil.ReadFile("persistent_store_3")
		if err != nil {
			panic(err)
		}

	} else if c.Id == 4 {

		dat, err := ioutil.ReadFile("persistent_store_4")
		if err != nil {
			panic(err)
		}

	} else if c.Id == 5 {

		dat, err := ioutil.ReadFile("persistent_store_5")
		if err != nil {
			panic(err)
		}
	}

	p := strings.Split(string(dat), " ")
	p[0] = strings.Trim(p[0])
	p[1] = strings.Trim(p[1])
	sm = &StateMachine{id: c.Id, peers: c.cluster, term: p[1], PrevLogIndex: 5, PrevLogTerm: 2, CommitIndex: 2, LastLogIndex: 6, LastLogTerm: 2, votedfor: p[0], status: "Follower"}

}

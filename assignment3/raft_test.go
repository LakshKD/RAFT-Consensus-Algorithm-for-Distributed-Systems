package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
    "time"
	"sync"
	"github.com/cs733-iitb/cluster"
)
var RaftNode_Obj []Node    //array containing objects of Raft node

machine1 ,err1 := cluster.New(1, "cluster_test_config1.json")
 if err1 != nil {
		panic(err)
	}

machine2 ,err2 := cluster.New(2, "cluster_test_config1.json")
 if err2 != nil {
		panic(err)
	}
machine3 ,err3 := cluster.New(3, "cluster_test_config1.json")
if err3 != nil {
		panic(err)
	}
machine4 ,err4 := cluster.New(4, "cluster_test_config1.json")
 if err4 != nil {
		panic(err)
}
machine5 ,err5 := cluster.New(5, "cluster_test_config1.json")
 
 if err5 != nil {
		panic(err)
}


RaftNode_obj[0]=New(Config{machine1,1,"log_file_1",150,100})
RaftNode_obj[1]=New(Config{machine2,2,"log_file_2",150,100})
RaftNode_obj[2]=New(Config{machine3,3,"log_file_3",150,100})
RaftNode_obj[3]=New(Config{machine4,4,"log_file_4",150,100})
RaftNode_obj[4]=New(Config{machine5,5,"log_file_5",150,100})


// Simple serial check of getting and setting
func TestTCPSimple(t *testing.T) {
	//peers, err := cluster.New(1, "cluster_test_config1.json")
  
    rafts := makeRafts() // array of []raft.Node
	
	ldr := getLeader(rafts)
	ldr.Append("foo")
	time.Sleep(1 time.Second)
	for _, node:= rafts {
		select {
		// to avoid blocking on channel.
		case ci := <- node.CommitChannel():
		if ci.err != nil {t.Fatal(ci.err)}
		if string(ci.data) != "foo" {
			t.Fatal("Got different data")
			}
		default: t.Fatal("Expected message on all nodes")
		}
	}	
}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

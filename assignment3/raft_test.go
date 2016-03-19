package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	//cluster "github.com/cs733-iitb/cluster"
	//log1 "github.com/cs733-iitb/log"
)

func makeRafts() [5]Node {
	var RaftNodeObj [5]Node
	//fmt.Println("adfad")
	RaftNodeObj[0], _ = New1(1, "cluster_test_config1.json")
	RaftNodeObj[1], _ = New1(2, "cluster_test_config1.json")
	RaftNodeObj[2], _ = New1(3, "cluster_test_config1.json")
	RaftNodeObj[3], _ = New1(4, "cluster_test_config1.json")
	RaftNodeObj[4], _ = New1(5, "cluster_test_config1.json")

	return RaftNodeObj

}

var rafts [5]Node

func getLeader(RaftNodeObj [5]Node) Node {

	//Raft := new(RaftNode)

	//var config *Config
	idcountarray := [6]int{0, 0, 0, 0, 0, 0}

	for i := 0; i <= 4; i++ {
		//Raft := new(RaftNode)
		//RaftNodeObj[i]=Raft
		p := RaftNodeObj[i].(*RaftNode).LeaderId()
		idcountarray[p]++
	}
	max := idcountarray[0]
	index := 0
	for j := 1; j <= 5; j++ {
		if idcountarray[j] > max {
			index = j
		}
	}
	//fmt.Println("Inside getleader")
	//fmt.Println(index)
	for k := 0; k <= 4; k++ {
		//Raft1 := new(RaftNode)
		//RaftNodeObj[k]=Raft1
		if RaftNodeObj[k].(*RaftNode).config.Id == index {
			return RaftNodeObj[k]
		}

	}

	return nil
}
func TestMakingLeader(t *testing.T) {
	fmt.Println("----------------------------------Leader Formation----------------------------------------")

	rafts = makeRafts()
	fmt.Println("Elections are going on")
	// fmt.Println(rafts[1].(*RaftNode).sm)   Intially the state is follower
	actions := rafts[1].(*RaftNode).sm.ProcessEvent(TimeoutEv{rand.Intn(10)})
	//fmt.Println(rafts[1].(*RaftNode).sm) Then the state becomes Candidate
	rafts[1].(*RaftNode).doActions(actions) //Node 2 will send the voterequest to all the peer nodes
	time.Sleep(10 * time.Second)            // Elections are going on here
	ldr := getLeader(rafts)                 // Finally got the leader and checking the status and leader Id of all the other nodes

	expectedcountLeaderID := 3 //LeaderId of majority will be 2
	Actualcountofleaderid := 0
	for i := 0; i <= 4; i++ {
		if rafts[i].(*RaftNode).sm.LeaderID == 2 {
			Actualcountofleaderid++
		}
	}
	expectedstatus := "Leader"
	expect1(t, ldr.(*RaftNode).sm.status, expectedstatus, Actualcountofleaderid, expectedcountLeaderID)

}

func TestAppend_ON_Leader(t *testing.T) {

	fmt.Println("--------------------------Checking Append On Leader----------------------------------")

	ldr := getLeader(rafts) //ID number 2 becomes the leader and send the append request to all the machines
	time.Sleep(5 * time.Second)
	data := []byte{90}
	ldr.(*RaftNode).Append(data)

	time.Sleep(10 * time.Second)
	for _, node := range rafts {
		select {
		// to avoid blocking on channel.
		case ci := <-node.CommitChannel():
			if len(ci.error) > 0 {
				t.Fatal(ci.error)
			}
			if string(ci.Data) != "90" {
				t.Fatal("Got different data")
			}
		default:
			t.Fatal("Expected message on all nodes")

		} //select end
	} //for end
}
func TestGetfunction(t *testing.T) {
	fmt.Println("------------------------------------Checking Get---------------------------------------")
	data, _ := rafts[1].(*RaftNode).lg.Get(int64(0))
	expecteddata := []byte{90}
	expect3(t, expecteddata, data)
}
func TestShutdownfunction(t *testing.T) {

	fmt.Println("---------------------------------Checking Shutdown------------------------------------")
	rafts[1].(*RaftNode).Shutdown()
	expected := true
	expect4(t, expected, rafts[1].(*RaftNode).shutdown)
}

func TestLeaderIDfunction(t *testing.T) {

	fmt.Println("--------------------------Checking Leader ID function----------------------------------")
	id := rafts[2].(*RaftNode).sm.LeaderID //Checking Leader ID of node 3
	expected := 2
	expect5(t, expected, id)

	id = rafts[0].(*RaftNode).sm.LeaderID //Checking Leader ID of node 1
	expect5(t, expected, id)

	id = rafts[3].(*RaftNode).sm.LeaderID //Checking Leader ID of node 4
	expect5(t, expected, id)

	id = rafts[4].(*RaftNode).sm.LeaderID //Checking Leader ID of node 5
	expect5(t, expected, id)

}

func TestCommitIndex(t *testing.T) {

	fmt.Println("--------------------------Checking Commit Index Function----------------------------------")
	//Checking the commitindex of Leader
	a1 := rafts[1].(*RaftNode).sm.CommitIndex
	expected := 1

	expect6(t, expected, a1)

}

func TestID(t *testing.T) {

	fmt.Println("----------------------------Checking ID function----------------------------------------")
	//Checking the Id of a Node
	a1 := rafts[0].(*RaftNode).Id()
	expected := 1

	expect7(t, expected, a1)

}
func expect1(t *testing.T, a string, b string, c int, d int) {

	p := a == b
	q := c >= d

	if p == false || q == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`

		t.Error(fmt.Sprintf("Expected %v, found %v", d, c))
	} else {
		fmt.Println("Test case 1 (Formation of Leader)Passed")
	}

}

func expect2(t *testing.T, a string, b string) {

	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Testcase 2 (Append call on leader)passed")
	}

}
func expect3(t *testing.T, a []byte, d []byte) {

	p := reflect.DeepEqual(a, d)

	if p == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", a, d)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Test case 3 (Get data)Passed")
	}

}

func expect4(t *testing.T, a bool, d bool) {

	p := a == d

	if p == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", a, d)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Test case 4 (Shutting Down the node)Passed")
	}

}
func expect5(t *testing.T, a int, d int) {

	p := a == d

	if p == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", a, d)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Test case 5 (LeaderID)Passed")
	}

}

func expect6(t *testing.T, a int, d int) {

	p := a == d

	if p == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", a, d)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Test case 6 (LeaderID)Passed")
	}

}

func expect7(t *testing.T, a int, d int) {

	p := a == d

	if p == false {
		t.Error(fmt.Sprintf("Expected %v, found %v", a, d)) // t.Error is visible when running `go test -verbose`

	} else {
		fmt.Println("Test case 7 (ID of a Node)Passed")
	}

}

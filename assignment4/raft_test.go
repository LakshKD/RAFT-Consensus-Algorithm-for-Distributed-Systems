package main

import (
	"fmt"
	"math/rand"
	_"reflect"
	"testing"
	"time"
	"os"
	"strconv"
	_"io/ioutil"
	_"strings"

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
	fmt.Println("<----------------------------------Leader Formation---------------------------------------->")
logfile := []string{"log_File_1","log_File_2","log_File_3","log_File_4","log_File_5"}
	for k:=0;k<=4;k++{
	    os.RemoveAll(logfile[k])
	 }  
   s:=[]string{"persistent_store_1","persistent_store_2","persistent_store_3","persistent_store_4","persistent_store_5"}

   for m:=0;m<=4;m++{
   	 err := os.Remove(s[m])
   	 if err !=nil{
   	 	fmt.Println(err)
   	 }
   }


    for i:=1 ;i<=5 ;i++ {

     switch(i){

    	 case 1: //a:=rand.Intn(5)
    			f, err := os.Create("persistent_store_1")
    			if(err!=nil){
    			fmt.Println("Error opening file")
    			}
    			_,err = f.Write([]byte(strconv.Itoa(0)))
    			b:=rand.Intn(5)
    			_,err = f.Write([]byte(" "))
    			_,err = f.Write([]byte(strconv.Itoa(b)))
    			f.Close()
      	 
      	 case 2: //a:=rand.Intn(5)
    			f, err := os.Create("persistent_store_2")
    			if(err!=nil){
    			fmt.Println("Error opening file")
    			}
    			_,err = f.Write([]byte(strconv.Itoa(0)))
    			b:=rand.Intn(10)
    			_,err = f.Write([]byte(" "))
    			_,err = f.Write([]byte(strconv.Itoa(b)))
    			f.Close()

    	 case 3: //a:=rand.Intn(5)
    			f, err := os.Create("persistent_store_3")
    			if(err!=nil){
    			fmt.Println("Error opening file")
    			}
    			_,err = f.Write([]byte(strconv.Itoa(0)))
    			b:=rand.Intn(15)
    			_,err = f.Write([]byte(" "))
    			_,err = f.Write([]byte(strconv.Itoa(b)))
    			f.Close()		
          
          case 4://a:=rand.Intn(5)
    			f, err := os.Create("persistent_store_4")
    			if(err!=nil){
    			fmt.Println("Error opening file")
    			}
    			_,err = f.Write([]byte(strconv.Itoa(0)))
    			b:=rand.Intn(20)
    			_,err = f.Write([]byte(" "))
    			_,err = f.Write([]byte(strconv.Itoa(b)))
    			f.Close() 

    	  case 5://a:=rand.Intn(5)
    			f, err := os.Create("persistent_store_5")
    			if(err!=nil){
    			fmt.Println("Error opening file")
    			}
    			_,err = f.Write([]byte(strconv.Itoa(0)))
    			b:=rand.Intn(25)
    			_,err = f.Write([]byte(" "))
    			_,err = f.Write([]byte(strconv.Itoa(b)))
    			f.Close()		


      	}	
    }	
    rafts = makeRafts()
	//fmt.Println("Elections are going on")
	time.Sleep(5 * time.Second)            // Elections are going on here
	ldr := getLeader(rafts)                 // Finally got the leader and checking the status and leader Id of all the other nodes

	expectedcountLeaderID := 3 //LeaderId of majority will be 2
	Actualcountofleaderid := 0
    /*for i := 0; i <= 4; i++ {
		//fmt.Printf("%+v\n",rafts[i].(*RaftNode).sm) 
	}*/



	for i := 0; i <= 4; i++ {
		if rafts[i].(*RaftNode).sm.LeaderID == ldr.(*RaftNode).sm.LeaderID {
			Actualcountofleaderid++
		}
	}
	expectedstatus := "Leader"
	expect1(t, ldr.(*RaftNode).sm.status, expectedstatus, Actualcountofleaderid, expectedcountLeaderID)
}




func TestAppend_ON_Leader(t *testing.T) {
	 
	fmt.Println("<--------------------------Checking Single Append On Leader---------------------------------->")
    var flag int
    flag=0
    ldr := getLeader(rafts) 
    //var j int64
   

	time.Sleep(5 * time.Second)
	ldr.(*RaftNode).Append([]byte("90"))

	time.Sleep(10 * time.Second)
	
	for _, node := range rafts {
		select {
		// to avoid blocking on channel.
		case ci := <-node.CommitChannel():
			//fmt.Println(node.(*RaftNode).sm.id,ci)
			if len(ci.error) > 0 {
				t.Fatal(ci.error)
				flag=1
			}
			if string(ci.Data) != "90" {
				t.Fatal("Got different data")
				flag=1
			}
	     case <- time.After(time.Second * time.Duration(10)):

			t.Fatal("Expected message on all nodes")
			flag=1

		} 
	} 
    if(flag!=1){
    	fmt.Println("Testcase 2 (Append-1 on Leader) Passed")
    }else{
    	fmt.Println("Testcase 2 (Append-1 on Leader) Failed")
    }
    /*for i:=0; i<=4;i++{
    	
       
        fmt.Println(rafts[i].(*RaftNode).sm.status)	
         
        for j=0;j<=rafts[i].(*RaftNode).sm.log.GetLastIndex();j++ {
         data,_:=rafts[i].(*RaftNode).sm.log.Get(j)
         fmt.Println(data)
        }
    	
    }
    for i := 0; i <= 4; i++ {
		fmt.Printf("%+v\n",rafts[i].(*RaftNode).sm) 
	}*/
    
}
/*
func TestMultiple_Append_ON_Leader(t *testing.T){
   fmt.Println("-----------------------------Checking Multiple Append On Leader--------------------------")
   var flag int 
   var j int64
   flag = 0
   data := []string{"test1","test2","test3"}
   var num_of_Test=3
  for i:=0; i<=2;i++ {
   ldr := getLeader(rafts) 
   time.Sleep(5*time.Second)
   ldr.(*RaftNode).Append([]byte(data[i]))
    
   time.Sleep(10*time.Second)
   for _, node := range rafts {
		select {
		// to avoid blocking on channel.
		case ci := <-node.CommitChannel():
			fmt.Print(ci)
			if len(ci.error) > 0 {
				t.Fatal(ci.error)
				flag=1
			}
			if string(ci.Data) != data[i] {
				t.Fatal("Got different data")
				flag=1
			}
		case <- time.After(time.Second * time.Duration(10)):
			t.Fatal("Expected message on all nodes")
			flag=1

		} 
	}
    if(flag!=1){
    	fmt.Println("Testcase" ,num_of_Test,"(Append",num_of_Test,"on Leader) Passed")
    }else{
    	fmt.Println("Testcase" ,num_of_Test,"(Append",num_of_Test,"on Leader) Failed")
    	flag = 0
    }
    num_of_Test++
   }
  for i:=0; i<=4;i++{
    	
       //rafts[i].(*RaftNode).lg.Open(dpath)
        fmt.Println(rafts[i].(*RaftNode).sm.status)	
         
        for j=0;j<=rafts[i].(*RaftNode).sm.log.GetLastIndex();j++ {
         data,_:=rafts[i].(*RaftNode).sm.log.Get(j)
         fmt.Println(data)
        }
    	//data,err:=rafts[i].(*RaftNode).lg.Get(int64(0))
    	//fmt.Println(data,err)
    }
  for i := 0; i <= 4; i++ {
		fmt.Printf("%+v\n",rafts[i].(*RaftNode).sm) 
	}


} 
*/
/*
func TestGetfunction(t *testing.T) {
	fmt.Println("<---------------------------------------Checking Get----------------------------------------->")
	data, _ := rafts[1].(*RaftNode).lg.Get(int64(0))
	//expecteddata := []byte{90}
	expect3(t, []byte("foo"), data)
}
func TestShutdownfunction(t *testing.T) {

	fmt.Println("<-------------------------------------Checking Shutdown-------------------------------------->")
	rafts[1].(*RaftNode).Shutdown()
	expected := true
	expect4(t, expected, rafts[1].(*RaftNode).shutdown)
}

func TestLeaderIDfunction(t *testing.T) {

	fmt.Println("<-------------------------------Checking Leader ID function---------------------------------->")
	id := rafts[2].(*RaftNode).sm.LeaderID //Checking Leader ID of node 3
	ldr := getLeader(rafts)
    expected := ldr.(*RaftNode).sm.LeaderID
	expect5(t, expected, id)

	id = rafts[0].(*RaftNode).sm.LeaderID //Checking Leader ID of node 1
	expect5(t, expected, id)

	id = rafts[1].(*RaftNode).sm.LeaderID //Checking Leader ID of node 2
	expect5(t, expected, id)

	id = rafts[4].(*RaftNode).sm.LeaderID //Checking Leader ID of node 5
	expect5(t, expected, id)

	id = rafts[3].(*RaftNode).sm.LeaderID  //Checking Leader ID of node 4
	expect5(t,expected,id)

}

func TestCommitIndex(t *testing.T) {

	fmt.Println("<----------------------------Checking Commit Index Function------------------------------------>")
	a1 := rafts[1].(*RaftNode).sm.CommitIndex
	expected := 1

	expect6(t, expected, a1)

}

func TestID(t *testing.T) {

	fmt.Println("<-----------------------------Checking ID function--------------------------------------------->")
	//Checking the Id of a Node
	a1 := rafts[0].(*RaftNode).Id()
	expected := 1

	expect7(t, expected, a1)

}

/*func TestStateMachine(t *testing.T){
   fmt.Println("------------------------------Checking StateMachines-------------------------------------------")
   for i:=0;i<=4;i++{
    fmt.Println(rafts[i].(*RaftNode).sm)
    }
}*/



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
/*
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
*/
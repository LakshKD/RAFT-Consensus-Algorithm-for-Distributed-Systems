Name:- Lakshya Kumar'
Roll Number:- 153050051

Problem Definition

In this assignment we have to build a Raftnode which encapsulates the Raft statemachine(created in assignment 2). We have to build the functionalities like 
1 Function to intialize the RaftNodes and start them.
2.Function to select a leader and performs the append function.




Code Description:- 

*There is a node interface which contains the functions :-
        1.Append([]byte)

	2. // A channel for client to listen on. What goes into Append must come out of here at some point.
	  CommitChannel() <-chan CommitInfo

	3. // Last known committed index in the log.
	   CommittedIndex() int //This could be -1 until the system stabilizes.

	4.// Returns the data at a log index, or an error.
	   Get(index int64) ([]byte, error)

	5.// Node's id
	   Id()  int

	6.// Id of leader. -1 if unknown
	  LeaderId() int

	7.// Signal to shut down all goroutines, stop sockets, flush log and close it, cancel timers.
	  Shutdown()

* New1() This function is creating the Raftnodes and intializing them and also starting their processevents function.

* DoActions() Each RaftNode is having the DoActions function that will help the RaftNode to perform the Actions that will come through Event channel or through the input channel.

*ProcessEvent() This function of the raftnode.go file is helping the Raftnode to perform the actions that will come to the Raftnode.


Testing of Raftnode

1.func TestMakingLeader():- Testing the formation of the Leader.

2. func TestAppend_ON_Leader():- Testing the Append functionality on the Leader and it is showing some error in the commit at leader.

3.TestGetfunction() :- Testing to get the data from the Log 

4.TestShutdownfunction():- Testing Shutdown functionality of the Raftnode.

5.TestLeaderIDfunction():-Testing to get the LeaderId from the nodes in the cluster.

6.TestCommitIndex():- Testing to get the commit index(This test case is failing due to some reason).

Note :- I have not tested the mock cluster functionality due to shortage of time.



 





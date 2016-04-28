package main

import (
	"bufio"
	"fmt"
	"github.com/LakshKD/cs733/assignment4/raft/fs"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
)
var id int 
var crlf = []byte{'\r', '\n'}
var port = 8080

func check(obj interface{}) {
	if obj != nil {
		//fmt.Println(obj)
		os.Exit(1)
	}
}

func reply(conn *net.TCPConn, msg *fs.Msg) bool {
	var err error
	write := func(data []byte) {
		if err != nil {
			return
		}
		_, err = conn.Write(data)
	}
	var resp string
	switch msg.Kind {
	case 'C': // read response
		resp = fmt.Sprintf("CONTENTS %d %d %d", msg.Version, msg.Numbytes, msg.Exptime)
	case 'O':
		resp = "OK "
		if msg.Version > 0 {
			resp += strconv.Itoa(msg.Version)
		}
	case 'F':
		resp = "ERR_FILE_NOT_FOUND"
	case 'V':
		resp = "ERR_VERSION " + strconv.Itoa(msg.Version)
	case 'M':
		resp = "ERR_CMD_ERR"
	case 'I':
		resp = "ERR_INTERNAL"
	case 'R':
		resp = "ERR_REDIRECT " + strconv.Itoa(msg.Version)
	default:
		fmt.Printf("Unknown response kind '%c'", msg.Kind)
		return false
	}
	resp += "\r\n"
	write([]byte(resp))
	if msg.Kind == 'C' {
		write(msg.Contents)
		write(crlf)
	}
	return err == nil
}

type msg_client struct {

	ID int64
	//abc []string
	Request *fs.Msg
	//LeaderId int

}


func serve(conn *net.TCPConn,clientHandlerid int64,ch chan *fs.Msg,raft Node) {
	var a msg_client
	reader := bufio.NewReader(conn)
	for {
		msg, msgerr, fatalerr := fs.GetMsg(reader)
		//fmt.Println("message from client",msg)
		if fatalerr != nil || msgerr != nil {
			reply(conn, &fs.Msg{Kind: 'M'})
			conn.Close()
			break
		}

		if msgerr != nil {
			if (!reply(conn, &fs.Msg{Kind: 'M'})) {
				conn.Close()
				break
			}
		}
		a = msg_client{ID:clientHandlerid,Request:msg}
		//fmt.Println("message before sending to raft node ",a.ID,*(a.Request))
		b, err := json.Marshal(a)
		if err != nil {
			//fmt.Println("Error in encoding message")
		}
		//fmt.Println("after encoding through json",string(b))
		raft.Append(b)
		
			d1 := <- ch
		   //fmt.Println(d1)
		if !reply(conn, d1) {
			conn.Close()
			break
		}
	}
}
func File_Map(raft Node, f *fs.FS){
    //fmt.Println(raft)
      for {
			d1 := <-raft.CommitChannel()
			mm := strings.Split(d1.error," ")
			//fmt.Println("Leader Id received ",mm[4])
			if mm[1]== "Redirect"{
				Leader_ID,_:= strconv.Atoi(mm[4])
				//fmt.Println("Leader Id",Leader_ID)
				response := &fs.Msg{
					Kind:     'R' ,
					Version:  Leader_ID,  
		           	}
		      var e1 msg_client
			 if err := json.Unmarshal(d1.Data, &e1); err != nil {
             panic(err)
               }    	
		       ClientCh := f.FSMAP[e1.ID]    	
		       ClientCh <- response    	
			}else{
			
			//fmt.Println("message from commit channel of Raft Node",d1)
			var e msg_client
			 if err := json.Unmarshal(d1.Data, &e); err != nil {
             panic(err)
            }
            //fmt.Println("Decoded Message",e.ID,e.Request,e.LeaderId)
			response1 := f.ProcessMsg(e.Request)
			ClientCh := f.FSMAP[e.ID]
			ClientCh <- response1
		 }
		}
}
func serverMain(raftnode Node, clientport int) {
	var Handler_ID  int64
	Handler_ID = 2
	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", clientport))
	check(err)
	tcp_acceptor, err := net.ListenTCP("tcp", tcpaddr)
	check(err)
	var f = &fs.FS{Dir: make(map[string]*fs.FileInfo, 1000)}
	 f.FSMAP = make(map[int64](chan *fs.Msg), 1000)
	go File_Map(raftnode,f)
	for {
		
		tcp_conn, err := tcp_acceptor.AcceptTCP()
		//fmt.Println("Inside tcp")
		check(err)
		ClientHandlerCh := make(chan *fs.Msg,1)
		f.FSMAP[Handler_ID] = ClientHandlerCh
		go serve(tcp_conn,Handler_ID,ClientHandlerCh,raftnode)   //Client Handler
		Handler_ID = Handler_ID + 1
	}
}

func main() {

	var RaftNodeObj Node
    
	Nodedetails := os.Args
	id,_:= strconv.Atoi(Nodedetails[1])
	port,_ := strconv.Atoi(Nodedetails[2])
	RaftNodeObj, _ = New1(id, "cluster_test_config1.json")
	//fmt.Println(RaftNodeObj)
	serverMain(RaftNodeObj,port)
	
}




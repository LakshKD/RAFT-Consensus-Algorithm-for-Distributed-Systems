package main

import (  
    "fmt"
    "net"
    "os"
    //"io"
    "io/ioutil"
   //"strconv"
     "strings" 
    //"bytes"
     "strconv"	
)

const (  
    HOST = ""
    PORT = "8089"
    TYPE = "tcp"
    lb = "<"
    rb = ">"  
)

type Filestats struct {
     version int64
     numofbytes int
     exptime int	

}

var files map[string]Filestats

func serverMain(){
  files = make(map[string]Filestats)
  files["test1.txt"] = Filestats {
     0,6,0,
  }
  files["test2.txt"] = Filestats {
     0,8,0,
  }
  files["test3.txt"] = Filestats {
     0,10,0,
  }
  files["test4.txt"] = Filestats {
     0,4,0,
 }
   p, err := net.Listen(TYPE, ":"+PORT)   // Listen for incoming connections.
    if err != nil {
        fmt.Println("Error on listening:", err.Error())
        os.Exit(1)
    }
    
    defer p.Close()  // Close the listener when the application closes.
    fmt.Println("Listening on " + HOST + ":" + PORT)
 for {
        
        conn, err := p.Accept()  // Listen for an incoming connection.
        if err != nil {
            fmt.Println("Error accepting connections ", err.Error())
            os.Exit(1)
        }

        go RequestHandler(conn)// Handle request
	
    }
}


func RequestHandler(conn net.Conn) {  
  readbuf := make([]byte, 1024)   // Make a buffer to hold incoming data.
  size, err := conn.Read(readbuf)  // Read the incoming connection into the buffer.
  if err != nil {
    fmt.Println("Error in reading from the read buffer:", err.Error())
  }
   s := string(readbuf[:size])    //Converting into the string
   p := strings.Fields(s)        //Splitting the string s into the fields in order to extract the Command

  if p[0] == "read"{
   data, err := ioutil.ReadFile(p[1])   //Reading the file,the filename is in the p[1]
     if err != nil {
         
            conn.Write([]byte("ERR_FILE_NOT_FOUND\n"))    //Sending the error message on the channel to the client
           os.Exit(1)
     } 
    conn.Write([]byte("CONTENTS "))
    conn.Write([]byte(lb + strconv.FormatInt(files[p[1]].version,10) + rb + " ")) 
    conn.Write([]byte(lb + strconv.Itoa(files[p[1]].numofbytes) + rb + " "))
    conn.Write([]byte(lb + strconv.Itoa(files[p[1]].exptime) + rb + "\\r\\n")) 
    conn.Write([]byte("\n")) 
    conn.Write([]byte(data)) // Write the file content in the connection channel.
    conn.Write([]byte("\\r\\n"))
     }else if p[0] == "write"{
          
         if _, err := os.Stat(p[1]); err == nil {
	 z := strings.Split(s,"\\r\\n")
         ioutil.WriteFile(p[1],[]byte(z[1]),0644)
	 conn.Write([]byte("OK <version>\r\n"))     
         } else {
          conn.Write([]byte("ERR_FILE_NOT_FOUND\n"))
	    os.Exit(1)
          } 	   
     }else if p[0] == "cas"{


     }else if p[0] == "delete"{

  
     }
  // Close the connection when you're done with it.
  conn.Close()


}

func main() {  
   serverMain()    //Calling FileServer 
}



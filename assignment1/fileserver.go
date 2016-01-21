package main

import (  
    "fmt"
    "net"
    "os"
    "io/ioutil"
    "strings" 
    "strconv"
    "sync"
    "time"
    //"sync/atomic"	
)

const (  
    HOST = "127.0.0.1"
    PORT = "8080"
    TYPE = "tcp"  
)
type Filestats struct {
     version int64
     numofbytes int
     exptime int64
     ctime int64	
}
var files map[string]Filestats

func FileExpiryhandler(m map[string]Filestats) (){
    var mutex = &sync.Mutex{}
   for k, _ := range m {
     if(m[k].exptime!=0){
                                                                          //Delete file when
          if((m[k].ctime+m[k].exptime)<time.Now().Unix()){  
                mutex.Lock()                                                          //creation time plus expiry time is less than current time
                os.Remove(k)
                delete(m,k) 
                mutex.Unlock() 
          }
     }

   }   

}


 
func serverMain(){

     
   
   files = make(map[string]Filestats)
   files["test1.txt"] = Filestats {
     0,6,0,time.Now().Unix(),
  }
  files["test2.txt"] = Filestats {
     0,8,0,time.Now().Unix(),
  }
  files["test3.txt"] = Filestats {
     0,10,0,time.Now().Unix(),
  }
  files["test4.txt"] = Filestats {
     0,4,0,time.Now().Unix(),
 }
   
   
   p, err := net.Listen(TYPE, ":"+PORT)   // Listen for incoming connections.
    if err != nil {
        fmt.Println("Error on listening:", err.Error())
        //os.Exit(1)
    }
    
    defer p.Close()  // Close the listener when the application closes.
    fmt.Println("Listening on " + HOST + ":" + PORT)
 for {
        
        conn, err := p.Accept()  // Listen for an incoming connection.
        if err != nil {
            fmt.Println("Error accepting connections ", err.Error())
            //os.Exit(1)
        }

         go RequestHandler(conn)// Handle request
	       go FileExpiryhandler(files)
    }
}


func RequestHandler(conn net.Conn) { 
  var mutex = &sync.Mutex{}
  //mux sync.Mutex{};
  for {    
  readbuf := make([]byte, 2048)
  size, err := conn.Read(readbuf)  
  if err != nil {
    fmt.Println("Error in reading from the read buffer:", err.Error())
    conn.Close()
    return 
  }
   s := string(readbuf[:size])    
      //fmt.Println(s)                           
   //j := strings.Split(s,"\r\n") 
     j := strings.SplitN(s,"\r\n",2)
     //fmt.Println(j)
     j[1] = strings.Trim(j[1], "\r\n")
      //fmt.Println(j[1])
   p := strings.Split(j[0]," ")	
      //fmt.Println(p)
   if p[0] == "read"{

    cmdlen := len(p)
    //fmt.Println(cmdlen)
    if cmdlen==2{
          //fmt.Println("Inside read")
         
           mutex.Lock()
           data, err := ioutil.ReadFile(strings.TrimSpace(p[1]))
          if err != nil {
         
            conn.Write([]byte("ERR_FILE_NOT_FOUND\n")) 
            conn.Close()   
          } 
          
          // fmt.Println("dsfsdfds")
           
           ver     :=   strconv.FormatInt(files[p[1]].version,10)
           nmbytes :=   strconv.Itoa(files[p[1]].numofbytes)
           exp     :=   strconv.FormatInt(files[p[1]].exptime,10) 
            
           conn.Write([]byte("CONTENTS"+" "+ver+" "+nmbytes+" "+exp+" "+"\r\n"))
           conn.Write([]byte(data))
          //fmt.Println("heyo")
          //fmt.Println(files)
          conn.Write([]byte("\r\n"))
          mutex.Unlock()
        }else{

          conn.Write([]byte("ERR_CMD_ERR\n"))
        }

     }else if p[0] == "write"{
           cmdlen := len(p)
           if cmdlen == 3 || cmdlen == 4 {

                if cmdlen == 3{

                  if _, err := os.Stat(p[1]); err == nil {   //File Exists on the server
	                 //fmt.Println("inside write")
                  mutex.Lock()
                  ioutil.WriteFile(p[1],[]byte(j[1]),0777)    
		              b, _ := strconv.Atoi(p[2])
		              expt := files[p[1]].exptime
                  c1time:= files[p[1]].ctime
                  files[p[1]] = Filestats {
                    //0,b,expt,time.Now().Unix(),
                    0,b,expt,c1time,
                  }       	   
             
                  s := strconv.FormatInt(files[p[1]].version,10)
                  conn.Write([]byte("OK"+" "+s+"\r\n"))    
                  mutex.Unlock()

                  } else { 
                            mutex.Lock()  
                            //fmt.Println("hello")         //File does not Exists on the server
                            ioutil.WriteFile(p[1],[]byte(j[1]),0777)   
						                                       
                            b, _ := strconv.Atoi(p[2])
		                        expt := int64(0)
                            //expt , _ := strconv.ParseInt(p[3],10,64)                 
		           
                    
                            files[p[1]] = Filestats {
                                0,b,expt,time.Now().Unix(),
                            }  
                    
		                       s := strconv.FormatInt(files[p[1]].version,10)     
		                        conn.Write([]byte("OK" +" " +s+"\r\n"))
	                         mutex.Unlock()                                                                     
                          } 
                 }else if cmdlen==4{
                       if _, err := os.Stat(p[1]); err == nil {   //File Exists on the server
                   //fmt.Println("inside write")
                  mutex.Lock()
                  ioutil.WriteFile(p[1],[]byte(j[1]),0777)    
                  b, _ := strconv.Atoi(p[2])
                  expt := files[p[1]].exptime
                  c1time := files[p[1]].ctime
                  files[p[1]] = Filestats {
                    //0,b,expt,time.Now().Unix(),
                    0,b,expt,c1time,
                  }            
             
                  s := strconv.FormatInt(files[p[1]].version,10)
                  conn.Write([]byte("OK"+" "+s+"\r\n"))    
                  mutex.Unlock()

                  } else { 
                            mutex.Lock()  
                            //fmt.Println("hello")         //File does not Exists on the server
                            ioutil.WriteFile(p[1],[]byte(j[1]),0777)   
                                                   
                            b, _ := strconv.Atoi(p[2])
                            
                            expt , _ := strconv.ParseInt(p[3],10,64)                 
                             //expt = expt + time.Now().Unix()
                    
                            files[p[1]] = Filestats {
                                0,b,expt,time.Now().Unix(),
                            }  
                    
                           s := strconv.FormatInt(files[p[1]].version,10)     
                            conn.Write([]byte("OK" +" " +s+"\r\n"))
                           mutex.Unlock()  

                         }

                 }
           }else{
                  conn.Write([]byte("ERR_CMD_ERR\n"))
            }

             
     }else if p[0] == "cas"{

        cmdlen := len(p)
           if cmdlen == 4 || cmdlen == 5 {
          if cmdlen == 4{
            if _, err := os.Stat(p[1]); err == nil {

              s := strconv.FormatInt(files[p[1]].version,10)
              if s == strings.TrimSpace(p[2]){
                   mutex.Lock()
                   ioutil.WriteFile(p[1],[]byte(j[1]),0777)
                    v   := files[p[1]].version
                    b,_ := strconv.Atoi(strings.TrimSpace(p[3]))            //numofbytes
                    e,_ := strconv.ParseInt(strings.TrimSpace(p[4]),10,64)           //exptime
                    delete(files,p[1])
                    v=v+1                                 //Updating the version
                    files[p[1]] = Filestats {            //Creating a new entry for the file in the map
                     v,b,e,time.Now().Unix(),
                  } 
                  s1 := strconv.FormatInt(files[p[1]].version,10)     
                  conn.Write([]byte("OK" +" "+s1+"\r\n"))
                    mutex.Unlock()
               }else {
                   conn.Write([]byte("ERR_VERSION "))
                   conn.Write([]byte(s+"\r\n"))
                   conn.Close()
                   return
               }
          }else{
           conn.Write([]byte("ERR_FILE_NOT_FOUND\n"))    
             conn.Close()
             return     
          }
        }else if cmdlen == 5{
            if _, err := os.Stat(p[1]); err == nil {

              s := strconv.FormatInt(files[p[1]].version,10)
              if s == strings.TrimSpace(p[2]){
                   mutex.Lock()
                   ioutil.WriteFile(p[1],[]byte(j[1]),0777)
                    v   := files[p[1]].version
                    b,_ := strconv.Atoi(strings.TrimSpace(p[3]))            //numofbytes
                    e,_ := strconv.ParseInt(strings.TrimSpace(p[4]),10,64)           //exptime
                    delete(files,p[1])
                    v=v+1                                 //Updating the version
                    files[p[1]] = Filestats {            //Creating a new entry for the file in the map
                     v,b,e,time.Now().Unix(),
                  } 
                  s1 := strconv.FormatInt(files[p[1]].version,10)     
                  conn.Write([]byte("OK" +" "+s1+"\r\n"))
                    mutex.Unlock()
               }else {
                   conn.Write([]byte("ERR_VERSION "))
                   conn.Write([]byte(s+"\r\n"))
                   conn.Close()
                   return
               }
          }else{
           conn.Write([]byte("ERR_FILE_NOT_FOUND\n"))    
             conn.Close()
             return     
          }
           

        }
      }else{
           conn.Write([]byte("ERR_CMD_ERR\n"))    
          }
     
     }else if p[0] == "delete"{

        cmdlen := len(p)
        if cmdlen == 2{ 
       if _, err := os.Stat(strings.TrimSpace(p[1])); err == nil{
             mutex.Lock()
             os.Remove(p[1])
              conn.Write([]byte("OK" + "\r\n")) 
             mutex.Unlock()
         }else{
            conn.Write([]byte("ERR_FILE_NOT_FOUND\n"))
            conn.Close()
            return
         } 
       }else{
        conn.Write([]byte("ERR_CMD_ERR\n"))
       }

     }else {

        conn.Write([]byte("ERR_CMD_ERR\n"))

     }
   }


}
func main() {  
  serverMain()    
}



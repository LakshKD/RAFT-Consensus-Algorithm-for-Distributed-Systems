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
)

// Simple serial check of getting and setting
func TestTCPSimple(t *testing.T) {
	var mutex = &sync.Mutex{}
        go serverMain()
          var wg sync.WaitGroup
        time.Sleep(1 * time.Second) // one second is enough time for the server to start
         
         wg.Add(5)
 
         for i:=1 ; i<=5 ; i++ {
          
         go func (j int) {
        //fmt.Println("dfjnsdjfn")
        name     := "hi.txt"
	contents := "Hey"
	exptime  := 300000
	conn, err:= net.Dial("tcp", "localhost:8080")

		if err != nil {
			t.Error(err.Error()) // report error through testing framework
		}
        //fmt.Println(conn)
	scanner  := bufio.NewScanner(conn)

	// Writing to a file
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents+strconv.Itoa(j))
	scanner.Scan() // read first line
	resp := scanner.Text() // extract the text from the buffer
	arr := strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	version, err := strconv.ParseInt(arr[1], 10, 64) // parse version as number
	if err != nil {
		t.Error("Non-numeric version found")
	}
	


        // try a read now
	fmt.Fprintf(conn, "read %v\r\n", name) 
	scanner.Scan()

	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "CONTENTS")
	expect(t, arr[1], fmt.Sprintf("%v", version)) // expect only accepts strings, convert int version to string
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))	
	scanner.Scan()
	//fmt.Println(scanner.Text())
	expect(t, contents, scanner.Text())

        //Compare And Swap Testing
       mutex.Lock()
       contents = "New Contents"
       version1 := 0
       fmt.Fprintf(conn, "cas %v %v %v %v\r\n%v\r\n", name,version1,len(contents), exptime, contents+strconv.Itoa(j))
       scanner.Scan()
       resp = scanner.Text()
       fmt.Println(resp)
       arr = strings.Split(resp, " ")
       expect(t, arr[0], "OK")
	 version, err = strconv.ParseInt(arr[1], 10, 64) // parse version as number
	 if err != nil {
		t.Error("Non-numeric version found")
	 }
         time.Sleep(2 * time.Second)
    mutex.Unlock()
     // Delete the file
        
	/* fmt.Fprintf(conn, "delete %v\r\n", name) 
	 scanner.Scan()
	 resp = scanner.Text()
	 //fmt.Println(resp)
     arr = strings.Split(resp, " ")
     expect(t, arr[0], "OK")*/

      wg.Done()
        }(i)
     
     }
    wg.Wait()
   //time.Sleep(2 * time.Second)
  
}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

Name    :- Lakshya Kumar
Roll No :- 153050051

Description of the FileServer Written in Go language

-> The below mentioned operations can be performed on the file server :-

 	1. read   (reading the contents of the file)
 	2. write  (Writing into the file )
 	3. cas    (Compare and swap)
 	4. delete (Delete the file from the file server)


-> All the files are stored in the file server's secondary storage and there is a Map data structure that is used to index the files. The structure of the map is as follows:-

   type Filestats struct {
     version int64
     numofbytes int
     exptime int64
     ctime int64	
   }
  var files map[string]Filestats

-> So each file name has 4 fields that are :-

	1.version.
	2.numofbytes(size of the file).
	3.exptime(time at which the file get expired).
	4.ctime (the creation time of the file).

-> Inside serverMain() function there are two functions in the fileserver.go file :-
  1. RequestHandler()     :-   Handle request of each client
  2. FileExpiryhandler()  :-   Time to time deletes the files form the server's secondary storage whose 
                                   creation time + expiry time < current time

Note :- Inside the RequestHandler all the four commands i.e read,write,cas and delete are implemented. 

-> Concurrency
     *In order to handle the concurrency i have used the sync package that will provide the lock() and unlock() function. Whenever there is some updations peformed on the file then the client will acquire the locks.

     *Whenever the updations are performed on the map then also the clients will acquire the locks().
   (Tested the server on 5 multiple clients simultaneously)
      *I have used wait group to generate 5 clients that will try to access the server at the same time and read and writes are performed successfully and when the 5 clients try to do the cas at the same time only one of them will succeed and finally 4 have ERR_VERSION error and only one of them will succeed.

* I have applied the mutex locks on the entire map so whenever any of the operation is performed then the clients will lock the entire map.

Assumptions:-

-> Map is always present in the primary memory and the server is always up. The changes which we are making in the map  are not persistent so when the server goes down then the changes done to the map will get lost. 

-> Initially, there are 4 files present on the server that are test1.txt,test2.txt,test3.txt and test4.txt.
   There entries are also present in the map.

Assumptions on file expiry time:-

In case of write command 
   -> When the write command is like "write filename numbytes \r\n content bytes\r\n"
       
     case-1(File exists on the server)
        *> Then directly write to the file and the expiry time and creation time will remain same.

     case-2(File does not exists on the server)
        *> Then file will be created with 0 expiry time and creation time will be the current time of the server.  

   -> When the write command is like "write filename numbytes exptime\r\n content bytes\r\n"   

     case-1(File exists on the server)   
        *> Then directly write to the file and the expiry time and creation time will remain same.

     case-1(File does not exists on the server)   
        *> Then file will be created with expiry time as "exptime" and creation time will be the current time of the server. This file will gets deleted when the creation time (ctime) + exptime of this file is less than the current time of the server. 

*Testing the Server:-

 1. I have tested the server with the file basic_test.go.
 2. I have checked the file server with "go test" then OK is coming.
 3. I have checked the file server with "go test -race" then there is one Race condition is coming. 
 
*Bugs
 1.When i have tested the server with the basic_test.go file then sometimes it is allowing two clients to enter into the cas function at the file server but sometimes it is depicting the correct behaviour.

 Reference:-

 1. https://tour.golang.org/moretypes/16
 2. https://golang.org/pkg/strings/#example_LastIndex
 3. https://gobyexample.com/epoch
 4. https://blog.golang.org/go-maps-in-action
 5. https://gobyexample.com/maps
 6. http://loige.co/simple-echo-server-written-in-go-dockerized/
 7. http://stackoverflow.com/questions/19208725/example-for-sync-waitgroup-correct

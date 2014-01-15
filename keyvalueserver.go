package main

import (
        "bufio"
        "fmt"
        "io"
        "net"
        "strings"
	"os"
	"log"
	"io/ioutil"
	"os/exec"
)

var (
	CACHE = map[string]string{}
	count = 0 
)

var (
    TRACE   *log.Logger
    INFO    *log.Logger
    WARNING *log.Logger
    ERROR   *log.Logger
)

//logging function
//http://www.goinggo.net/2013/11/using-log-package-in-go.html
func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    TRACE = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    INFO = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    WARNING = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    ERROR = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

//https://github.com/grahamking/Key-Value-Polyglot/blob/master/memg.go
func main() {
    	
	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		ERROR.Println(err)
	}
	defer logfile.Close()
	Init(ioutil.Discard, logfile, os.Stdout, os.Stderr)


	_, err = exec.Command("./createfile.py").Output()
    	if err != nil {
    		ERROR.Println(err)
    	}
    	


        CACHE = make(map[string]string)

        service := ":1204"
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)


	readfile()

	for {
		go stopServer()		

		conn, err := listener.Accept()
		
		if err != nil {
			continue
		}
		INFO.Println("Client "+conn.RemoteAddr().String()+" is connected.....")
		// run as a goroutine
		go handleConn(conn)
	}

}

//stop server on type of "EXIT"
func stopServer(){
	read := bufio.NewReader(os.Stdin)
	input, _ := read.ReadString('\n')
	inputstring:=string([]byte(input)[0:])		
	inputstring = inputstring[:len(inputstring)-1]
			
	if inputstring=="EXIT" {
		writefile()	
		INFO.Println("Server is shutting down...All client connections are closed.")	
		os.Exit(1)
	}
}


//read key values from file and insert in map variable
//http://rosettacode.org/wiki/Read_a_file_line_by_line#Go
func readfile(){
	f, err := os.OpenFile("input.txt", os.O_CREATE|os.O_RDONLY,0600)
    	if err != nil {
		ERROR.Println(err)        	
	    }
	bf := bufio.NewReader(f)
	for {
        	switch line, err := bf.ReadString('\n'); err {
        		case nil:
        		    	// valid line, echo it.  note that line contains trailing \n.
				if line!="\n" {
					parts := strings.Split(line, " ")
				        key := parts[0]
					value := ""
					for i:=1;i<len(parts);i++ {
						if i==(len(parts)-1) {
							value+=parts[i]
						} else {
							value+=parts[i]+" "
						}
					}
	        	        	CACHE[key] = string(value[:len(value)-1])
					count=count+1	
	     			}
			case io.EOF:
        		  	if line > "" {
 	       		        	// last line of file missing \n, but still valid
        		        	fmt.Println(line)
        		    	}
        		    	return
        		
			default:
        		    ERROR.Println(err)
        	}
    	}
}


func handleConn(conn net.Conn) {
	defer closeConnection(conn)
        reader := bufio.NewReader(conn)

        for {
                content, err := reader.ReadString('\n')
                if err == io.EOF {
                        break
                } else if err != nil {
                        ERROR.Println(err)
                        return
                }

                content = content[:len(content)-2] // Chop \r\n

                parts := strings.Split(content, " ")
		
                cmd := parts[0]
                switch cmd {

                case "get":
                        key := parts[1]
                        val, ok := CACHE[key]
                        if ok {
                                conn.Write([]uint8("VALUE ----> " + val + "\r\n"))
				INFO.Println("Client "+conn.RemoteAddr().String()+" Get ---> Key: "+key+" Value: "+string(val))
                        }else{
				conn.Write([]uint8("KEY NOT PRESENT\r\n"))	
				INFO.Println("Client "+conn.RemoteAddr().String()+" Get ---> Key: "+key+" KEY NOT PRESENT ")
			}
                        //conn.Write([]uint8("END\r\n"))

                case "set":
                        key := parts[1]
			value := ""
			for i:=2;i<len(parts);i++ {
				if i==(len(parts)-1) {
					value+=parts[i]
				} else {
					value+=parts[i]+" "
				}
			}
			
			

			_, ok := CACHE[key]
                        if ok {
				oldval:=CACHE[key]
				if oldval==value {
					conn.Write([]uint8("KEY ALREADY PRESENT \r\n"))				
				} else {
					CACHE[key] = string(value)
					conn.Write([]uint8("KEY UPDATED \r\n"))
				}
			
			} else {
				CACHE[key] = string(value)	
				conn.Write([]uint8("STORED\r\n"))
				INFO.Println("Client "+conn.RemoteAddr().String()+" Set ---> Key: "+key+" Value: "+string(value))
			}
			if len(CACHE)>count {
				appendfile(string(key),string(value))
			}else {
				//writefile()
			}
                        
		
		case "delete":
			key := parts[1]
			_, ok := CACHE[key]
                        if ok {
                                delete(CACHE, key)
				INFO.Println("Client "+conn.RemoteAddr().String()+" Delete ---> Key: "+key)
				conn.Write([]uint8("DELETED\r\n"))
                        }else{
				conn.Write([]uint8("KEY NOT PRESENT\r\n"))	
				INFO.Println("Client "+conn.RemoteAddr().String()+" Delete ---> Key: "+key+" KEY NOT PRESENT ")
			}

			//writefile()
		case "list":
			conn.Write([]uint8("-------------------------------\r\n"))	
			conn.Write([]uint8("List of Keys\r\n"))	
			conn.Write([]uint8("-------------------------------\r\n"))	
			for key, _ := range CACHE {
				conn.Write([]uint8(key+"\r\n"))	
			}			
			conn.Write([]uint8("-------------------------------\r\n"))

		case "rename":
			key := parts[1]
			value, ok := CACHE[key]
                        if ok {
                                delete(CACHE, key)
				_, check := CACHE[parts[2]]
				if check {
					conn.Write([]uint8("NEWKEY FOR RENAME ALREADY PRESENT \r\n"))	
				} else {
					CACHE[parts[2]] = string(value)
					INFO.Println("Client "+conn.RemoteAddr().String()+" Rename ---> OldKey: "+key+ " NewKey: "+parts[2])
					conn.Write([]uint8("RENAMED\r\n"))
				}
                        }else{
				conn.Write([]uint8("KEY NOT PRESENT\r\n"))	
				INFO.Println("Client "+conn.RemoteAddr().String()+" Rename ---> Key: "+key+" KEY NOT PRESENT ")
			}


		case "help":
			conn.Write([]uint8("-----------------------------------------------------\r\n"))	
			conn.Write([]uint8("Following Commands are used\r\n"))	
			conn.Write([]uint8("-----------------------------------------------------\r\n"))	
			conn.Write([]uint8("Store/Update key value --> set <key> <value>\r\n"))
			conn.Write([]uint8("Retrive value --> get <key>\r\n"))
			conn.Write([]uint8("Delete key value --> delete <key>\r\n"))
			conn.Write([]uint8("Rename key --> rename <oldkey> <newkey>\r\n"))
			conn.Write([]uint8("List all keys --> list\r\n"))
			conn.Write([]uint8("-------------------------------\r\n\n\n"))			
		
		default:
			conn.Write([]uint8("NOT VALID OPTION\r\n"))
			INFO.Println("Client "+conn.RemoteAddr().String()+" NOT VALID OPTION")	
                }
        }
}

//client close
func closeConnection(conn net.Conn){
	INFO.Println("Client "+conn.RemoteAddr().String()+" is disconnected.....")
	conn.Close()
}

//write key values to file
func writefile(){
	os.Remove("input.txt")
	f, err := os.OpenFile("input.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
    		ERROR.Println(err)
	}

	defer f.Close()
	for key, value := range CACHE {
		val := strings.Contains(value, "\n")
		if val==true {
			fmt.Print(value)
		}
		if _, err = f.WriteString(key+" "+value+"\n"); err != nil {
    			ERROR.Println(err)
		}else{
			count=count+1
		}
	}	
}


//append new key value pair in file 
//http://stackoverflow.com/questions/7151261/append-to-a-file-in-go
func appendfile(ke string,va string){
	f, err := os.OpenFile("input.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
    		ERROR.Println(err)
	}

	defer f.Close()

	if _, err = f.WriteString(ke+" "+va+"\n"); err != nil {
    		ERROR.Println(err)
	}else{
		count=count+1
	}
}

//check errors
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		ERROR.Println(err.Error())
		os.Exit(1)
	}
}


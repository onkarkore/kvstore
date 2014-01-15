
The Key-Value store

----------------------------------------------------------------------------------------------------------------------------------------

1. Start the server by running following command
	go run keyvalueserver.go


2. Then client can connect to the server by using following command
	telnet <ip address of server> <port number>


3. Multiple clients can allow to connect by using same command


4. Server is stopped by typing 'EXIT' on server side


5. Operations performed on client side
	- Type 'help' and press enter to see all commands.
	- Store key-value pair by using follwing command (space is not allowed in key)
		set <key> <value>
	- Update key-value pair by using follwing command
		set <key> <value>
	- Retrive value by using follwing command
		get <key>
	- Delete key-value pair by using follwing command
		delete <key>
	- Rename key by using follwing command
		rename <oldkey> <newkey>
	- List all keys by using follwing command
		list

6. Operation performed on server side
	- Type 'EXIT' to stop server and close all connections with client
	- Main function creates log file "log.txt"
	- Logging is done at every step to recover key-value pairs after failure (failure like if we updated or deleted key-value pair 		  then it will only affect map but "input.txt" contain old data. At the same time if server crashed or power failure cause 		  server to stop abnormally then our data become inconsistent.)
	- Log file contains information like
		- date and time when client connect and disconnect
		- operations performed by clients with date and time	
		- connection errors, io errors are also logged	
	- Read "input.txt" file and load all key-value pairs in map when server starts
	- Execute "createfile.py" python script to check whether server stoped normally by pressing 'EXIT' command last time. If did 		  not stopped normally then "input.txt" file is recovered by using "log.txt" file.
	- If new key-value pair is added then write immediately to "input.txt" file
	- If key-value pair is updated or deleted then only map will updated. The updated key-value pair in map will be stored in file 		  "input.txt" after pressing 'EXIT'
	

References -
1. Logging - http://www.goinggo.net/2013/11/using-log-package-in-go.html
2. Main idea - https://github.com/grahamking/Key-Value-Polyglot/blob/master/memg.go
3. Read file - http://rosettacode.org/wiki/Read_a_file_line_by_line#Go
4. Append data in file - http://stackoverflow.com/questions/7151261/append-to-a-file-in-go


















# Go Chat 

## Required GO v1.11+

Go implementation CLI chat

## Features
- Both client and server should be console apps. Users send messages to chat via stdin prompt
- Users must specify their names when joining chat
- Each message is broadcasted to all chat members
- Same user (identified by same name) can have multiple simultaneous clients running.
- Online status calculation: user is "online" when he/she has at least one client running, otherwise user is "offline".
- Server must notify all chat members when some user comes online or goes offline.
- Client and server part should be written using Golang
- Network protocol must be extensible. It should be possible to add new features in future.
- Clean, readable code. Simplicity is a plus
- All errors must be handled, some examples:
    - network failures
    - slow clients
    - invalid messages (like attempt to send message on behalf of other user)
    - etc
    
## HowTo 

### Run server
Run command from project home directory `go run cmd/server/main.go`. It will download required libraries and start server 
by default at `localhost:8080`. For custom port specify port with flag `--port`. (Example `go run cmd/server/main.go --port 3443`)

### Run client 
Run command from project home directory `go run cmd/client/main.go`. It will download required libraries and start CLI client. After start 
it asks Username and after that connects to server (by default `localhost:8080/api/ws`). To change `hostname`, `port` or `endpoint` You 
can use flags: `--hostname`, `--port` and `--path` respectively. (Example `go run cmd/client/main.go --hostname localhost --port 3443 --path /api/ws`)

## For better experience use terminal that supports ASCII Control Characters. 


  


  
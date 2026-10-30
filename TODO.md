
# TODO

# /!\ !!! LEARN GO LANGUAGE / PROJECT MAIN CONCEPTS !!! /!\

## All Tasks

### Fixs
- Add topic name to CmdMessageDatas
- Update IpcIdentifier to client IpcPath (check in ToServer.Path)
- Check error codes usage (1, 2 and 3) and priorities (1 > 3 > 2)
- Update client.IsActive value depending on connection state

### Feats
- Add messages replay on subscribers if needed
- Add subscriber go routine receiving messages from topic go routine then send it to client
- Implement client subscribe logic, waiting for new messages
- Add SIGINT and SIGTERM management
- Check message length in producer mode (1024 bytes)
- Add Mutex on topic messages list
- Add tests on every new global features
- Convert text messages into little endian on stdout if --raw
- Handle Graceful shutdown of subscribers/producers if the server crash ...
- Return new offset after new message (check inside subject)


## Possible repartition

### Part1

#### Fixs
- Update IpcIdentifier to client IpcPath (check in ToServer.Path)
- Update client.IsActive value depending on connection state

#### Feats
- Add subscriber go routine receiving messages from topic go routine then send it to client
- Add messages replay on subscribers if needed
- Check message length in producer mode (1024 bytes)
- Convert text messages into little endian on stdout if --raw
- Return new offset after new message (check inside subject)


### Part 2

#### Fixs
- Add topic name to CmdMessageDatas
- Check error codes usage (1, 2 and 3) and priorities (1 > 3 > 2)

#### Feats
- Implement client subscribe logic, waiting for new messages
- Add Mutex on topic messages list
- Add SIGINT and SIGTERM management
- Handle Graceful shutdown of subscribers/producers if the server crash ...
- Add tests on every new global features

# hw3
## Readout
```
cm_Primary:
Record    Owner_ID  copy_set
0         0         map[]
1         1         map[]
2         2         map[]
3         0         map[]
4         1         map[]
cm_Backup:
Record    Owner_ID  copy_set
0         0         map[]
1         1         map[]
2         2         map[]
3         0         map[]
4         1         map[]
process_0:
Page      isOwner   Access
0         true      ReadWrite
1         false     Nil
2         false     Nil
3         true      ReadWrite
4         false     Nil
process_1:
Page      isOwner   Access
0         false     Nil
1         true      ReadWrite
2         false     Nil
3         false     Nil
4         true      ReadWrite
process_2:
Page      isOwner   Access
0         false     Nil
1         false     Nil
2         true      ReadWrite
3         false     Nil
4         false     Nil
```
Each process and central manager maintains a page table with the same number of rows as there are pages The information they keep is slightly different. CM keeps track of each page's owner and the set of processes that have copies of them. Each process tracks if they own the page and what access rights they have. If they own the page, they have RW access, and if they have a copy, they only have read access.

```
Type             From      To        Page      Requester
ReadRequest      2         3         1         1
ReadForward      3         1         1         1
ReadPage         1         2         1         1
ReadConfirmation 2         3         1         1
Done
```
## Normal Ivy
Implemented in hw3 folder

## Fault Tolerant Ivy

Implemented in hw3_ft folder

Primary sends heartbeat to backup every 450ms (arbitrary number), so if primary dies, backup will change the cm reference that all clients use to itself and continue from there. The Primary and Backup shares the same channel that the processes send to, so that no packet is lost. They also share another heartbeat channel, so that the monitor coroutine does not consume the normal messages.

The Primary just copies all its state to the backup on every change to state
(simulated, since a real life primary isnt restricted by channel structs)

when the primary comes alive, backup detects heartbeat again, and copies the state over. Primary can resume from last state.

## Caveats

The single shared channel and cm reference represents a hidden and infallible single point of failure. It is like a hidden load balancer. To get around it, it would be best to fall back to Paxos like protocol so that processes can timeout and re-request and re-elect CM if necessary. However, Paxos like protocol would probably have higher overhead, which would not be suitable for shared memory performance.


## test cases for normal and fault tolerant Ivy

### 100 read requests for random pages by random processes
Normal:
Run `make r`.

FT:

Run `CGO_ENABLED=1 go run -race cm.go main.go mailbox.go message.go page.go process.go records.go -r={number of read requests}`.

### 100 write requests for random pages by random processes
Normal:
Run `make w`.

FT:

Run `CGO_ENABLED=1 go run -race cm.go main.go mailbox.go message.go page.go process.go records.go -w={number of write requests}`.

You can use the provided makefile to add specify the cli args.


## testing performance
all testing is done for 100 requests, by random processes looking for random pages. (chosen by random choice)

normal ivy takes under 10 ms to read/write 100 pages
### FT performance

At 0% crashes, ft sometimes takes 10-20 ms but some times takes 100ms for read/write.
At >0% chance to crash and recover, ft ivy takes 500ms for 100 requests whether they are read or write.

Clearly, normal ivy is more performant, and when ft kicks in, performance reduces by order of 100.


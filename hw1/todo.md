# Assignment 1 TODO
- refactor output messages
- add tabwriter
- refactor helper functions to another file/package
- decide clock adjustment output
- 

## soln
- use unbuffered channels of type Msg as means of communication. 
- server has 

## 1_2

- reciever adjusts clock according to incoming message
- print to file instead of terminal

- saves messages that clients receive to some object
- we can sort the object.

## 1_3
- currently printf shows msg contents, but we also want to know if its Receive, Transmit or Forward and who did it [also](#1_2)
- how to sort vector clock?
    - AdjustClock() needs to be changed Done
- causality violation detection
- Total order even for concurrent ts

## 2_1
- goroutine may not be suitable for timeout waiting
- does coordinator need to check incoming messages for new coordinator? yes.
    - channel doesnt need mutex i think

- channel needs to send in a goroutine or a dead node will cause the broadcast to block (done)
- how to wait for timeouts? (done)


## troubleshooting
- bully worst and bully normal
    - symptoms: some nodes dont seeme to respond to victory message, elect themselves and nobody responds to them 
    - node 9 sends the victory message so early that 
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
- causality violation detection
# Assignment 1

Consider some client-server architecture as follows. Several clients are registered
to the server. Periodically, each client sends message to the server. Upon receiving a
message, the server flips a coin and decides to either forward the message to all other
registered clients (excluding the original sender of the message) or drops the message
altogether. To solve this question, you will do the following:
1. Simulate the behaviour of both the server and the registered clients via GO routines. 

## soln
- use unbuffered channel as means of communication. 
- By default, sends and receives block until the other side is ready. This allows goroutines to synchronize without explicit locks or condition variables.
- value: client_id, msg
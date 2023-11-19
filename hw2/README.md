# 50.041 Distributed Systems Homework Assignment 2
1005285 Joshua Ng Tze Wee 

## Requirements
- Go 1.21.1

## 2_1_1 Lamport’s Shared Priority Queue (Protocol)

### How to run
```
cd 2_1_1/
make
```
- Runs all nodes (default is 10)

### Syntax
```
client_id->-1 @LamportClock: Data                // sending message syntax
adjust clock_some_id: LamportClock->LamportClock // adjust clock
s drop: '6->-1 @44: 7462'                        // server drop message

From To Timestamp Data
id   id ts        Data
```

## 2_1_2 Ricart and Agrawala Lamport’s Shared Priority Queue (Protocol)
# 50.041 Distributed Systems Homework Assignment
1005285 Joshua Ng Tze Wee 

## Requirements
- Go 1.21
## 1_1

### How to run
```
go run hw1/1_1/hw1.go
```
- Runs clients and server

### How to read output
#### Syntax
```
self.id [T]ransmit/[R]ecieve msg.id msg.data 
-1      Tx                   2        898
2       Rx                   -1       898
```
#### Location
- stdout
#### Notes
- Main function creates array of channels type Msg, 

## 1_2

### How to run
```
go run hw1/1_2/hw1_2.go
```
- Runs clients and server

### How to read output
#### Syntax
```
self.id [T]ransmit/[R]ecieve msg.id msg.data msg.clock
-1      Tx                   2      898      15
2       Rx                   -1     898      15
2       Ad                   nil    nil      15->15
2       Rx                   -1     898      15
```
#### Location
- stdout
#### Notes
- Main function creates array of channels type Msg, 
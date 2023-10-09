# 50.041 Distributed Systems Homework Assignment
1005285 Joshua Ng Tze Wee 

## Requirements
- Go 1.21.1
## 1_1

### How to run
```
cd 1_1/
make runs
```
- Runs clients and server

### Syntax
```
self.id [T]ransmit/[R]ecieve msg.id msg.data 
-1      Tx                   2        898
2       Rx                   -1       898
```
### Location
- stdout

## 1_2

### How to run
```
cd 1_2/
make runs
```
- Runs clients and server

### Syntax
```
self.id [T]ransmit/[R]ecieve msg.id msg.data msg.clock
-1      Tx                   2      898      15
2       Rx                   -1     898      15
2       Ad                   nil    nil      15->15
2       Rx                   -1     898      15
```
### Location
- stdout

## 1_3

### How to run
```
cd 1_3/
make runs
```
- Runs clients and server

### Syntax
```
From To Timestamp Data
1    -1 2         1750
...
```
### Location
- stdout

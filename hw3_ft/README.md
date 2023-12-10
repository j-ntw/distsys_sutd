# hw3
## Readout
```
cm3:
Record    Owner_ID  isLocked  copy_set
0         0         false     map[]
1         1         false     map[]
2         2         false     map[]
3         0         false     map[]
4         1         false     map[]
process_0:
Page      isOwner   isLocked  Access
0         true      false     ReadWrite
1         false     false     Nil
2         false     false     Nil
3         true      false     ReadWrite
4         false     false     Nil
process_1:
Page      isOwner   isLocked  Access
0         false     false     Nil
1         true      false     ReadWrite
2         false     false     Nil
3         false     false     Nil
4         true      false     ReadWrite
process_2:
Page      isOwner   isLocked  Access
0         false     false     Nil
1         false     false     Nil
2         true      false     ReadWrite
3         false     false     Nil
4         false     false     Nil
...
```
Each process and central manager maintains a page table with the same number of rows as there are pages The information they keep is slightly different. CM keeps track of each page's owner and the set of processes that have copies of them. Each process tracks if they own the page and what access rights they have. If they own the page, they have RW access, and if they have acopy, they only have read access.
## Fault Tolerant Ivy
we slow down CM processing to 1 second, and primary sends  heart beat to backup  0,5s, so if primry dies, backup will change the cm reference that all clients use to itself and continue from there

(simulated, its cheating a bit since im basically having an invisible, infallible load balancer that changes internal ip to process facing ip addresses)
and the primary just copies all its state to the backup on every change to state
(simulated, since a real life primary isnt restricted by channel structs)

and the primary just copies all its state to the backup on every change to state
(simulated, since a real life primary isnt restricted by channel structs)


when primary comes alive, backup detects heartbeat again, and copies the state over

Joshua Ng, [12/10/2023 12:05 PM]
is this cheating too much lol
## test cases

### P3 wants to read page 1 (page fault at P3)
Run `make r`.

### P3 wants to write page 1 (page fault at P3)
Run `make w`.


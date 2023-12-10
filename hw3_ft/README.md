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

## Fault Tolerant Ivy
we slow down CM processing to 1 second, and primary sends heartbeat to backup  0.5s, so if primary dies, backup will change the cm reference that all clients use to itself and continue from there

The Primary just copies all its state to the backup on every change to state
(simulated, since a real life primary isnt restricted by channel structs)

when the primary comes alive, backup detects heartbeat again, and copies the state over

## test cases

### P3 wants to read page 1 (page fault at P3)
Run `make r`.

### P3 wants to write page 1 (page fault at P3)
Run `make w`.


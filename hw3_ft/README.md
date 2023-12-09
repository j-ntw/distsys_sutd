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

## test cases

### P3 wants to read page 1 (page fault at P3)
Run `make r`.

### P3 wants to write page 1 (page fault at P3)
Run `make w`.


[ivy ppt](http://www0.cs.ucl.ac.uk/staff/B.Karp/gz03/f2011/lectures/gz03-lecture5-Ivy.pdf)
• RQ (read query, reader to MGR)
• RF (read forward, MGR to owner)
• RD (read data, owner to reader)
• RC (read confirm, reader to MGR)
• WQ (write query, writer to MGR)
• IV (invalidate, MGR to copy_set)
• IC (invalidate confirm, copy_set to MGR)
• WF (write forward, MGR to owner)
• WD (write data, owner to writer)
• WC (write confirm, writer to MGR)

## notes

remove clock, its ok not to be ordered.
is isLocked needed?
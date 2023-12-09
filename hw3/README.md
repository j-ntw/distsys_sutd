# hw3
CM is like a hotel manager. he keeps track of rooms/keys and who checked them out.

CM knows about every key.

## test cases

### P3 wants to read page 1 (page fault at P3)

### P3 wants to write page 1 (page fault at P3)


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
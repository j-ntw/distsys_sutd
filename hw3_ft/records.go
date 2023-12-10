package main

import "sync"

type Records struct {
	records []CM_Record
	sync.RWMutex
}

func (records *Records) Get() []CM_Record {
	records.RLock()
	defer records.RUnlock()
	return records.records
}

func (records *Records) GetOwner(page_no int) int {
	records.RLock()
	defer records.RUnlock()
	return records.records[page_no].owner_id
}

func (records *Records) GetCopySet(page_no int) map[int]bool {
	records.RLock()
	defer records.RUnlock()
	return records.records[page_no].copy_set
}

func (records *Records) IsCopySetEmpty(page_no int) bool {
	records.RLock()
	defer records.RUnlock()
	return len(records.records[page_no].copy_set) == 0
}

func (records *Records) SetRequester(page_no int, requester_id int) {
	records.Lock()
	defer records.Unlock()
	records.records[page_no].copy_set[requester_id] = true
}

func (records *Records) DeleteCopyHolder(page_no int, requester_id int) {
	records.Lock()
	defer records.Unlock()
	delete(records.records[page_no].copy_set, requester_id)
}

func (records *Records) Set(newRecords []CM_Record) {
	records.Lock()
	defer records.Unlock()
	records.records = newRecords
}

func newRecords(records []CM_Record) *Records {
	return &Records{records: records}
}

type CM_Record struct {
	owner_id int
	copy_set map[int]bool // set of process_ids who have read only copies

}

func newRecord(id int) *CM_Record {
	copy_set := make(map[int]bool)
	return &CM_Record{
		owner_id: id,

		copy_set: copy_set,
	}
}

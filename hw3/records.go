package main

import "sync"

type Records struct {
	records []CM_Record
	sync.RWMutex
}

func (records *Records) Get() []CM_Record {
	records.Lock()
	defer records.Unlock()
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
	Copy(newRecords, records.records)

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
func (records *Records) DeepCopy() []CM_Record {
	// locks source records (self) and creates a deep copy
	records.Lock()
	defer records.Unlock()
	src := records.records
	dst := make([]CM_Record, len(src))

	for i, record := range src {
		// Copy owner_id
		dst[i].owner_id = record.owner_id

		// Copy copy_set
		dst[i].copy_set = make(map[int]bool)
		for key, value := range record.copy_set {
			dst[i].copy_set[key] = value
		}
	}

	return dst
}

func Copy(src []CM_Record, dst []CM_Record) {
	// copies the values from source to destination, without creating new objects or changing pointers
	for i, record := range src {
		// Copy owner_id
		dst[i].owner_id = record.owner_id

		// Copy copy_set
		for key, value := range record.copy_set {
			dst[i].copy_set[key] = value
		}
	}
}

package main

import "sync"

type CM_REF struct {
	cm *CM
	sync.RWMutex
}

func (cm_ref *CM_REF) GetRef() *CM {
	cm_ref.RLock()
	defer cm_ref.RUnlock()
	return cm_ref.cm
}

func (cm_ref *CM_REF) SetRef(new_ref *CM) {
	cm_ref.Lock()
	defer cm_ref.Unlock()
	cm_ref.cm = new_ref
}

func newCM_REF(new_ref *CM) *CM_REF {
	return &CM_REF{cm: new_ref}
}

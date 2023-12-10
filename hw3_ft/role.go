package main

import "fmt"

type RoleType int

const (
	Primary RoleType = iota
	Backup
	Unused
)

func (a RoleType) String() string {
	return [...]string{"Primary", "Backup", "Unused"}[a]
}

func copyState(from RoleType, to RoleType) {
	cm_arr[0].Lock()
	defer cm_arr[0].Unlock()
	cm_arr[1].Lock()
	defer cm_arr[1].Unlock()
	// TODO: hacky, non safe
	cm_arr[int(to)].records = cm_arr[int(from)].records
	fmt.Printf("cm%d: copy state to cm%d\n", from+numCM+1, from+numCM+2)
}

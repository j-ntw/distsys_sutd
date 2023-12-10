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

	// TODO: hacky, non safe
	cm_arr[int(to)].records.Set(cm_arr[int(from)].records.Get())
	fmt.Printf("cm_%s: copy state to cm_%s\n", from.String(), to.String())
}

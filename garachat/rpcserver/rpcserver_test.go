package rpcserver

import (
	"strings"
	"testing"
)

// test Create and auth calls
func TestCreateAuth(t *testing.T) {
	svr := InitServer()
	users := []string{"max", "max", "max1", "max2", "max3", " ", ""}
	wantErr := []bool{false, true, false, false, false, true, true}
	for i := 0; i < len(users); i++ {
		usr := users[i]
		args := Args{usr, 0, 0, 0, ""}
		var pw uint64
		err := svr.NewUser(args, &pw)
		switch {
		case i == 1 && err == nil:
			t.Fatal("no error on duplicate user")
		case strings.Trim(usr, "s\t\n\r") == "" && err == nil:
			t.Fatalf("no error on empty username")
		case strings.ContainsAny(usr, "s\t\n\r") && err == nil:
			t.Fatalf("no error on illegal chars")
		case err != nil && !wantErr[i]:
			t.Fatalf("failed to create user:%v", err)
		}
		if wantErr[i] {
			continue
		}
		//test auth
		args.T = pw
		if !svr.auth(args) {
			t.Fail()
		}
	}
}

func TestSubmitRead(t *testing.T) {
	svr := InitServer()
	args := Args{"max", 0, 0, 0, "test"}
	var pw uint64
	svr.NewUser(args, &pw)
	args.T = pw
	var resp string //not really used
	if err := svr.Submit(args, &resp); err != nil {
		t.Fail()
	}
}

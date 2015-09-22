package gobreak

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Pwd  string
}

func getUsers() []User{
	var lst []User
	lst = append(lst, User{"abc", "abc"})
	lst = append(lst, User{"xyz", "xyz"})
	return lst
}

func Test_Find(t *testing.T) {
	lst:=getUsers()
	idx:= Slice(&lst).Find(User{"abc", "abc"})
	fmt.Println(idx)
}

func Test_RemoveAt(t *testing.T) {
	lst:=getUsers()
	Slice(&lst).RemoveAt(1)
	fmt.Println(lst)
}

func Test_Remove(t *testing.T) {
	lst:=getUsers()
	Slice(&lst).Remove(User{"abc", "abc"})
	fmt.Println(lst)
}
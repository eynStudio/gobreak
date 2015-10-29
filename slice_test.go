package gobreak

import (
	"fmt"
	"testing"
)

type User struct {
	Id   GUID
	Name string
	Pwd  string
}
func (p User) ID() GUID { return p.Id }

func getUsers() []User {
	var lst []User
	lst = append(lst, User{"1", "abc", "abc"})
	lst = append(lst, User{"2", "xyz", "xyz"})
	return lst
}

func Test_Find(t *testing.T) {
	lst := getUsers()
	idx := Slice(&lst).Find(User{"1","abc", "abc"})
	fmt.Println(idx)
}

func Test_FindEntity(t *testing.T) {
	lst := getUsers()
	entity := Slice(&lst).FindEntity("1")
	fmt.Println(entity)
}

func Test_ReplaceEntity(t *testing.T) {
	lst := getUsers()
	Slice(&lst).ReplaceEntity(User{"3","333","333"})
	Slice(&lst).ReplaceEntity(User{"1","111","111"})
	fmt.Println(lst)
}
func Test_RemoveEntity(t *testing.T) {
	lst := getUsers()
	Slice(&lst).RemoveEntity("1")
	fmt.Println(lst)
}
func Test_RemoveAt(t *testing.T) {
	lst := getUsers()
	Slice(&lst).RemoveAt(1)
	fmt.Println(lst)
}

func Test_Remove(t *testing.T) {
	lst := getUsers()
	Slice(&lst).Remove(User{"1","abc", "abc"})
	fmt.Println(lst)
}

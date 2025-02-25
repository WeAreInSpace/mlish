package main

import (
	"fmt"
	"strings"

	"github.com/WeAreInSpace/mlish"
)

type TestModel struct {
	name string
	age  int
}

type NewTestModel struct {
	name            string
	age             int
	nameInLowercase string
}

func main() {
	userModal := mlish.NewModel[TestModel]()

	userModal.Add(
		[]TestModel{
			{
				name: "Gorn",
				age:  14,
			},
			{
				name: "Gornny",
				age:  15,
			},
		}...,
	)
	userModal.ForEach(
		func(item *mlish.ForParams[TestModel]) TestModel {
			return item.GetData()
		},
	)
	userModal.Add(
		[]TestModel{
			{
				name: "Aran",
				age:  26,
			},
			{
				name: "Thanachai",
				age:  27,
			},
		}...,
	)

	newUserModal := mlish.And(
		userModal,
		func(item *mlish.ForParams[TestModel]) NewTestModel {
			newTestModel := NewTestModel{}
			newTestModel.age = item.GetData().age
			newTestModel.name = item.GetData().name
			newTestModel.nameInLowercase = strings.ToLower(item.GetData().name)
			return newTestModel
		},
	)
	newUserModal.For(
		func(item *mlish.ForParams[NewTestModel]) {
			fmt.Printf("name %s, age %d, spacialName %s\n", item.GetData().name, item.GetData().age, item.GetData().nameInLowercase)
		},
	)

	newUserModal.Filter(func(item *mlish.ForParams[NewTestModel]) NewTestModel {
		if len(item.GetData().nameInLowercase) > 4 {
			return item.GetData()
		}
		return NewTestModel{}
	})
}

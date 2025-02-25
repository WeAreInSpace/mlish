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
		[]*TestModel{
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
		func(item *mlish.ForParams[TestModel]) *TestModel {
			return item.DataAddr()
		},
	)
	userModal.Add(
		[]*TestModel{
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

	newUserModal := mlish.Migrate(
		userModal,
		func(item *mlish.ForParams[TestModel]) *NewTestModel {
			newTestModel := &NewTestModel{}
			newTestModel.age = item.DataAddr().age
			newTestModel.name = item.DataAddr().name
			newTestModel.nameInLowercase = strings.ToLower(item.DataAddr().name)
			return newTestModel
		},
	)
	newUserModal.For(
		func(item *mlish.ForParams[NewTestModel]) {
			fmt.Printf("name %s, age %d, spacialName %s\n", item.DataAddr().name, item.DataAddr().age, item.DataAddr().nameInLowercase)
		},
	)

	newUserModal = newUserModal.Filter(
		func(item *mlish.ForParams[NewTestModel]) *NewTestModel {
			if len(item.DataAddr().nameInLowercase) > 4 {
				return item.DataAddr()
			}
			return nil
		},
	)
	newUserModal.For(
		func(item *mlish.ForParams[NewTestModel]) {
			fmt.Printf("new: name %s, age %d, spacialName %s\n", item.DataAddr().name, item.DataAddr().age, item.DataAddr().nameInLowercase)
		},
	)

	newUserModal = newUserModal.FilterByRegex(
		"[a-z,A-Z]",
		func(item *mlish.ForParams[NewTestModel]) string {
			return item.DataAddr().name
		},
	)

	newUserModal.For(
		func(item *mlish.ForParams[NewTestModel]) {
			fmt.Println("Regex", item.DataAddr().name)
		},
	)
}

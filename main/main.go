package main

import (
	"os"
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
	mlish.Settings.DebugMode = false
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
			{
				name: "Thanachai",
				age:  27,
			},
		}...,
	)

	userModal.Push(
		os.Stdout,
		func(item *mlish.ForParams[TestModel]) []byte {
			return []byte(item.DataAddr().name + "\n")
		},
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

	newUserModal = newUserModal.Filter(
		func(item *mlish.ForParams[NewTestModel]) *NewTestModel {
			if len(item.DataAddr().nameInLowercase) > 4 {
				return item.DataAddr()
			}
			return nil
		},
	)

	newUserModal = newUserModal.FilterByRegex(
		`^[a-zA-Z\s]+`,
		func(item *mlish.ForParams[NewTestModel]) string {
			return item.DataAddr().nameInLowercase
		},
	)

	newUserModal.Push(
		os.Stdout,
		func(item *mlish.ForParams[NewTestModel]) []byte {
			return []byte(item.Data().nameInLowercase + "\n")
		},
	)
}

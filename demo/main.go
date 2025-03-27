package main

import (
	"fmt"

	"github.com/WeAreInSpace/mlish"
)

type UserModel struct {
	nickName string
	name     string
	age      int8
}

type UserModelWithProfile struct {
	nickName string
	name     string
	age      int8
	Profile
}

type Profile struct {
	description string
	image       string
}

func main() {
	userModel := mlish.NewModel[UserModel]()
	userModel.Add(
		[]*UserModel{
			{
				nickName: "Gorn",
				name:     "Aranthanachai Tangseng",
				age:      15,
			},
			{
				nickName: "John",
				name:     "John Doe",
				age:      18,
			},
			{
				nickName: "Jane",
				name:     "Jane Doe",
				age:      15,
			},
		}...,
	)
	userModel.ForEach(
		func(item *mlish.ForParams[UserModel]) *UserModel {
			newUserModel := &UserModel{
				nickName: item.DataAddr().nickName,
				name:     item.DataAddr().name,
				age:      item.DataAddr().age,
			}
			return newUserModel
		},
	)

	userModelWithProfile := mlish.Migrate(
		userModel,
		func(item *mlish.ForParams[UserModel]) *UserModelWithProfile {
			newItem := &UserModelWithProfile{
				nickName: item.DataAddr().nickName,
				name:     item.DataAddr().name,
				age:      item.DataAddr().age,

				Profile: Profile{
					description: "",
					image:       "",
				},
			}
			return newItem
		},
	)

	userModelWithProfile.For(
		func(item *mlish.ForParams[UserModelWithProfile]) {
			fmt.Printf(
				"Nickname: %s\nName: %s\nAge: %d\nProfile: description: %s, image: %s\n\n",
				item.DataAddr().nickName,
				item.DataAddr().name,
				item.DataAddr().age,
				item.DataAddr().Profile.description,
				item.DataAddr().Profile.image,
			)
		},
	)

	userModel.For(
		func(item *mlish.ForParams[UserModel]) {
			fmt.Printf(
				"Nickname: %s\nName: %s\nAge: %d\n\n",
				item.DataAddr().nickName,
				item.DataAddr().name,
				item.DataAddr().age,
			)
		},
	)
}

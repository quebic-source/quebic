//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"os"
	"quebic-faas/types"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var username string
var firstname string
var password string
var role string

func init() {
	setupUserCmds()
	setupUserCompFlags()
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User commonds",
	Long:  `User commonds`,
}

func setupUserCmds() {
	userCmd.AddCommand(userLoginCmd)
	userCmd.AddCommand(userChangePWDCmd)
	userCmd.AddCommand(currentAuthCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userUpdateCmd)
	userCmd.AddCommand(userGetAllCmd)
}

func setupUserCompFlags() {
	//login
	userLoginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	userLoginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")

	//change-password
	userChangePWDCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")

	//create
	userCreateCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	userCreateCmd.PersistentFlags().StringVarP(&firstname, "firstname", "f", "", "firstname")
	userCreateCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
	userCreateCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "role")

	//update
	userUpdateCmd.PersistentFlags().StringVarP(&firstname, "firstname", "f", "", "firstname")
}

var userLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "user : login",
	Long:  `user : login`,
	Run: func(cmd *cobra.Command, args []string) {
		userLogin(cmd, args)
	},
}

var userChangePWDCmd = &cobra.Command{
	Use:   "change-password",
	Short: "user : change-password",
	Long:  `user : change-password`,
	Run: func(cmd *cobra.Command, args []string) {
		userChangePWD(cmd, args)
	},
}

var currentAuthCmd = &cobra.Command{
	Use:   "auth-context",
	Short: "user : auth-context",
	Long:  `user : auth-context`,
	Run: func(cmd *cobra.Command, args []string) {
		currentAuth(cmd, args)
	},
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "user : create",
	Long:  `user : create`,
	Run: func(cmd *cobra.Command, args []string) {
		userSave(cmd, args, true)
	},
}

var userUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "user : update",
	Long:  `user : update`,
	Run: func(cmd *cobra.Command, args []string) {
		userSave(cmd, args, false)
	},
}

var userGetAllCmd = &cobra.Command{
	Use:   "ls",
	Short: "user : get-all",
	Long:  `user : get-all`,
	Run: func(cmd *cobra.Command, args []string) {
		userGetAll(cmd, args)
	},
}

func userLogin(cmd *cobra.Command, args []string) {
	authDTO := &types.AuthDTO{Username: username, Password: password}

	mgrService := appContainer.GetMgrService()

	authToken, errResponse := mgrService.UserLogin(authDTO)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	appContainer.GetAppConfig().Auth.AuthToken = authToken.Token
	appContainer.SaveConfiguration()

	color.Green("Successfully logged in. Welcome to quebic !!!")
}

func userChangePWD(cmd *cobra.Command, args []string) {
	user := &types.User{Password: password}
	mgrService := appContainer.GetMgrService()

	errResponse := mgrService.UserChangePassword(user)

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("Successfully changed your password")
}

func currentAuth(cmd *cobra.Command, args []string) {
	mgrService := appContainer.GetMgrService()

	user, errResponse := mgrService.UserCurrentAuth()

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	username := user.Username
	firstname := user.Firstname
	role := user.Role

	color.Green("Auth context \nusername : %s, firstname : %s, role : %s", username, firstname, role)
}

func userSave(cmd *cobra.Command, args []string, isCreate bool) {
	user := &types.User{
		Username:  username,
		Password:  password,
		Firstname: firstname,
		Role:      role,
	}
	mgrService := appContainer.GetMgrService()

	var errResponse *types.ErrorResponse
	if isCreate {
		errResponse = mgrService.UserCreate(user)
	} else {
		errResponse = mgrService.UserUpdate(user)
	}

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("Successfully saved user")
}

func userGetAll(cmd *cobra.Command, args []string) {
	mgrService := appContainer.GetMgrService()

	users, errResponse := mgrService.UserGetAll()

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"UserName",
		"FirstName",
		"Role",
	})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareUsersTable(users))
	table.Render()

}

func prepareUsersTable(data []types.User) [][]string {

	var rows [][]string

	for _, val := range data {

		_username := val.Username
		_firstname := val.Firstname
		_role := val.Role

		rows = append(rows, []string{_username, _firstname, _role})

	}

	return rows

}

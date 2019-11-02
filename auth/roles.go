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

package auth

const RoleAdmin = "admin"
const RoleDeveloper = "developer"
const RoleTester = "tester"
const RoleAny = ""

const RoleAccessLevelAdmin = 1
const RoleAccessLevelDeveloper = 2
const RoleAccessLevelTester = 3

func RoleAccessLevelGraph() map[string]int {

	graph := make(map[string]int)

	graph[RoleAdmin] = RoleAccessLevelAdmin
	graph[RoleDeveloper] = RoleAccessLevelDeveloper
	graph[RoleTester] = RoleAccessLevelTester

	return graph
}

func CheckRoleAccess(expectedRole string, currentRole string) bool {

	expectedAccessLevel := RoleAccessLevelGraph()[expectedRole]
	currentAccessLevel := RoleAccessLevelGraph()[currentRole]

	return currentAccessLevel <= expectedAccessLevel
}

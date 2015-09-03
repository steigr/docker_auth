/*
   Copyright 2015 Mathias Kaufmann

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package authn

import (
  "io"
  "io/ioutil"
  "os/exec"
  "bytes"
  "strings"
)

type Script struct {
	Command *string `yaml:"command,omitempty" json:"command:omitempty"`
}

type shellUsersAuth struct {
	command *string
}

func NewShellUserAuth(script *Script) *shellUsersAuth {
	return &shellUsersAuth{command: script.Command}
}

func (sua *shellUsersAuth) Authenticate(user string, password PasswordString) (bool,error) {
	cmd := exec.Command("sh","-c",*sua.command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false,nil
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false,nil
	}
	err = cmd.Start()
	io.WriteString(stdin,strings.Join([]string{user,string(password)},"\n"))
	out,err := ioutil.ReadAll(stdout)
	if err != nil {
		return false, NoMatch
	}
	if bytes.Equal(out,[]byte("OK")) {
		return false,nil
	}
	return true,nil
}

func (sua *shellUsersAuth) Stop() {
}

func (sua *shellUsersAuth) Name() string {
	return "shell"
}

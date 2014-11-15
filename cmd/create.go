/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       create.go
 *  @date       2014-11-02
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *  ------------------------------------------------------------
 */
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/web/crypto"
)

func generateSecret() string {
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_=+)")
	return crypto.RandomString(64, chars)
}

func createProject(path string) error {
	return nil
}

func Create(context *cli.Context) {
	args := context.Args()
	if len(args) != 1 {
		log.Print("Valid Project Name Missing")
	} else {
		if cwd, err := os.Getwd(); err != nil {
			panic(err)
		} else {
			path := filepath.Join(cwd, args[0])
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				panic(err)
			}
			createProject(path)
		}
	}
}

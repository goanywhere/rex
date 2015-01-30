/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2015 GoAnywhere (http://goanywhere.io).
 * ----------------------------------------------------------------------
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
 * ----------------------------------------------------------------------*/
package internal

import (
	"fmt"
	"time"
)

// Loading provides a simple dotted loading prompt while
// processing task. NOTE that task needs to manually notify
// the prompt once the task is done.
func Loading(done chan bool) {
	var (
		index      = 0
		indicators = []string{"-", "\\", "|", "/"}
	)
	go func() {
		for {
			select {
			case <-done:
				fmt.Print(" \b")
				close(done)
				return
			default:
				fmt.Printf("%s\b", indicators[index])
				if index < (len(indicators) - 1) {
					index++
				} else {
					index = 0
				}
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()
}

func Prompt(message string, values ...interface{}) {
	fmt.Printf(fmt.Sprintf("* %s", message), values...)
}

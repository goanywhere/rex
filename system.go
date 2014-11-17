/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       system.go
 *  @date       2014-10-23
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
package web

import (
	"os"
)

// Copy recursively copies files/(sub)directoires into the given path.
// *NOTE* It walks the file tree rooted at root, calling walkFn for
// each file or directory in the tree, including root, means that for
// very large directories it can be inefficient. Copy does not follow symbolic links.
func Copy(src, dest string) {

}

// Exists check if the given path exists.
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir checks if the given path is a directory.
func IsDir(path string) bool {
	src, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return src.IsDir()
}

// IsDir checks if the given path is a file.
func IsFile(path string) bool {
	src, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !src.IsDir()
}

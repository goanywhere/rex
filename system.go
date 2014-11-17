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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// AbsDir finds the absolute path for the given path.
// Supported Formats:
//	* empty path  => current working directory.
//	* '.', '..' & '~'
func AbsDir(path string) string {
	var abs string
	cwd, _ := os.Getwd()

	if path == "" || path == "." {
		abs = cwd
	} else if path == ".." {
		abs = filepath.Join(cwd, path)
	} else if strings.HasPrefix(path, "~/") {
		abs = filepath.Join(UserDir(), path[2:])
	} else if strings.HasPrefix(path, "./") {
		abs = filepath.Join(cwd, path[2:])
	} else if strings.HasPrefix(path, "../") {
		abs = filepath.Join(cwd, "..", path[2:])
	} else {
		return path
	}
	return abs
}

// Copy recursively copies files/(sub)directoires into the given path.
// *NOTE* It walks the file tree rooted at root, calling walkFn for
// each file or directory in the tree, including root, means that for
// very large directories it can be inefficient. Copy does not follow symbolic links.
func Copy(src, dest string) error {
	if !Exists(src) || strings.HasPrefix(dest, src) {
		return fmt.Errorf("Operation aborted, either the source does not exist or the destination is inside the source.")
	}

	fn := func(path string, f os.FileInfo, err error) error {
		log.Printf("[*COPY*] %s with %d bytes\n", path, f.Size())
		return nil
	}
	return filepath.Walk(src, fn)
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

// UserDir finds base path of current system user.
func UserDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

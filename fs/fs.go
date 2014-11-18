/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       fs.go
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
package fs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Abs finds the absolute path for the given path.
// Supported Formats:
//	* empty path  => current working directory.
//	* '.', '..' & '~'
// *NOTE* Abs does NOT check the existence of the path.
func Abs(path string) string {
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
// *NOTE* It uses platform's native copy commands (windows: copy, *nix: rsync).
func Copy(src, dst string) (err error) {
	var cmd *exec.Cmd
	src, dst = Abs(src), Abs(dst)
	// Determine the command we need to use.
	if runtime.GOOS == "windows" {
		// *NOTE* Not sure this will work correctly, we don't have Windows to test.
		if IsFile(src) {
			cmd = exec.Command("copy", src, dst)
		} else {
			cmd = exec.Command("xcopy", src, dst, "/S /E")
		}
	} else {
		cmd = exec.Command("rsync", "-a", src, dst)
	}

	if stdout, err := cmd.StdoutPipe(); err == nil {
		if stderr, err := cmd.StderrPipe(); err == nil {
			// Start capturing the stdout/err.
			err = cmd.Start()
			io.Copy(os.Stdout, stdout)
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(stderr)
			cmd.Wait()
			if cmd.ProcessState.String() != "exit status 0" {
				err = fmt.Errorf("\t%s\n", buffer.String())
			}
		}
	}
	return
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

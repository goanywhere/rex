/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       logger.go
 *  @date       2014-10-11
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
	"github.com/agtorre/gocolorize"
	"github.com/op/go-logging"
)

const formatter = "%{color} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"

func Logger() *logging.Logger {
	l := logging.MustGetLogger("web.go")
	logging.SetFormatter(logging.MustStringFormatter(formatter))
	return l
}

func Debug(message string) {
	colorize := gocolorize.NewColor("blue").Paint
	fmt.Println(colorize("[DEBUG] " + message))
}

func Error(message string) {
	colorize := gocolorize.NewColor("red").Paint
	fmt.Println(colorize("[ERROR] " + message))
}

func Fatal(message string) {
	colorize := gocolorize.NewColor("white:red").Paint
	fmt.Println(colorize("[FATAL] " + message))
}

func Info(message string) {
	colorize := gocolorize.NewColor("green").Paint
	fmt.Println("▶ " + colorize("[INFO] "+message))
}

func Warn(message string) {
	colorize := gocolorize.NewColor("yellow").Paint
	fmt.Println(colorize("[WARN] " + message))
}

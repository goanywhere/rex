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
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[Web.go]", 0)

func Debug(format string, values ...interface{}) {
	logger.Printf("[DEBUG] "+format, values...)
}

func Error(format string, values ...interface{}) {
	logger.Printf("[ERROR] "+format, values...)
}

func Fatal(format string, values ...interface{}) {
	logger.Fatalf("[FATAL] "+format, values...)
}

func Info(format string, values ...interface{}) {
	logger.Printf("[INFO] "+format, values...)
}

func Warn(format string, values ...interface{}) {
	logger.Printf("[WARN] "+format, values...)
}

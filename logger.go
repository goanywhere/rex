/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

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

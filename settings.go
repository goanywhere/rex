/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       settings.go
 *  @date       2014-10-16
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
	"time"

	"github.com/spf13/viper"
)

type settings struct {
	SupportedFormats []string
}

func configure(filename string) *settings {
	cwd, _ := os.Getwd()
	// --------------------
	// Application Defaults
	// --------------------
	viper.SetDefault("address", ":3000")
	viper.SetDefault("application", "webapp")
	viper.SetDefault("version", "0.0.1")
	viper.SetDefault("folder", map[string]string{
		"templates": "templates",
	})
	viper.SetDefault("XSRF", map[string]interface{}{
		"enabled": true,
	})
	// --------------------
	// User Settings
	// --------------------
	viper.AddConfigPath(cwd)      // User settings file path.
	viper.SetConfigName(filename) // Application settings file name.
	viper.ReadInConfig()

	return &settings{SupportedFormats: viper.SupportedExts}
}

func (self *settings) Get(key string) interface{} {
	return viper.Get(key)
}

func (self *settings) Set(key string, value interface{}) {
	viper.Set(key, value)
}

func (self *settings) SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func (self *settings) AllKeys() []string {
	return viper.AllKeys()
}

func (self *settings) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

func (self *settings) AutomaticEnv() {
	viper.AutomaticEnv()
}

func (self *settings) BindEnv(input ...string) (err error) {
	return viper.BindEnv(input...)
}

func (self *settings) ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

func (self *settings) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (self *settings) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (self *settings) GetInt(key string) int {
	return viper.GetInt(key)
}

func (self *settings) GetString(key string) string {
	return viper.GetString(key)
}

func (self *settings) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (self *settings) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func (self *settings) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (self *settings) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (self *settings) InConfig(key string) bool {
	return viper.InConfig(key)
}

func (self *settings) IsSet(key string) bool {
	return viper.IsSet(key)
}

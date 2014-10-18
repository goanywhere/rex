/**
 *  ------------------------------------------------------------
 *  @project	webapp
 *  @file       config.go
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
package webapp

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	SupportedFormats []string
}

func init() {
	// --------------------
	// Application Defaults
	// --------------------
	viper.SetDefault("address", ":3000")
	viper.SetDefault("application", "webapp")
	viper.SetDefault("version", "0.0.1")
}

func configure(filename string) *config {
	cwd, _ := os.Getwd()
	// --------------------
	// User Settings
	// --------------------
	viper.AddConfigPath(cwd)      // User settings file path.
	viper.SetConfigName(filename) // Application settings file name.
	viper.ReadInConfig()
	return &config{SupportedFormats: viper.SupportedExts}
}

func (config *config) Get(key string) interface{} {
	return viper.Get(key)
}

func (config *config) Set(key string, value interface{}) {
	viper.Set(key, value)
}

func (config *config) SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func (config *config) AllKeys() []string {
	return viper.AllKeys()
}

func (config *config) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

func (config *config) AutomaticEnv() {
	viper.AutomaticEnv()
}

func (config *config) BindEnv(input ...string) (err error) {
	return viper.BindEnv(input...)
}

func (config *config) ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

func (config *config) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (config *config) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (config *config) GetInt(key string) int {
	return viper.GetInt(key)
}

func (config *config) GetString(key string) string {
	return viper.GetString(key)
}

func (config *config) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (config *config) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func (config *config) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (config *config) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (config *config) InConfig(key string) bool {
	return viper.InConfig(key)
}

func (config *config) IsSet(key string) bool {
	return viper.IsSet(key)
}

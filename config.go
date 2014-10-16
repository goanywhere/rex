/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       viper.go
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

type Config struct {
	SupportedFormats []string
}

func init() {
	// --------------------
	// Application Defaults
	// --------------------
	viper.SetDefault("address", ":9394")
	viper.SetDefault("application", "webapp")
	viper.SetDefault("version", "0.0.1")
}

func Configure(filename string) *Config {
	cwd, _ := os.Getwd()
	// --------------------
	// User Settings
	// --------------------
	viper.AddConfigPath(cwd)      // User settings file path.
	viper.SetConfigName(filename) // Application settings file name.
	viper.ReadInConfig()
	return &Config{SupportedFormats: viper.SupportedExts}
}

func (config *Config) Get(key string) interface{} {
	return viper.Get(key)
}

func (config *Config) Set(key string, value interface{}) {
	viper.Set(key, value)
}

func (config *Config) SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func (config *Config) AllKeys() []string {
	return viper.AllKeys()
}

func (config *Config) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

func (config *Config) AutomaticEnv() {
	viper.AutomaticEnv()
}

func (config *Config) BindEnv(input ...string) (err error) {
	return viper.BindEnv(input...)
}

func (config *Config) ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

func (config *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (config *Config) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (config *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

func (config *Config) GetString(key string) string {
	return viper.GetString(key)
}

func (config *Config) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (config *Config) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func (config *Config) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (config *Config) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (config *Config) InConfig(key string) bool {
	return viper.InConfig(key)
}

func (config *Config) IsSet(key string) bool {
	return viper.IsSet(key)
}

// Copyright 2018. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"strings"

	"github.com/go-ini/ini"
)

type Config struct {
	Path  string
	Ini   *ini.File
	dirty bool
}

func NewConfig(path string) (config *Config, err error) {
	config = &Config{Path: path}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		config.Ini = ini.Empty()
		return
	}

	config.Ini, err = ini.Load(path)
	return
}

func (c *Config) Save() error {
	if c.dirty {
		err := c.Ini.SaveTo(c.Path)
		c.dirty = false
		return err
	}
	return nil
}

func (c *Config) Get(sectionName string, keyName string) string {
	section := c.Ini.Section(sectionName)
	key := section.Key(keyName)
	if key != nil {
		return key.String()
	}

	return ""
}

func (c *Config) Set(sectionName string, key string, value string) {
	section := c.Ini.Section(sectionName)
	section.Key(key).SetValue(value)
	c.dirty = true
}

func (c *Config) Unset(sectionName string, key string) {
	section := c.Ini.Section(sectionName)
	section.DeleteKey(key)
	c.dirty = true
}

func (c *Config) GetIni() *ini.File {
	return c.Ini
}

func (c *Config) ExportEnv() {
	for _, section := range c.Ini.Sections() {
		for _, key := range section.Keys() {
			envVar := "AKAMAI_" + strings.ToUpper(section.Name()) + "_"
			envVar += strings.ToUpper(strings.Replace(key.Name(), "-", "_", -1))
			os.Setenv(envVar, key.String())
		}
	}
}

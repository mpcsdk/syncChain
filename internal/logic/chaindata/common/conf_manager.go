// Copyright 2020 The RangersProtocol Authors
// This file is part of the RocketProtocol library.
//
// The RangersProtocol library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The RangersProtocol library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the RangersProtocol library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"os"
	"strings"
	"sync"

	ini "github.com/glacjay/goini"
)

type ConfManager interface {
	GetStrings(section string) map[string]string

	// GetString read basic conf from tas.conf file
	GetString(section string, key string, defaultValue string) string
	GetBool(section string, key string, defaultValue bool) bool
	GetDouble(section string, key string, defaultValue float64) float64
	GetInt(section string, key string, defaultValue int) int

	//set basic conf to tas.conf file
	SetString(section string, key string, value string)
	SetBool(section string, key string, value bool)
	SetDouble(section string, key string, value float64)
	SetInt(section string, key string, value int)

	//delete basic conf
	Del(section string, key string)

	GetSectionManager(section string) SectionConfManager
}

type SectionConfManager interface {
	//read basic conf from tas.conf file
	GetString(key string, defaultValue string) string
	GetBool(key string, defaultValue bool) bool
	GetDouble(key string, defaultValue float64) float64
	GetInt(key string, defaultValue int) int

	//set basic conf to tas.conf file
	SetString(key string, value string)
	SetBool(key string, value bool)
	SetDouble(key string, value float64)
	SetInt(key string, value int)

	//delete basic conf
	Del(key string)
}

type ConfFileManager struct {
	path string
	dict ini.Dict
	lock sync.RWMutex
}

type SectionConfFileManager struct {
	section string
	cfm     ConfManager
}

var GlobalConf ConfManager

func InitConf(path string) {
	if GlobalConf == nil {
		GlobalConf = newConfINIManager(path)
	}
}

func newConfINIManager(path string) ConfManager {
	cs := &ConfFileManager{
		path: path,
	}

	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	cs.dict = ini.MustLoad(path)

	return cs
}

func (cs *ConfFileManager) GetSectionManager(section string) SectionConfManager {
	return &SectionConfFileManager{
		section: section,
		cfm:     cs,
	}
}

func (sfm *SectionConfFileManager) GetString(key string, defaultValue string) string {
	return sfm.cfm.GetString(sfm.section, key, defaultValue)
}

func (sfm *SectionConfFileManager) GetBool(key string, defaultValue bool) bool {
	return sfm.cfm.GetBool(sfm.section, key, defaultValue)
}

func (sfm *SectionConfFileManager) GetDouble(key string, defaultValue float64) float64 {
	return sfm.cfm.GetDouble(sfm.section, key, defaultValue)
}

func (sfm *SectionConfFileManager) GetInt(key string, defaultValue int) int {
	return sfm.cfm.GetInt(sfm.section, key, defaultValue)
}

func (sfm *SectionConfFileManager) SetString(key string, value string) {
	sfm.cfm.SetString(sfm.section, key, value)
}

func (sfm *SectionConfFileManager) SetBool(key string, value bool) {
	sfm.cfm.SetBool(sfm.section, key, value)
}

func (sfm *SectionConfFileManager) SetDouble(key string, value float64) {
	sfm.cfm.SetDouble(sfm.section, key, value)
}

func (sfm *SectionConfFileManager) SetInt(key string, value int) {
	sfm.cfm.SetInt(sfm.section, key, value)
}

func (sfm *SectionConfFileManager) Del(key string) {
	sfm.cfm.Del(sfm.section, key)
}

func (cs *ConfFileManager) GetStrings(section string) map[string]string {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.dict[section]
}

func (cs *ConfFileManager) GetString(section string, key string, defaultValue string) string {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	if v, ok := cs.dict.GetString(strings.ToLower(section), strings.ToLower(key)); ok {
		return v
	}
	return defaultValue
}

func (cs *ConfFileManager) GetBool(section string, key string, defaultValue bool) bool {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	if v, ok := cs.dict.GetBool(strings.ToLower(section), strings.ToLower(key)); ok {
		return v
	}
	return defaultValue
}

func (cs *ConfFileManager) GetDouble(section string, key string, defaultValue float64) float64 {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	if v, ok := cs.dict.GetDouble(strings.ToLower(section), strings.ToLower(key)); ok {
		return v
	}
	return defaultValue
}

func (cs *ConfFileManager) GetInt(section string, key string, defaultValue int) int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	if v, ok := cs.dict.GetInt(strings.ToLower(section), strings.ToLower(key)); ok {
		return v
	}
	return defaultValue
}

func (cs *ConfFileManager) SetString(section string, key string, value string) {
	cs.update(func() {
		cs.dict.SetString(strings.ToLower(section), strings.ToLower(key), value)
	})
}

func (cs *ConfFileManager) SetBool(section string, key string, value bool) {
	cs.update(func() {
		cs.dict.SetBool(strings.ToLower(section), strings.ToLower(key), value)
	})
}

func (cs *ConfFileManager) SetDouble(section string, key string, value float64) {
	cs.update(func() {
		cs.dict.SetDouble(strings.ToLower(section), strings.ToLower(key), value)
	})
}

func (cs *ConfFileManager) SetInt(section string, key string, value int) {
	cs.update(func() {
		cs.dict.SetInt(strings.ToLower(section), strings.ToLower(key), value)
	})
}

func (cs *ConfFileManager) Del(section string, key string) {
	cs.update(func() {
		cs.dict.Delete(strings.ToLower(section), strings.ToLower(key))
	})
}

func (cs *ConfFileManager) update(updator func()) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	updator()
	cs.store()
}

func (cs *ConfFileManager) store() {
	err := ini.Write(cs.path, &cs.dict)
	if err != nil {

	}
}

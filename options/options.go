/*
 * Copyright (c) 2020 Juan Medina.
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

// Package options contains our game Options
package options

import (
	"encoding/json"
	"github.com/juan-medina/gosge/components/color"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

// Settings allow to get and set in game settings
type Settings interface {
	// Get an in game setting with a default value
	Get(setting string, _default interface{}) interface{}
	// Set an in game setting
	Set(setting string, value interface{})

	// GetFloat32 get an in game float32 setting with a default value
	GetFloat32(setting string, _default float32) float32
	// SetFloat32 set an in game float32 setting
	SetFloat32(setting string, value float32)

	// GetString get an in game string setting with a default value
	GetString(setting string, _default string) string
	// SetString set an in game string setting
	SetString(setting string, value string)

	// GetIn32 get an in game int32 setting with a default value
	GetIn32(setting string, _default int32) int32
	// SetInt32 set an in game int32 setting
	SetInt32(setting string, value int32)
}

// Options are our game options
type Options struct {
	Title      string                 // Title is the game title
	BackGround color.Solid            // BackGround is the background color.Color
	Monitor    int                    // Monitor is the monitor that we will use for full screen
	Icon       string                 // Icon is a path for a PNG containing the application icon
	Windowed   bool                   // Windowed will indicate if we want the game on a window
	Width      int                    // Width is the desired width
	Height     int                    // Height is the desired height
	Settings   map[string]interface{} // Settings store our game settings
}

// Get an in game setting with a default value
func (o *Options) Get(setting string, _default interface{}) interface{} {
	var v interface{}
	var ok bool

	if v, ok = o.Settings[setting]; !ok {
		o.Settings[setting] = _default
		return _default
	}

	return v
}

// Set an in game setting
func (o *Options) Set(setting string, value interface{}) {
	o.Settings[setting] = value
}

// GetFloat32 get an in game float32 setting with a default value
func (o *Options) GetFloat32(setting string, _default float32) float32 {
	result := _default
	value := o.Get(setting, _default)
	switch v := value.(type) {
	case float32:
		result = v
	case float64:
		result = float32(v)
	case int32:
		result = float32(v)
	case int:
		result = float32(v)
	case int64:
		result = float32(v)
	case string:
		var err error
		var f64 float64
		if f64, err = strconv.ParseFloat(v, 32); err == nil {
			result = float32(f64)
		}
	}
	return result
}

// SetFloat32 set an in game float32 setting
func (o *Options) SetFloat32(setting string, value float32) {
	o.Set(setting, value)
}

// GetString get an in game string setting with a default value
func (o *Options) GetString(setting string, _default string) string {
	result := _default
	value := o.Get(setting, _default)
	switch v := value.(type) {
	case float32:
		result = strconv.FormatFloat(float64(v), 'E', -1, 32)
	case float64:
		result = strconv.FormatFloat(v, 'E', -1, 32)
	case int32:
		result = strconv.Itoa(int(v))
	case int:
		result = strconv.Itoa(v)
	case int64:
		result = strconv.FormatInt(v, 10)
	case string:
		result = v
	}
	return result
}

// SetString set an in game string setting
func (o *Options) SetString(setting string, value string) {
	o.Set(setting, value)
}

// GetIn32 get an in game int32 setting with a default value
func (o *Options) GetIn32(setting string, _default int32) int32 {
	result := _default
	value := o.Get(setting, _default)
	switch v := value.(type) {
	case float32:
		result = int32(v)
	case float64:
		result = int32(v)
	case int32:
		result = v
	case int:
		result = int32(v)
	case int64:
		result = int32(v)
	case string:
		var err error
		var i int
		if i, err = strconv.Atoi(v); err == nil {
			result = int32(i)
		}
	}
	return result
}

// SetInt32 set an in game int32 setting
func (o *Options) SetInt32(setting string, value int32) {
	o.Set(setting, value)
}

// Save current options
func (o Options) Save() error {
	var err error
	var fileName string

	if fileName, err = o.getFilePath(); err != nil {
		return err
	}

	if o.Settings == nil {
		o.Settings = make(map[string]interface{}, 0)
	}

	var content []byte
	if content, err = json.MarshalIndent(o, "", "    "); err != nil {
		return err
	}
	if err = ioutil.WriteFile(fileName, content, 0644); err != nil {
		return err
	}
	return err
}

func (o Options) getFilePath() (string, error) {
	var err error
	var home string
	if home, err = os.UserHomeDir(); err != nil {
		return "", err
	}

	name := strings.ToLower(o.Title)
	name = "." + strings.Replace(name, " ", "_", -1)
	gameDir := path.Join(home, name)

	if _, err = os.Stat(gameDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(gameDir, 0755); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return path.Join(gameDir, "options.json"), nil
}

func (o Options) doConfigExist() (bool, error) {
	var err error
	var filePath string

	if filePath, err = o.getFilePath(); err != nil {
		return false, err
	}

	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Load options
func (o *Options) Load() error {
	var err error
	var exists bool

	if exists, err = o.doConfigExist(); err != nil {
		return err
	}

	if !exists {
		if o.Settings == nil {
			o.Settings = make(map[string]interface{}, 0)
		}
		if err = o.Save(); err != nil {
			return err
		}
	}

	var data []byte
	var filePath string

	if filePath, err = o.getFilePath(); err != nil {
		return err
	}

	if data, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}

	if err = json.Unmarshal(data, o); err != nil {
		return err
	}

	return nil
}

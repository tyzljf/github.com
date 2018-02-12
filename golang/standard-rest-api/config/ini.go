package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	defaultSection = "General"
	bNumComment    = []byte{'#'}
	bSemComment    = []byte{';'}
	bEmpty         = []byte{}
	bEqual         = []byte{'='}
	sectionStart   = []byte{'['}
	sectionEnd     = []byte{']'}
	lineBreak      = "\n"
)

// IniConfig implements config to parse ini file
type IniConfig struct {
}

func (ini *IniConfig) Parse(filename string) (Configer, error) {
	return ini.parseFile(filename)
}

func (ini *IniConfig) parseFile(filename string) (*IniConfigerContainer, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ini.parseData(filepath.Dir(filename), data)
}

func (ini *IniConfig) parseData(dir string, data []byte) (*IniConfigerContainer, error) {
	cfg := &IniConfigerContainer{
		data:           make(map[string]map[string]string),
		sectionComment: make(map[string]string),
		keyComment:     make(map[string]string),
		RWMutex:        sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()

	var comment bytes.Buffer
	buf := bufio.NewReader(bytes.NewBuffer(data))
	section := defaultSection

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if _, ok := err.(*os.PathError); ok {
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, bEmpty) {
			continue
		}
		var bComment []byte
		switch {
		case bytes.HasPrefix(line, bNumComment):
			bComment = bNumComment
		case bytes.HasPrefix(line, bSemComment):
			bComment = bSemComment
		}
		if bComment != nil {
			line = bytes.TrimLeft(line, string(bComment))
			if comment.Len() > 0 {
				comment.WriteByte('\n')
			}
			comment.Write(line)
			continue
		}

		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			section = strings.ToLower(string(line[1 : len(line)-1]))
			if comment.Len() > 0 {
				cfg.sectionComment[section] = comment.String()
				comment.Reset()
			}
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
			continue
		}

		if _, ok := cfg.data[section]; !ok {
			cfg.data[section] = make(map[string]string)
		}
		keyValue := bytes.SplitN(line, bEqual, 2)

		key := string(bytes.TrimSpace(keyValue[0]))
		key = strings.ToLower(key)

		// handle include "other.conf"
		if len(keyValue) == 1 && strings.HasPrefix(key, "include") {
			includefiles := strings.Fields(key)

			if includefiles[0] == "include" && len(includefiles) == 2 {
				otherfile := strings.Trim(includefiles[1], "\"")
				if filepath.IsAbs(otherfile) {
					otherfile = filepath.Join(dir, otherfile)
				}

				i, err := ini.parseFile(otherfile)
				if err != nil {
					return nil, err
				}

				for sec, dt := range i.data {
					if _, ok := cfg.data[sec]; !ok {
						cfg.data[sec] = make(map[string]string)
					}
					for k, v := range dt {
						cfg.data[sec][k] = v
					}
				}

				for sec, comm := range i.sectionComment {
					cfg.sectionComment[sec] = comm
				}

				for k, comm := range i.keyComment {
					cfg.keyComment[k] = comm
				}

				continue
			}
		}

		if len(keyValue) != 2 {
			return nil, errors.New("read the content error: \"" + string(line) + "\", should key = value")
		}

		val := string(bytes.TrimSpace(keyValue[1]))
		cfg.data[section][key] = val
		if comment.Len() > 0 {
			cfg.keyComment[section+"."+key] = comment.String()
			comment.Reset()
		}
	}

	return cfg, nil
}

func (ini *IniConfigerContainer) Bool(key string) (bool, error) {
	return ParseBool(ini.getData(key))
}

func (ini *IniConfigerContainer) DefaultBool(key string, defaultval bool) bool {
	v, err := ini.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ini *IniConfigerContainer) Int(key string) (int, error) {
	return strconv.Atoi(ini.getData(key))
}

func (ini *IniConfigerContainer) DefaultInt(key string, defaultval int) int {
	v, err := ini.Int(ini.getData(key))
	if err != nil {
		return defaultval
	}
	return v
}

func (ini *IniConfigerContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(ini.getData(key), 10, 64)
}

func (ini *IniConfigerContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := ini.Int64(ini.getData(key))
	if err != nil {
		return defaultval
	}
	return v
}

func (ini *IniConfigerContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(ini.getData(key), 64)
}

func (ini *IniConfigerContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := ini.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ini *IniConfigerContainer) String(key string) string {
	return ini.getData(key)
}

func (ini *IniConfigerContainer) DefaultString(key string, defaultval string) string {
	v := ini.String(key)
	if v == "" {
		return defaultval
	}
	return v
}

// section.key or key
func (ini *IniConfigerContainer) getData(key string) string {
	if len(key) == 0 {
		return ""
	}
	ini.RLock()
	defer ini.RUnlock()

	var (
		section, k string
		sectionKey = strings.Split(strings.ToLower(key), "::")
	)

	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = defaultSection
		k = sectionKey[0]
	}

	if v, ok := ini.data[section]; ok {
		if vv, ok := v[k]; ok {
			return vv
		}
	}

	return ""
}

func (ini *IniConfigerContainer) Set(key, value string) error {
	ini.Lock()
	defer ini.Unlock()
	if len(key) == 0 {
		return errors.New("Key is empty")
	}

	var (
		section, k string
		sectionKey = strings.Split(strings.ToLower(key), "::")
	)

	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = defaultSection
		k = sectionKey[0]
	}

	if _, ok := ini.data[section]; !ok {
		ini.data[section] = make(map[string]string)
	}

	ini.data[section][k] = value
	return nil
}

func (ini *IniConfigerContainer) saveConfigFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get section or key comments
	getCommentStr := func(section, key string) string {
		var (
			comment string
			ok      bool
		)
		if len(key) == 0 {
			comment, ok = ini.sectionComment[section]
		} else {
			comment, ok = ini.keyComment[section+"."+key]
		}

		if ok {
			// Empty comment
			if len(comment) == 0 || len(strings.TrimSpace(comment)) == 0 {
				return string(bNumComment)
			}

			prefix := string(bNumComment)
			return prefix + strings.Replace(comment, lineBreak, lineBreak+prefix, -1)
		}

		return ""
	}

	buf := bytes.NewBuffer(nil)
	// Save default section at first place
	if dt, ok := ini.data[defaultSection]; ok {
		for key, val := range dt {
			if key != " " {
				// Write key comments
				if v := getCommentStr(defaultSection, key); len(v) > 0 {
					if _, err := buf.WriteString(v + lineBreak); err != nil {
						return err
					}
				}

				// Write key and value
				if _, err := buf.WriteString(key + string(bEqual) + val + lineBreak); err != nil {
					return err
				}
			}
		}

		// Put a line between sections
		if _, err := buf.WriteString(lineBreak); err != nil {
			return err
		}
	}
	// Saved named section
	for section, dt := range ini.data {
		if section != defaultSection {
			// Write section comments
			if v := getCommentStr(section, ""); len(v) > 0 {
				if _, err := buf.WriteString(v + lineBreak); err != nil {
					return err
				}
			}

			// Write section name
			if _, err := buf.WriteString(string(sectionStart) + section + string(sectionEnd) + lineBreak); err != nil {
				return err
			}

			for key, val := range dt {
				if key != " " {
					// Write key comments
					if v := getCommentStr(section, key); len(v) > 0 {
						if _, err := buf.WriteString(v + lineBreak); err != nil {
							return err
						}
					}
					// Write key and value
					if _, err := buf.WriteString(key + string(bEqual) + val + lineBreak); err != nil {
						return err
					}
				}
			}
		}

		// Put a line between sections
		if _, err := buf.WriteString(lineBreak); err != nil {
			return err
		}
	}
	_, err = buf.WriteTo(f)
	return err
}

type IniConfigerContainer struct {
	data           map[string]map[string]string
	sectionComment map[string]string
	keyComment     map[string]string
	sync.RWMutex
}

func init() {
	Register("ini", &IniConfig{})
}

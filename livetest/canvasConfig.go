package main

// define operations on canvasConfig configuration settings object.

import "strings"

// Each string in the array is a single setting, made up of a setting path and value(s).
// Settings are added to the string with format "<keypath>=<value>;" or "<keypath>=<value0>,<value1>..;"
// literal semicolons and equals signs in a key value will be escaped with "\". TODO:implement escaping.
type canvasConfig []string

// get all properties for a given key sub-path. ex. Props("people.max") = "name=max;height=6ft3"
func (cc *canvasConfig) Props(path string) canvasConfig {
	var props canvasConfig
	for _, s := range *cc {
		pair := strings.Split(s, "=")
		key := pair[0]
		val := pair[1]
		if key == path {
			panic("ERROR: Cannot get Props of a fullly qualified configuration setting path.")
		}
		if !strings.HasSuffix(val, ";") {
			panic("illegal value format in property")
		}
		if key, ok := strings.CutPrefix(key, path); ok {
			subp := key + "=" + val
			props = append(props, subp)
		}
	}
	return props
}

// read the value(s) for a given key path. eg Get("people.max.height") = {"6ft3"}
func (cc *canvasConfig) Get(key string) []string {
	for _, s := range *cc {
		pair := strings.Split(s, "=")
		k := pair[0]
		v := pair[1]
		var ok bool
		if v, ok = strings.CutSuffix(v, ";"); !ok {
			panic("illegal value format in property")
		}
		if k == key {
			vals := strings.Split(v, ",")
			return vals
		}
	}
	return nil //no key found
}

// set configuration key to value. if key already exists, replace it. If a conflicting key exists, replace it.
// conflicts occur when a key path contains another key path as a sub-path.
// ex. people.max="max gara" and people.max.height="6ft3".
// the people.max "file" is defined explicitly to contain one element (the string "max gara") so it can't also contain other files.
// if this were allowed, then the properties list for a sub-path could contain properties without names.
func (cc *canvasConfig) Set(key string, val string) {
	var found bool
	prop := key + "=" + val + ";"
	for idx, s := range *cc {
		pair := strings.Split(s, "=")
		k := pair[0]
		if k == key {
			(*cc)[idx] = prop //replace previous setting
			found = true
			continue
		}
		if strings.HasPrefix(key, k) && key[len(k)] == '.' {
			//remove conflicting sub-path property. path diff must be a whole word (period at end)
			cc.Remove(k)
			continue
		}
		if strings.HasPrefix(k, key) && k[len(key)] == '.' {
			//remove conflicting extra-path property. path diff must be a whole word (period at end)
			cc.Remove(k)
			continue
		}
	}
	if !found {
		*cc = append(*cc, prop) //create new setting
	}
}

// remove a specific setting without affecting any others. (changes slice)
func (cc *canvasConfig) Remove(key string) {
	l := len(*cc)
	for idx, s := range *cc {
		pair := strings.Split(s, "=")
		k := pair[0]
		if k == key {
			copy((*cc)[idx:], (*cc)[idx+1:])
			*cc = (*cc)[:l-1] //drop slice end
			break
		}
	}
}

// add a new value for a key without removing existing values. Note: will still remove other conflicting values.
func (cc *canvasConfig) Append(key, val string) {
	for idx, s := range *cc {
		pair := strings.Split(s, "=")
		k := pair[0]
		if k == key {
			(*cc)[idx] = (*cc)[idx] + "," + val
			return
		}
	}
}

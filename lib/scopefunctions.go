/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"log"
	"strconv"
	"strings"
)

func getScopeFunction(scopeName string) (result Scope) {
	switch {
	case scopeName == "any":
		result = ScopeAny
	case scopeName == "all":
		result = ScopeAll
	case scopeName == "none":
		result = ScopeNone
	case strings.HasPrefix(scopeName, "max"):
		number, err := strconv.Atoi(strings.TrimSpace(scopeName[3:]))
		if err != nil {
			log.Fatal("cant parse scopename ", scopeName, err)
		}
		result = func(set []bool) bool {
			return ScopeMax(set, number)
		}
	case strings.HasPrefix(scopeName, "min"):
		number, err := strconv.Atoi(strings.TrimSpace(scopeName[3:]))
		if err != nil {
			log.Fatal("cant parse scopename ", scopeName, err)
		}
		result = func(set []bool) bool {
			return ScopeMin(set, number)
		}
	case isNumber(scopeName):
		number, err := strconv.Atoi(scopeName)
		if err != nil {
			log.Fatal("cant parse scopename ", scopeName, err)
		}
		result = func(set []bool) bool {
			return ScopeNumber(set, number)
		}
	default:
		log.Fatal("unknown scope name:", scopeName)
	}
	return
}

func isNumber(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

func ScopeAny(set []bool) bool {
	for _, element := range set {
		if element {
			return true
		}
	}
	return false
}

func ScopeAll(set []bool) bool {
	for _, element := range set {
		if !element {
			return false
		}
	}
	return true
}

func ScopeNone(set []bool) bool {
	for _, element := range set {
		if element {
			return false
		}
	}
	return true
}

func ScopeMax(set []bool, max int) bool {
	count := 0
	for _, element := range set {
		if element {
			count++
			if count > max {
				return false
			}
		}
	}
	return true
}

func ScopeMin(set []bool, min int) bool {
	count := 0
	for _, element := range set {
		if element {
			count++
			if count >= min {
				return true
			}
		}
	}
	return false
}

func ScopeNumber(set []bool, number int) bool {
	count := 0
	for _, element := range set {
		if element {
			count++
		}
	}
	return count == number
}

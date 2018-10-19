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
	"regexp"
	"strconv"
)

func getOperatorFunction(operatorName string) (result Operator) {
	switch operatorName {
	case "==":
		result = OperatorEq
	case "!=":
		result = OperatorUneq
	case ">=":
		result = OperatorBE
	case "<=":
		result = OperatorSE
	case ">":
		result = OperatorB
	case "<":
		result = OperatorS
	case "regex":
		result = OperatorRegex
	default:
		log.Fatal("unknown operator name:", operatorName)
	}
	return
}

func OperatorEq(ruleValue string, value string) bool {
	return ruleValue == value
}

func OperatorUneq(ruleValue string, value string) bool {
	return ruleValue != value
}

func OperatorBE(ruleValue string, value string) bool {
	ruleFloat, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	return valueFloat >= ruleFloat
}

func OperatorB(ruleValue string, value string) bool {
	ruleFloat, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	return valueFloat > ruleFloat
}

func OperatorSE(ruleValue string, value string) bool {
	ruleFloat, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	return valueFloat <= ruleFloat
}

func OperatorS(ruleValue string, value string) bool {
	ruleFloat, err := strconv.ParseFloat(ruleValue, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("error on string to float parsing in operator: ", err, ruleValue, value)
		return false
	}
	return valueFloat < ruleFloat
}

func OperatorRegex(ruleValue string, value string) bool {
	result, err := regexp.MatchString(ruleValue, value)
	if err != nil {
		log.Println("regex error: ", err)
		return false
	}
	return result
}

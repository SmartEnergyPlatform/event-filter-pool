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
	"github.com/SmartEnergyPlatform/event-filter-pool/util"
	"log"

	"strings"

	"github.com/JumboInteractiveLimited/jsonpath"
)

type CompiledRule struct {
	Paths    []*jsonpath.Path
	Scope    Scope
	Operator Operator
	Value    string
}

type Filter struct {
	FilterDesc
	Id            string
	compiledRules []CompiledRule
	scopeFn       Scope
}

type Scope func([]bool) bool
type Operator func(string, string) bool

func NewFilter(assignment FilterDeployment) (result Filter, err error) {
	result = Filter{FilterDesc: assignment.Filter, Id: assignment.FilterId, scopeFn: getScopeFunction(assignment.Filter.Scope)}
	result.compiledRules, err = compileRules(assignment.Filter.Rules)
	if result.Topic == "" {
		result.Topic = util.Config.FilterTopic
	}
	return
}

func compileRules(rules []Rule) (result []CompiledRule, err error) {
	for _, rule := range rules {
		compiled, err := compileRule(rule)
		if err != nil {
			return result, err
		}
		result = append(result, compiled)
	}
	return
}

func compileRule(rule Rule) (result CompiledRule, err error) {
	paths, err := jsonpath.ParsePaths(rule.Path + "+")
	if err != nil {
		log.Println("ERROR in jsonpath: ", err)
		return result, err
	}
	result = CompiledRule{
		Paths:    paths,
		Value:    rule.Value,
		Scope:    getScopeFunction(rule.Scope),
		Operator: getOperatorFunction(rule.Operator),
	}
	return
}

func (filter *Filter) Check(value string) (ok bool, jsonValue string) {
	jsonValue, err := filter.getJson(value)
	if err != nil {
		log.Println("ERROR: unable to get json format of event ", err, *filter)
		return false, jsonValue
	}
	ruleResults := []bool{}
	for _, rule := range filter.compiledRules {
		eval, err := jsonpath.EvalPathsInBytes([]byte(jsonValue), rule.Paths)
		if err != nil {
			log.Println("filter.check() error: ", err)
			return false, jsonValue
		}

		pathResults := []bool{}
		for pathResult, ok := eval.Next(); ok; pathResult, ok = eval.Next() {
			value := string(pathResult.Value)
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = value[1 : len(value)-1]
			}
			operatorResult := rule.Operator(rule.Value, value)
			log.Println(rule.Value, " | ", value, " --> ", operatorResult)
			pathResults = append(pathResults, operatorResult)
		}
		log.Println("local scope: ", pathResults)
		ruleResults = append(ruleResults, rule.Scope(pathResults))
	}
	log.Println("global scope: ", ruleResults)
	return filter.scopeFn(ruleResults), jsonValue
}

func (filter *Filter) getJson(str string) (string, error) {
	return str, nil
}

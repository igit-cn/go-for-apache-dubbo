/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"github.com/creasty/defaults"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
)

// MethodConfig defines method config
type MethodConfig struct {
	InterfaceId                 string
	InterfaceName               string
	Name                        string `yaml:"name"  json:"name,omitempty" property:"name"`
	Retries                     string `yaml:"retries"  json:"retries,omitempty" property:"retries"`
	LoadBalance                 string `yaml:"loadbalance"  json:"loadbalance,omitempty" property:"loadbalance"`
	Weight                      int64  `yaml:"weight"  json:"weight,omitempty" property:"weight"`
	TpsLimitInterval            string `yaml:"tps.limit.interval" json:"tps.limit.interval,omitempty" property:"tps.limit.interval"`
	TpsLimitRate                string `yaml:"tps.limit.rate" json:"tps.limit.rate,omitempty" property:"tps.limit.rate"`
	TpsLimitStrategy            string `yaml:"tps.limit.strategy" json:"tps.limit.strategy,omitempty" property:"tps.limit.strategy"`
	ExecuteLimit                string `yaml:"execute.limit" json:"execute.limit,omitempty" property:"execute.limit"`
	ExecuteLimitRejectedHandler string `yaml:"execute.limit.rejected.handler" json:"execute.limit.rejected.handler,omitempty" property:"execute.limit.rejected.handler"`
	Sticky                      bool   `yaml:"sticky"   json:"sticky,omitempty" property:"sticky"`
	RequestTimeout              string `yaml:"timeout"  json:"timeout,omitempty" property:"timeout"`
}

// nolint
func (mc *MethodConfig) Prefix() string {
	if len(mc.InterfaceId) != 0 {
		return constant.DUBBO + "." + mc.InterfaceName + "." + mc.InterfaceId + "." + mc.Name + "."
	}

	return constant.DUBBO + "." + mc.InterfaceName + "." + mc.Name + "."
}

func (mc *MethodConfig) Init() error {
	return mc.check()
}

func initProviderMethodConfig(sc *ServiceConfig) error {
	methods := sc.Methods
	if methods == nil {
		return nil
	}
	for _, method := range methods {
		if err := method.check(); err != nil {
			return err
		}
	}
	sc.Methods = methods
	return nil
}

// check set default value and verify
func (mc *MethodConfig) check() error {
	if err := defaults.Set(mc); err != nil {
		return err
	}
	return verify(mc)
}

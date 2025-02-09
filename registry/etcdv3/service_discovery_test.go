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

package etcdv3

//
//var testName = "test"
//
//func setUp() {
//	config.GetRootConfig().ServiceDiscoveries[testName] = &config.ServiceDiscoveryConfig{
//		Protocol:  "etcdv3",
//		RemoteRef: testName,
//	}
//
//	//config.GetRootConfig().Remotes[testName] = &config.RemoteConfig{
//	//	Address:    "localhost:2379",
//	//	TimeoutStr: "10s",
//	//}
//}
//
//func Test_newEtcdV3ServiceDiscovery(t *testing.T) {
//	name := constant.ETCDV3_KEY
//	_, err := newEtcdV3ServiceDiscovery()
//
//	// warn: log configure file name is nil
//	assert.NotNil(t, err)
//
//	sdc := &config.ServiceDiscoveryConfig{
//		Protocol:  "etcdv3",
//		RemoteRef: "mock",
//	}
//	config.GetRootConfig().ServiceDiscoveries[name] = sdc
//
//	_, err = newEtcdV3ServiceDiscovery()
//
//	// RemoteConfig not found
//	assert.NotNil(t, err)
//
//	//config.GetRootConfig().Remotes["mock"] = &config.RemoteConfig{
//	//	Address:    "localhost:2379",
//	//	TimeoutStr: "10s",
//	//}
//
//	res, err := newEtcdV3ServiceDiscovery()
//	assert.Nil(t, err)
//	assert.NotNil(t, res)
//}
//
//func TestEtcdV3ServiceDiscovery_GetDefaultPageSize(t *testing.T) {
//	setUp()
//	serviceDiscovery := &etcdV3ServiceDiscovery{}
//	assert.Equal(t, registry.DefaultPageSize, serviceDiscovery.GetDefaultPageSize())
//}

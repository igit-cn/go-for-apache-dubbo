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

package etcd

import (
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"

	"go.etcd.io/etcd/server/v3/embed"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/metadata/identifier"
)

const defaultEtcdV3WorkDir = "/tmp/default-dubbo-go-registry.etcd"

func initEtcd(t *testing.T) *embed.Etcd {
	DefaultListenPeerURLs := "http://localhost:2380"
	DefaultListenClientURLs := "http://localhost:2379"
	lpurl, _ := url.Parse(DefaultListenPeerURLs)
	lcurl, _ := url.Parse(DefaultListenClientURLs)
	cfg := embed.NewConfig()
	cfg.LPUrls = []url.URL{*lpurl}
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.Dir = defaultEtcdV3WorkDir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestEtcdMetadataReportFactory_CreateMetadataReport(t *testing.T) {
	e := initEtcd(t)
	url, err := common.NewURL("registry://127.0.0.1:2379", common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)))
	if err != nil {
		t.Fatal(err)
	}
	metadataReportFactory := &etcdMetadataReportFactory{}
	metadataReport := metadataReportFactory.CreateMetadataReport(url)
	assert.NotNil(t, metadataReport)
	e.Close()
}

func TestEtcdMetadataReport_CRUD(t *testing.T) {
	e := initEtcd(t)
	url, err := common.NewURL("registry://127.0.0.1:2379", common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)))
	if err != nil {
		t.Fatal(err)
	}
	metadataReportFactory := &etcdMetadataReportFactory{}
	metadataReport := metadataReportFactory.CreateMetadataReport(url)
	assert.NotNil(t, metadataReport)

	err = metadataReport.StoreConsumerMetadata(newMetadataIdentifier("consumer"), "consumer metadata")
	assert.Nil(t, err)

	err = metadataReport.StoreProviderMetadata(newMetadataIdentifier("provider"), "provider metadata")
	assert.Nil(t, err)

	serviceMi := newServiceMetadataIdentifier()
	serviceUrl, err := common.NewURL("registry://localhost:8848", common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)))
	assert.Nil(t, err)
	err = metadataReport.SaveServiceMetadata(serviceMi, serviceUrl)
	assert.Nil(t, err)

	subMi := newSubscribeMetadataIdentifier()
	urlList := make([]string, 0, 1)
	urlList = append(urlList, serviceUrl.String())
	urls, _ := json.Marshal(urlList)
	err = metadataReport.SaveSubscribedData(subMi, string(urls))
	assert.Nil(t, err)

	serviceUrl, _ = common.NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?interface=com.ikurento.user.UserProvider&group=&version=2.6.0")
	metadataInfo := common.NewMetadataInfo(subMi.Application, "", map[string]*common.ServiceInfo{
		"com.ikurento.user.UserProvider": common.NewServiceInfoWithURL(serviceUrl),
	})
	err = metadataReport.RemoveServiceMetadata(serviceMi)
	assert.Nil(t, err)
	err = metadataReport.PublishAppMetadata(subMi, metadataInfo)
	assert.Nil(t, err)

	mdInfo, err := metadataReport.GetAppMetadata(subMi)
	assert.Nil(t, err)
	assert.Equal(t, metadataInfo.App, mdInfo.App)
	assert.Equal(t, metadataInfo.Revision, mdInfo.Revision)
	assert.Equal(t, 1, len(mdInfo.Services))
	assert.NotNil(t, metadataInfo.Services["com.ikurento.user.UserProvider"])

	e.Close()
}

func TestEtcdMetadataReport_ServiceAppMapping(t *testing.T) {
	e := initEtcd(t)
	url, err := common.NewURL("registry://127.0.0.1:2379", common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER)))
	if err != nil {
		t.Fatal(err)
	}
	metadataReportFactory := &etcdMetadataReportFactory{}
	metadataReport := metadataReportFactory.CreateMetadataReport(url)
	assert.NotNil(t, metadataReport)

	_, err = metadataReport.GetServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping")
	assert.NotNil(t, err)

	err = metadataReport.RegisterServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping", "demo_provider")
	assert.Nil(t, err)
	set, err := metadataReport.GetServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping")
	assert.Nil(t, err)
	assert.Equal(t, 1, set.Size())

	err = metadataReport.RegisterServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping", "demo_provider2")
	assert.Nil(t, err)
	err = metadataReport.RegisterServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping", "demo_provider")
	assert.Nil(t, err)
	set, err = metadataReport.GetServiceAppMapping("com.apache.dubbo.sample.basic.IGreeter", "mapping")
	assert.Nil(t, err)
	assert.Equal(t, 2, set.Size())

	e.Close()
}

func newSubscribeMetadataIdentifier() *identifier.SubscriberMetadataIdentifier {
	return &identifier.SubscriberMetadataIdentifier{
		Revision: "subscribe",
		BaseApplicationMetadataIdentifier: identifier.BaseApplicationMetadataIdentifier{
			Application: "provider",
		},
	}
}

func newServiceMetadataIdentifier() *identifier.ServiceMetadataIdentifier {
	return &identifier.ServiceMetadataIdentifier{
		Protocol: "nacos",
		Revision: "a",
		BaseMetadataIdentifier: identifier.BaseMetadataIdentifier{
			ServiceInterface: "com.test.MyTest",
			Version:          "1.0.0",
			Group:            "test_group",
			Side:             "service",
		},
	}
}

func newMetadataIdentifier(side string) *identifier.MetadataIdentifier {
	return &identifier.MetadataIdentifier{
		Application: "test",
		BaseMetadataIdentifier: identifier.BaseMetadataIdentifier{
			ServiceInterface: "com.test.MyTest",
			Version:          "1.0.0",
			Group:            "test_group",
			Side:             side,
		},
	}
}

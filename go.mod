// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

module github.com/otcshare/edgecontroller

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/grpc-ecosystem/grpc-gateway v1.8.4
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/otcshare/common/log v0.0.0-20190819124907-99dcb8b50ed8
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c // indirect
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b
	google.golang.org/grpc v1.19.0
	gopkg.in/square/go-jose.v2 v2.3.1
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/client-go v0.0.0-20190501104856-ef81ee0960bf
	k8s.io/utils v0.0.0-20190520173318-324c5df7d3f0 // indirect
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190226215855-775f8194d0f9

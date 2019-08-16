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

//nolint:lll
//go:generate protoc -I../../schema/pb -I../../../grpc-ecosystem/grpc-gateway -I../../../grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc,paths=source_relative:auth ../../schema/pb/auth.proto

//nolint:lll
//go:generate protoc -I../../schema/pb -I../../../grpc-ecosystem/grpc-gateway -I../../../grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc,paths=source_relative:ela ../../schema/pb/ela.proto

//nolint:lll
//go:generate protoc -I../../schema/pb -I../../../grpc-ecosystem/grpc-gateway -I../../../grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc,paths=source_relative,Mela.proto=github.com/open-ness/edgecontroller/pb/ela:eva ../../schema/pb/eva.proto

package pb

// Copyright 2019 Intel Corporation. All rights reserved.
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

module kube-rsu

require (
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.0.0-20191016110246-af539daaa43a
	k8s.io/apimachinery v0.0.0-20191004115701-31ade1b30762
	k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6 // indirect
)

replace rsu v0.0.0 => ./cmd

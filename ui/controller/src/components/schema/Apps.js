// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
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

{/* 

GET /apps
List of applications.

*/}

export default {
  schema: {
    type: "object",
    title: "Node Apps",
    properties: {
      apps: {
        type: "array",
        title: "Apps",
        items: {
          type: "object",
          title: "App",
          required: [
            "type",
            "name",
            "version",
            "vendor"
          ],
          properties: {
            id: {
              type: "string",
              title: "ID",
              readonly: true
            },
            type: {
              type: "string",
              title: "Type",
              enum: [
                "container",
                "vm"
              ]
            },
            name: {
              type: "string",
              title: "Name"
            },
            version: {
              type: "string",
              title: "Version"
            },
            vendor: {
              type: "string",
              title: "Vendor"
            },
            description: {
              type: "string",
              title: "Description"
            }
          }
        }
      }
    },
  },
  form: [
    "*"
  ]
};

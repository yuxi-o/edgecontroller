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

GET /nodes
List of Nodes.

*/}

export default {
  schema: {
    type: "object",
    title: "Nodes",
    properties: {
      nodes: {
        type: "array",
        title: "Nodes",
        items: {
          type: "object",
          title: "Node",
          required: [
            "name",
            "serial",
            "location"
          ],
          properties: {
            id: {
              type: "string",
              title: "ID",
              readonly: true
            },
            name: {
              type: "string",
              title: "Name"
            },
            location: {
              type: "string",
              title: "Location"
            },
            serial: {
              type: "string",
              title: "Serial"
            }
          }
        }
      }
    }
  },
  form: [
    "*"
  ]
};

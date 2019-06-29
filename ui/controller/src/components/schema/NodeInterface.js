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

export default {
  schema: {
    type: "object",
    title: "Node Interface",
    properties: {
      id: {
        title: "ID",
        type: "string",
        readonly: true
      },
      description: {
        title: "Description",
        type: "string"
      },
      driver: {
        title: "Driver",
        type: "string",
        enum: [
          "kernel",
          "userspace"
        ]
      },
      type: {
        title: "Type",
        type: "string",
        enum: [
          "none",
          "upstream",
          "downstream",
          "bidirectional",
          "breakout"
        ]
      },
      mac_address: {
        title: "Mac Address",
        type: "string",
        readOnly: true
      },
      vlan: {
        title: "Vlan",
        type: "number",
        readOnly: true
      },
      zones: {
        title: "Zones",
        type: "array",
        items: {
          title: "Private Traffic",
          type: "string"
        }
      },
      fallback_interface: {
        title: "Fallback Interface",
        type: "string"
      }
    }
  },
  form: [
    "*"
  ]
};

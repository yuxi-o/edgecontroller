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
    title: "App",
    required: [
      "type",
      "name",
      "version",
      "vendor",
      "cores",
      "memory",
      "source"
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
      },
      cores: {
        type: "number",
        title: "Cores",
        minimum: 1,
        maximum: 8,
      },
      memory: {
        type: "number",
        title: "Memory (in MB)",
        minimum: 1,
        maximum: 16384,
      },
      ports: {
        type: "array",
        title: "Ports",
        items: {
          type: "object",
          title: "Port",
          properties: {
            port: {
              type: "number",
              title: "Port",
              minimum: 1,
              maximum: 65535,
            },
            protocol: {
              type: "string",
              title: "Protocol",
            },
          },
        },
      },
      source: {
        type: "string",
        title: "Source",
      },
      epafeatures: {
        type: "array",
        title: "EPA Features",
        items: {
          type: "object",
          title: "EPA Feature",
          properties: {
            key: {
              type: "string",
              title: "EPA Feature Key",
            },
            value: {
              type: "string",
              title: "EPA Feature Value",
            },
          },
        },
      }
    }
  },
  form: [
    "*"
  ]
};

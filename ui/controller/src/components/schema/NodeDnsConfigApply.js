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
    title: "DNS Config",
    required: [
      "name"
    ],
    properties: {
      name: {
        title: "Name",
        type: "string"
      },
      records: {
        type: "object",
        title: "Records",
        properties: {
          a: {
            type: "array",
            title: "A Records",
            items: {
              type: "object",
              title: "A Record",
              properties: {
                name: {
                  "title": "Name",
                  "type": "string"
                },
                description: {
                  "title": "Description",
                  "type": "string"
                },
                alias: {
                  "title": "Alias",
                  "type": "boolean"
                },
                values: {
                  "title": "Values",
                  type: "array",
                  items: {
                    "type": "string"
                  }
                }
              }
            }
          }
        }
      },
      // configurations: {
      //   type: "object",
      //   properties: {
      //     forwarders: {
      //       type: "array",
      //       items: {
      //         type: "object",
      //         properties: {
      //           name: {
      //             type: "string",
      //           },
      //           description: {
      //             type: "string",
      //           },
      //           value: {
      //             type: "string",
      //           },
      //         },
      //       },
      //     },
      //   },
      // },
    }
  },
  form: [
    "*"
  ]
};

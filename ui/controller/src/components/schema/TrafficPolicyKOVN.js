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

const cidrField = {
  type: 'string',
  format: 'regex',
  pattern:
    '^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\\/(\\d|[1-2]\\d|3[0-2]))$',
  validationMessage: 'Please, enter an IP in a valid CIDR notation.',
};

const IPBlocks = {
  items: {
    title: 'IP Block',
    type: 'object',
    properties: {
      cidr: {
        title: 'CIDR',
        ...cidrField,
      },
      except: {
        title: 'Excepts',
        type: 'array',
        items: {
          ...cidrField,
        },
      },
    },
  },
};

const Ports = {
  ports: {
    title: 'Ports',
    type: 'array',
    items: {
      type: 'object',
      properties: {
        port: {
          title: 'Port number',
          type: 'number',
          minimum: 0,
          maximum: 65535,
        },
        protocol: {
          title: 'Protocol',
          type: 'string',
          enum: ['tcp', 'udp', 'sctp'],
        },
      },
      dependencies: {
        port: ['protocol'],
        protocol: ['port'],
      },
    },
  },
};

export default {
  schema: {
    type: 'object',
    title: 'Traffic Policy',
    required: ['name'],
    properties: {
      id: {
        title: 'ID',
        type: 'string',
        readOnly: true,
      },
      name: {
        title: 'Name',
        type: 'string',
      },
      ingress_rules: {
        type: 'array',
        title: 'Ingress Rules',
        items: {
          type: 'object',
          title: 'New Ingress Rule',
          properties: {
            description: {
              title: 'Description',
              type: 'string',
            },
            from: {
              title: 'From',
              type: 'array',
              ...IPBlocks,
            },
            ...Ports,
          },
        },
      },
      egress_rules: {
        type: 'array',
        title: 'Egress Rules',
        items: {
          type: 'object',
          title: 'New Egress Rule',
          properties: {
            description: {
              title: 'Description',
              type: 'string',
            },
            to: {
              title: 'To',
              type: 'array',
              ...IPBlocks,
            },
            ...Ports,
          },
        },
      },
    },
  },
  form: ['*'],
};

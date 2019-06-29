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

export default [
  {
    id: "1",
    uuid: "a8a25f02-b20c-4a29-a42c-46087c64fb66",
    function: "SGWU",
    config: {
      sxa: {
        cp_ip_address: "1.1.1.1",
        up_ip_address: "1.1.1.1",
      },
      sxb: {
        cp_ip_address: "1.1.1.1",
        up_ip_address: "1.1.1.1",
      },
      s1u: {
        up_ip_address: "1.1.1.1",
      },
      s5u_sgw: {
        up_ip_address: "1.1.1.1",
      },
      s5u_pgw: {
        up_ip_address: "1.1.1.1",
      },
      sgi: {
        up_ip_address: "1.1.1.1",
      },
      breakout: [
        {
          up_ip_address: "1.1.1.1",
        }
      ],
      dns: [
        {
          up_ip_address: "1.1.1.1",
        }
      ],
    },
    selectors: [
      {
        id: "1",
        network: {
          mcc: "",
          mnc: "",
        },
        uli: {
          tai: {
            "tac": 0
          },
          ecgi: {
            "eci": 0
          }
        },
        pdn: {
          "apns": [
            ""
          ]
        },
      },
    ],
    entitlements: [
      {
        id: "1",
        apns: [
          "1.1.1.1"
        ],
        imsis: [
          {
            begin: "",
            end: ""
          }
        ]
      }
    ],
  },
];

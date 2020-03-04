/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

const pfdForm = [
  "*",
];

const pfdAppForm = [
  "*",
];

const pfdAppFormSchema = {
  type: "object",
  title: "Af Application",
  properties: {
    externalAppID: {
      title: "AF application ID",
      type: "string",
    },
    allowedDelay: {
      title: "Allowed delay (in seconds)",
      type: "number",
    },
    cachingTime: {
      title: "Caching time (in seconds)",
      type: "number",
    },
    pfds: {
      title: "PFD's",
      type: "array",
      minItems: 1,
      maxItems: 10,
      items: {
        title: "",
        type: "object",
        properties: { 
          pfd: {
            notitle: "true",
            type: "object",
            properties: {
             pfdID: {
              title: "PFD ID",
              type: "string",
            },
            flowDescType: {
              title: "PFD Rule Type",
              type: "string",
              enum: [
                "Flow Description",
                "URL",
                "Domain Name",
              ],
            },
            flowDescValue: {
              title: "Values",
              type: "array",
              items: {
                title: "PFD Rule value",
                type: "string",
              },
            },
          }
        }
      }
    }
  }
}
};


const pfdFormSchema = {
  type: "object",
  title: "Packet Flow Descriptor",
  properties: {
    pfdDatas:{
      title: "PFD Transactions",
      type: "array",
      items: {
        title: "",
        type: "object",
        properties: {
          apps: {
            title: "Af Application",
            type: "object",
            properties: {
              externalAppID: {
                title: "AF application ID",
                type: "string",
              },
              allowedDelay: {
                title: "Allowed delay (in seconds)",
                type: "number",
              },
              cachingTime: {
                title: "Caching time (in seconds)",
                type: "number",
              },
              pfds: {
                title: "PFD's",
                type: "array",
                minItems: 1,
                maxItems: 10,
                items: {
                  title: "",
                  type: "object",
                  properties: { 
                    pfd: {
                      notitle: "true",
                      type: "object",
                      properties: {
                       pfdID: {
                        title: "PFD ID",
                        type: "string",
                      },
                      flowDescType: {
                        title: "PFD Rule Type",
                        type: "string",
                        enum: [
                          "Flow Description",
                          "URL",
                          "Domain Name",
                        ],
                      },
                      flowDescValue: {
                        title: "Values",
                        type: "array",
                        items: {
                          title: "PFD Rule value",
                          type: "string",
                        },
                      },
                    }
                  } // pfd end
                }
              }
            } // pfds end
          }
        } // apps end
      }
    }
  } // pfdDatas end
 }
};


export {
  pfdAppForm,
  pfdForm,
  pfdAppFormSchema,
  pfdFormSchema,
};

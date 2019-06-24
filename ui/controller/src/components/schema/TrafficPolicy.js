const IPFilter = {
  ip_filter: {
    title: "IP Filter",
    type: "object",
    properties: {
      address: {
        title: "IP Address",
        type: "string"
      },
      mask: {
        title: "Mask",
        type: "number",
        minimum: 0,
        maximum: 128,
      },
      begin_port: {
        title: "Begin Port",
        type: "number",
        minimum: 0,
        maximum: 65535,
      },
      end_port: {
        title: "End Port",
        type: "number",
        minimum: 0,
        maximum: 65535,
      },
      protocol: {
        title: "Protocol",
        type: "string",
        enum: [
          "all",
          "tcp",
          "udp",
          "icmp",
          "sctp"
        ]
      },
    },
  },
};

const MACFilter = {
  mac_filter: {
    title: "MAC Filter",
    type: "object",
    properties: {
      mac_addresses: {
        title: "MAC Addresses",
        type: "array",
        items: {
          title: "MAC Address",
          type: "string"
        }
      }
    }
  },
};

const GTPFilter = {
  gtp_filter: {
    title: "GTP Filter",
    type: "object",
    properties: {
      address: {
        title: "Address",
        type: "string"
      },
      mask: {
        title: "Mask",
        type: "number",
        minimum: 0,
        maximum: 128,
      },
      imsis: {
        title: "IMSIs",
        type: "array",
        items: {
          title: "IMSI",
          type: "number"
        }
      }
    }
  }
};

export default {
  schema: {
    type: "object",
    title: "Traffic Policy",
    required: [
      "name",
      "traffic_rules"
    ],
    properties: {
      id: {
        title: "ID",
        type: "string",
        readOnly: true
      },
      name: {
        title: "Name",
        type: "string",
        required: true
      },
      traffic_rules: {
        type: "array",
        title: "Traffic Rules",
        items: {
          type: "object",
          title: "Traffic Rule",
          properties: {
            description: {
              title: "Description",
              type: "string",
            },
            priority: {
              title: "Priority",
              type: "number",
              min: 0,
              max: 65535,
              required: true,
            },
            source: {
              title: "Source",
              type: "object",
              properties: {
                description: {
                  title: "Description",
                  type: "string"
                },
                ...MACFilter,
                ...IPFilter,
                ...GTPFilter,
              }
            },
            destination: {
              title: "Destination",
              type: "object",
              properties: {
                description: {
                  title: "Description",
                  type: "string"
                },
                ...MACFilter,
                ...IPFilter,
                ...GTPFilter,
              }
            },
            target: {
              title: "Target",
              type: "object",
              properties: {
                description: {
                  title: "Description",
                  type: "string",
                },
                action: {
                  title: "Action",
                  type: "string",
                  default: "accept",
                  enum: [
                    "accept",
                    "reject",
                    "drop"
                  ]
                },
                mac_modifier: {
                  title: "MAC Modifier",
                  type: "string"
                },
                ip_modifier: {
                  title: "IP Modifier",
                  type: "object",
                  properties: {
                    address: {
                      title: "IP Address",
                      type: "string"
                    },
                    port: {
                      title: "Port",
                      type: "number"
                    },
                  },
                },
              }
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

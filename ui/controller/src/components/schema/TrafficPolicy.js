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
        type: "string"
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
              type: "string"
            },
            priority: {
              title: "Priority",
              type: "number"
            },
            source: {
              title: "Source",
              type: "object",
              properties: {
                description: {
                  title: "Description",
                  type: "string"
                },
                mac_filter: {
                  type: "object",
                  title: "MAC Filter",
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
                ip_filter: {
                  title: "IP Filter",
                  type: "object",
                  properties: {
                    mask: {
                      title: "Mask",
                      type: "number"
                    },
                    begin_port: {
                      title: "Begin Port",
                      type: "number"
                    },
                    end_port: {
                      title: "End Port",
                      type: "number"
                    },
                    protocol: {
                      title: "Protocol",
                      type: "string"
                    },
                  },
                },
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
                      type: "number"
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
                ip_filter: {
                  title: "IP Filter",
                  type: "object",
                  properties: {
                    mask: {
                      title: "Mask",
                      type: "number"
                    },
                    begin_port: {
                      title: "Begin Port",
                      type: "number"
                    },
                    end_port: {
                      title: "End Port",
                      type: "number"
                    },
                    protocol: {
                      title: "Protocol",
                      type: "string"
                    },
                  },
                },
                gtp_filter: {
                  title: "GTP Filter",
                  type: "object",
                  properties: {
                    mask: {
                      title: "Mask",
                      type: "number"
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
              }
            },
            target: {
              title: "Target",
              type: "object",
              properties: {
                description: {
                  title: "Description",
                  type: "string"
                },
                action: {
                  title: "Action",
                  type: "string",
                  enum: [
                    "accept",
                    "reject",
                    "drop"
                  ]
                },
                mac_modifier: {
                  title: "MAC Modifier",
                  type: "string"
                }
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

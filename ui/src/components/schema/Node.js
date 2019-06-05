export default {
  schema: {
    type: "object",
    title: "Node",
    required: [
      "name",
      "serial",
      "location",
    ],
    properties: {
      id: {
        type: "string",
        title: "ID",
      },
      name: {
        type: "string",
        title: "Name",
      },
      location: {
        type: "string",
        title: "Location",
      },
      serial: {
        type: "string",
        title: "Serial",
      },
      interfaces: {
        type: "array",
        title: "Interfaces",
        required: [
          "id",
          "driver",
          "type",
          "vlan",
        ],
        items: {
          type: "object",
          title: "Interface",
          properties: {
            id: {
              type: "string",
              title: "ID",
            },
            description: {
              type: "string",
              title: "Description",
            },
            driver: {
              type: "string",
              title: "Driver",
              enum: [
                "kernel",
                "userspace",
              ],
            },
            type: {
              type: "string",
              title: "Type",
              enum: [
                "none",
                "upstream",
                "downstream",
                "bidirectional",
                "breakout",
              ],
            },
            mac_address: {
              type: "string",
              title: "MAC Address",
              readonly: true,
            },
            vlan: {
              type: "number",
              title: "VLAN",
            },
            zones: {
              type: "array",
              title: "Zones",
              items: {
                type: "string",
                title: "Private Traffic",
              },
            },
            fallback_interfaces: {
              type: "string",
              title: "Fallback Interfaces",
            },
          },
        },
      },
      source: {
        type: "string",
        title: "Source",
      }
    },
  },
  form: [
    "*",
  ],
};

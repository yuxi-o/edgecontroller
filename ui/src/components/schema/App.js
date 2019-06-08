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
      },
      memory: {
        type: "number",
        title: "Memory",
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
      }
    }
  },
  form: [
    "*"
  ]
};

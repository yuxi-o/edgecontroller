export default {
  schema: {
    type: "object",
    title: "DNS Config",
    required: [
      "name"
    ],
    properties: {
      id: {
        title: "ID",
        type: "string",
        readonly: true
      },
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
      dns_forwarders: {
        type: "array",
        title: "DNS Forwarders",
        items: {
          title: "DNS Forwarder",
          type: "string"
        }
      }
    }
  },
  form: [
    "*"
  ]
};

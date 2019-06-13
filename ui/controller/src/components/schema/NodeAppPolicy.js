{/* 

GET /nodes/{node_id}/apps/{app_id}/policy
Gets the traffic policy ID associated with a node application from the controller.

*/}

export default {
  schema: {
    type: "object",
    title: "Node App Policy",
    properties: {
      id: {
        title: "ID",
        type: "string",
        readOnly: true
      }
    }
  },
  form: [
    "*"
  ]
};

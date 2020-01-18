var { grpc } = require("@improbable-eng/grpc-web");

if (typeof window === 'undefined') {
    var { NodeHttpTransport } = require("@improbable-eng/grpc-web-node-http-transport");
    grpc.setDefaultTransport(NodeHttpTransport());
}

module.exports.grpc = grpc;

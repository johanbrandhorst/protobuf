const resolve = require("path").resolve;

module.exports = {
  entry: "./node_modules/@improbable-eng/grpc-web/dist/grpc-web-client.js",
  output: {
    path: resolve("."),
    filename: "grpc.inc.js",
    libraryTarget: "this"
  },
  optimization: {
    minimize: true
  },
  mode: "production"
};

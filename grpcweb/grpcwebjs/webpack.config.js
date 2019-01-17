const resolve = require("path").resolve;

module.exports = {
  entry: "./node_modules/grpc-web-client/dist/index.js",
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

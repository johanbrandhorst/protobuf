const resolve = require("path").resolve;

module.exports = {
  entry: "./index.js",
  output: {
    path: resolve("."),
    filename: "grpc.inc.js",
    libraryTarget: "this"
  },
  target: 'node',
  optimization: {
    minimize: true
  },
  mode: "production"
};

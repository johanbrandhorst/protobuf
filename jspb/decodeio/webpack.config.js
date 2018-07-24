const webpack = require("webpack");

module.exports = {
    entry: {
        decodeio: [
            "./node_modules/protobufjs/minimal.js",
            "./node_modules/long/index.js"
        ]
    },
    output: {
        filename: "decodeio.inc.js",
        libraryTarget: "this",
        path: __dirname,
    },
    optimization: {
        minimize: false
    },
    mode: "production",
};

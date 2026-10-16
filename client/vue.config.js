const { defineConfig } = require("@vue/cli-service");
const packageJson = require("./package.json");

// Vue CLI exposes every VUE_APP_* variable to the client; fall back to package.json
process.env.VUE_APP_VERSION = process.env.VUE_APP_VERSION || packageJson.version;

module.exports = defineConfig({
  transpileDependencies: true,
});

module.exports = {
  devServer: {
    proxy: "http://app:8080",
  },
  configureWebpack: {
    devtool: "source-map",
  },
};

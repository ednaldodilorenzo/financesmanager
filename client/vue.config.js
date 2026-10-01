const { defineConfig } = require("@vue/cli-service");
const packageJson = require("./package.json");

module.exports = defineConfig({
  transpileDependencies: true,
});

module.exports = {
  devServer: {
    proxy: "http://app:8080",
  },
  configureWebpack: {
    devtool: "source-map",
    plugins: [
      new (require("webpack").DefinePlugin)({
        "process.env.VUE_APP_VERSION": JSON.stringify(packageJson.version),
      }),
    ],
  },
};

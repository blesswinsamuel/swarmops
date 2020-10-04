module.exports = {
  devServer: {
    proxy: `http://localhost:${process.env.API_SERVER_PORT || "8080"}`,
  },
};

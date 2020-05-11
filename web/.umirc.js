// ref: https://umijs.org/config/
export default {
  treeShaking: true,
  plugins: [
    // ref: https://umijs.org/plugin/umi-plugin-react.html
    ['umi-plugin-react', {
      antd: true,
      dva: true,
      dynamicImport: false,
      title: 'demoweb',
      dll: false,
      
      routes: {
        exclude: [
          /models\//,
          /services\//,
          /model\.(t|j)sx?$/,
          /service\.(t|j)sx?$/,
          /components\//,
        ],
      },
    }],
  ],
  hash: true,
  "proxy": {
    "/admin": {
      // "target": "http://10.101.1.152:9090/",
      "target": "http://localhost:8080",
      request_timeout: 12000,
      "changeOrigin": true,
      // pathRewrite: {
        // '^/api': ''
      // }
    },
  },
}

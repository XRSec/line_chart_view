const {defineConfig} = require('@vue/cli-service')
module.exports = defineConfig({
    transpileDependencies: true,
    pages: {
        index: {
            title: "网络攻击流量墙",
            entry: "src/main.js",
        },
    },
    publicPath: process.env.NODE_ENV === 'production' ? '/home' : '/',
    devServer: {
        port: 8080,
        proxy: process.env.NODE_ENV === 'production' ? null : {
            '/api': {
                target: 'http://localhost:8081',
                changeOrigin: true,
            }
        }
    },
    outputDir: '../server/static',
})

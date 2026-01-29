// frontend/config-overrides.js
const webpack = require('webpack');

module.exports = function override(config) {
  // Node 核心模块 polyfill
  config.resolve.fallback = {
    ...(config.resolve.fallback || {}),
    crypto: require.resolve('crypto-browserify'),
    stream: require.resolve('stream-browserify'),
    http: require.resolve('stream-http'),
    https: require.resolve('https-browserify'),
    zlib: require.resolve('browserify-zlib'),
    url: require.resolve('url/'),
    vm: require.resolve('vm-browserify'),
    process: require.resolve('process/browser.js'),
  };

  // 兼容像 axios / walletconnect 里直接写 'process/browser' 的引用
  config.resolve.alias = {
    ...(config.resolve.alias || {}),
    'process/browser': require.resolve('process/browser.js'),
  };

  config.plugins = [
    ...(config.plugins || []),
    new webpack.ProvidePlugin({
      process: 'process/browser',
      Buffer: ['buffer', 'Buffer'],
    }),
  ];

  return config;
};
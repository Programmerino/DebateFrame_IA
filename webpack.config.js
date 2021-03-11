const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const WorkboxPlugin = require('workbox-webpack-plugin');
require('wasm-loader');

module.exports = {
  mode: 'development',
  devtool: 'inline-source-map',
  entry: {
    app: './static/index.js',
  },
  plugins: [
    new WorkboxPlugin.GenerateSW({
      //include: [/\.html$/, /\.js$/, /\.wasm$/, /\.woff2$/, /\.svg$/],
      importWorkboxFrom: 'local',
      importsDirectory: 'wb-assets',
      navigateFallback: '/dist/',
      cacheId: 'debateframe',
      clientsClaim: true,
      skipWaiting: true,
      maximumFileSizeToCacheInBytes: 15 * 1024 * 1024,
      globDirectory: "dist",
      globPatterns: ["client.wasm"],
    }),
  ],
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  resolve: {
    alias: {
      'parchment': path.resolve(__dirname, 'node_modules/parchment/src/parchment.ts'),
      'quill$': path.resolve(__dirname, 'node_modules/quill/quill.js'),
    },
    extensions: ['.js', '.ts', '.svg']
  },
  module: {
    rules: [
      {
        test: /\.wasm$/,
        loaders: ['wasm-loader']
      },
      { test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "url-loader?limit=10000&mimetype=application/font-woff" },
      { test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "file-loader" },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.html$/,
        use: ['file-loader?name=[name].[ext]', 'extract-loader', 'html-loader'],
      },
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [['@babel/preset-env', {
              "useBuiltIns": "usage"
            }]],
            plugins: [
              "@babel/plugin-syntax-dynamic-import",
              "@babel/proposal-class-properties"
            ],
          }
        }
      },
      {
        test: /\.less$/,
        use: [
          {
            loader: "style-loader"
          },
          {
            loader: "css-loader"
          },
          {
            loader: "less-loader"
          }
        ]
      },
      {
        test: /\.ts$/,
        use: [{
          loader: 'ts-loader',
          options: {
            compilerOptions: {
              declaration: false,
              target: 'es5',
              module: 'commonjs'
            },
            transpileOnly: true
          }
        }]
      },
      {
        test: /\.svg$/,
        use: [{
          loader: 'html-loader',
          options: {
            minimize: true
          }
        }]
      },
    ],
  },
  node: {
    fs: 'empty',
  },
};
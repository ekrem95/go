const webpack = require('webpack');
const path = require('path');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');

const production = process.env.NODE_ENV === 'production';
const jsDev = [];
const jsProduction = [
  new UglifyJSPlugin(),
];

const APP_DIR = path.resolve(__dirname, 'app');
const BUILT_DIR = path.resolve(__dirname, 'public');

const config = {
  entry: APP_DIR + '/app.js',
  output: {
    path: BUILT_DIR,
    filename: 'bundle.js',
  },
  devtool: production ? '' : 'inline-sourcemap',
  module: {
    loaders: [
      {
        test: /\.js$/,
        include: APP_DIR,
        exclude: /node_modules/,
        loader: 'babel-loader',
        query: {
          presets: ['react'],
        },
      },
      {
        test: /\.scss$/,
        use: [
          { loader: 'style-loader' },
          { loader: 'css-loader' },
          { loader: 'sass-loader' },
        ],
      },
    ],
  },
  plugins: production ? jsProduction : jsDev,
};

module.exports = config;

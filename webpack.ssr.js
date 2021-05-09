const path = require("path");
const webpack = require("webpack");
const TerserPlugin = require("terser-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

const pathToLibrary = (name) =>
    path.resolve(__dirname, `./node_modules/${name}`);

module.exports = {
    entry: "./frontend/ssr.js",
    target: "node",
    mode: "production",
    output: {
        path: path.resolve(__dirname, "build"),
        filename: "ssr.js",
        libraryTarget: "commonjs2",
    },
    resolve: {
        alias: {
            "apollo-utilities": pathToLibrary("apollo-utilities"),
            "core-js": pathToLibrary("core-js"),
            "prop-types": pathToLibrary("prop-types"),
            "hoist-non-react-statics": pathToLibrary("hoist-non-react-statics"),
            "@emotion/memoize": pathToLibrary("@emotion/memoize"),
            "@emotion/stylis": pathToLibrary("@emotion/stylis"),
            "@emotion/hash": pathToLibrary("@emotion/hash"),
        },
    },
    module: {
        rules: [
            {
                test: /\.js?$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                },
            },
            {
                test: /\.svg$/,
                use: ["@svgr/webpack"],
            },
            {
                test: /\.less$/i,
                loaders: [
                    MiniCssExtractPlugin.loader,
                    {
                        loader: "css-loader",
                        options: {
                            modules: false,
                        },
                    },
                    {
                        loader: "less-loader",
                        options: {
                            lessOptions: {
                                javascriptEnabled: true,
                            },
                        },
                    },
                ],
            },
        ],
    },
    plugins: [
        new MiniCssExtractPlugin(),
        // new webpack.IgnorePlugin(/canvas/), // jsdom
        new webpack.DefinePlugin({
            "process.env.SERVER": JSON.stringify(true),
        }),
    ],
    optimization: {
        minimize: false,
        minimizer: [
            new CssMinimizerPlugin({
                minimizerOptions: {
                    preset: [
                        // "advanced",
                        // "default",
                        "advanced",
                        {
                            discardComments: { removeAll: true },
                        },
                    ],
                },
            }),
            new TerserPlugin({
                terserOptions: {
                    ie8: false,
                    output: {
                        comments: false,
                    },
                },
                sourceMap: true,
                extractComments: false,
            }),
        ],
    },
    stats: {
        maxModules: Number.MAX_VALUE,
    },
};

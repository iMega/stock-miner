const webpack = require("webpack");
const TerserPlugin = require("terser-webpack-plugin");
// const CopyPlugin = require("copy-webpack-plugin");

const multipleManifest = (manifestNames) =>
    manifestNames.map(
        (name) =>
            new webpack.DllReferencePlugin({
                context: __dirname,
                manifest: require(`./manifest.${name}.json`),
                name: name,
            })
    );

const version = process.env.VERSION.split("/")[2];

module.exports = {
    entry: {
        client: "./frontend/index.js",
    },
    target: "web",
    mode: "production",
    output: {
        filename: `[name]-${version}.js`,
        path: __dirname + "/build",
        publicPath: "/",
    },
    externals: {
        react: "React",
        "react-dom": "ReactDOM",
        // "cross-fetch": "cross-fetch",
        // "@apollo/client": "@apollo/client",
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx|mjs)$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                },
            },
            {
                test: /\.svg$/,
                use: ["@svgr/webpack"],
            },
        ],
    },
    plugins: [
        // new BundleAnalyzerPlugin(),
        // ...multipleManifest([
        // "vendor",
        // "core",
        // "react",
        // "graphql",
        // "apollo",
        // "uuid",
        // "emotion"
        // ]),
        // new CopyPlugin({
        //     patterns: [{ from: "src/icons" }],
        // }),
        new webpack.DefinePlugin({
            "process.env.STORYBOOK_GRAPHQL_HOST": JSON.stringify(
                process.env.STORYBOOK_GRAPHQL_HOST || "http://localhost/query"
            ),
            WS_HOST: JSON.stringify(
                process.env.WS_HOST || "ws://localhost/query"
            ),
        }),
        new webpack.BannerPlugin(`ver.: ${version}`),
    ],
    optimization: {
        minimize: true,
        minimizer: [
            new TerserPlugin({
                terserOptions: {
                    ie8: false,
                    // output: {
                    //     comments: false,
                    // },
                },
                sourceMap: true,
                extractComments: false,
            }),
        ],
        // splitChunks: {
        //     chunks: "all",
        // },
    },
    stats: {
        maxModules: Number.MAX_VALUE,
    },
};

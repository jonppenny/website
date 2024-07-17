const path = require('path');
const webpack = require('webpack');
const {merge} = require('webpack-merge');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const TerserPlugin = require("terser-webpack-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = (env) => {
    const rootPath = process.cwd();
    const isProduction = env.production || false;

    const userConfig = require('./config');
    const config = merge({
        paths: {
            root: rootPath,
            assets: path.join(rootPath, "/resources/"),
            dist: path.join(rootPath, "/dist/")
        },
        enabled: {
            sourceMaps: !isProduction,
            optimize: isProduction,
            // cacheBusting: isProduction
        },
    }, userConfig);

    let webpackConfig = {
        entry: {
            main: [
                './resources/scripts/app.js',
                './resources/styles/app.scss'
            ]
        },
        context: config.paths.root,
        output: {
            filename: 'js/[name].js',
            path: config.paths.dist,
            pathinfo: !isProduction,
            clean: true,
        },
        devtool: (config.enabled.sourceMaps ? 'eval-source-map' : undefined),
        module: {
            rules: [
                {
                    enforce: 'pre',
                    test: /\.(js|s?[ca]ss)$/,
                    include: config.paths.assets,
                    loader: 'import-glob',
                },
                {
                    test: /\.js$/,
                    exclude: [/node_modules(?![/|\\](bootstrap|foundation-sites))/],
                    use: {
                        loader: "babel-loader",
                        options: {
                            presets: ['@babel/preset-env'],
                        }
                    }
                },
                {
                    test: /\.css$/,
                    include: config.paths.assets,
                    use: [
                        {
                            loader: MiniCssExtractPlugin.loader,
                        },
                        {
                            loader: 'css-loader',
                            options: {
                                sourceMap: config.enabled.sourceMaps
                            }
                        },
                        {
                            loader: 'postcss-loader',
                            options: {
                                sourceMap: config.enabled.sourceMaps
                            }
                        },
                    ]
                },
                {
                    test: /\.scss$/,
                    include: config.paths.assets,
                    use: [
                        {
                            loader: MiniCssExtractPlugin.loader,
                        },
                        {
                            loader: 'css-loader',
                            options: {
                                sourceMap: config.enabled.sourceMaps
                            }
                        },
                        {
                            loader: 'postcss-loader',
                            options: {
                                sourceMap: config.enabled.sourceMaps
                            }
                        },
                        {
                            loader: "sass-loader",
                            options: {
                                sourceMap: config.enabled.sourceMaps
                            }
                        }
                    ]
                },
                /*{
                  test: /\.(jpg|jpe?g|gif|png|svg)$i/,
                  type: "asset/resource",
                  generator: {
                    filename: 'images/[name].[hash][ext][query]'
                  }
                },
                {
                  test: /\.(woff|woff2|eot|ttf|otf)$/i,
                  type: "asset/resource",
                  generator: {
                    filename: 'fonts/[name].[hash][ext][query]'
                  }
                }*/
            ]
        },
        resolve: {
            modules: [
                "node_modules"
            ],
            enforceExtension: false
        },
        externals: {
            jquery: 'jQuery',
        },
        plugins: [
            new MiniCssExtractPlugin({
                filename: "css/[name].css"
            }),
            new webpack.ProvidePlugin({
                $: 'jquery',
                jQuery: 'jquery',
                'window.jQuery': 'jquery',
            }),
            new CopyPlugin({
                patterns: [
                    { from: "resources/fonts", to: "fonts" },
                    {from: "resources/images", to: "images"},
                ],
            }),
        ],
        stats: {
            errorDetails: true,
            children: true
        }
    };

    if (config.enabled.optimize) {
        webpackConfig = merge(webpackConfig, {
            optimization: {
                minimize: true,
                minimizer: [
                    new TerserPlugin(),
                    new CssMinimizerPlugin(),
                ],
            },
        });
    }

    return webpackConfig;
};

const fs = require("fs");
const { PurgeCSS } = require("purgecss");
import React from "react";
import ReactDOM from "react-dom/server";
import { renderToStringWithData } from "@apollo/client/react/ssr";
import { CacheProvider } from "@emotion/react";
import createCache from "@emotion/cache";
import createEmotionServer from "create-emotion-server";
import "antd/dist/antd.less";

import MainPage from "./Main.page";
import SignInPage from "./Signin.page";

const key = "custom";
const cache = createCache({ key });
const { extractCritical } = createEmotionServer(cache);

const r = renderToStringWithData(
    <CacheProvider value={cache}>
        <MainPage />
    </CacheProvider>
).then((content) => {
    const { ids, css, html } = extractCritical(content);

    return ReactDOM.renderToString(
        <html>
            <head>
                <meta charset="UTF-8" />
                <meta
                    name="viewport"
                    content="width=device-width,initial-scale=1"
                />
                <style
                    data-emotion-css={ids.join(" ")}
                    dangerouslySetInnerHTML={{ __html: css }}
                />
                <title>Title</title>
                <link rel="stylesheet" href="main.css" />
            </head>
            <body>
                <div id="root" dangerouslySetInnerHTML={{ __html: html }} />
                <script src="https://cdn.jsdelivr.net/npm/react@16.14.0/umd/react.production.min.js" />
                <script src="https://cdn.jsdelivr.net/npm/react-dom@16.14.0/umd/react-dom.production.min.js" />
                <script src={`client.js`} />
            </body>
        </html>
    );
});

// <script src="https://cdn.jsdelivr.net/npm/cross-fetch@3.1.2/dist/browser-ponyfill.min.js" />
// <script src="https://cdn.jsdelivr.net/npm/@apollo/client@3.3.13/apollo-client.cjs.min.js" />

r.then((res) => {
    fs.writeFile("./build/index.htm", `<!DOCTYPE html>${res}`, (err) => {
        if (err) {
            console.error("Failed to write index.htm", err);
            return console.error("============ Failed to write index.htm", err);
        }
        console.log("The file index.htm was saved!");

        const purgeCSSResult = new PurgeCSS().purge({
            content: ["./build/index.htm"],
            css: ["./build/main.css"],
        });

        purgeCSSResult
            .then((res) => {
                // res.map(({ file, css }) =>
                //     fs.writeFile(file, css, (err) => {
                //         if (err) {
                //             return console.error(
                //                 "failed to write CSS file, " + err
                //             );
                //         }
                //         console.log("purgeCSS: The file was saved!");
                //     })
                // );
            })
            .catch((err) => {
                console.log("purgeCSS: failed to purge, ", err);
            });
    });
}).catch((err) => {
    console.error("Failed to render Main page:", err);
    process.exit(1);
});

const signinPage = renderToStringWithData(
    <CacheProvider value={cache}>
        <SignInPage />
    </CacheProvider>
).then((content) => {
    const { ids, css, html } = extractCritical(content);

    return ReactDOM.renderToString(
        <html>
            <head>
                <meta charset="UTF-8" />
                <meta
                    name="viewport"
                    content="width=device-width,initial-scale=1"
                />
                <style
                    data-emotion-css={ids.join(" ")}
                    dangerouslySetInnerHTML={{ __html: css }}
                />
                <title>Sign in</title>
            </head>
            <body>
                <div id="root" dangerouslySetInnerHTML={{ __html: html }} />
            </body>
        </html>
    );
});

signinPage
    .then((res) => {
        fs.writeFile("./build/signin.htm", `<!DOCTYPE html>${res}`, (err) => {
            if (err) {
                return console.error("Failed to write signin.htm", err);
            }
            console.log("The file signin.htm was saved!");
        });
    })
    .catch((err) => {
        console.error("Failed to render Signin page:", err);
        process.exit(1);
    });

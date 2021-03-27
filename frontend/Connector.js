import React from "react";
import fetch from "cross-fetch";
import {
    ApolloProvider,
    ApolloClient,
    InMemoryCache,
    HttpLink,
} from "@apollo/client";

const cache = new InMemoryCache();

const client = new ApolloClient({
    ssrMode: false,
    link: new HttpLink({ uri: "http://127.0.0.1/query", fetch }),
    cache,
    credentials: "same-origin",
});

const Connector = ({ children }) => (
    <ApolloProvider client={client}>{children}</ApolloProvider>
);

export default Connector;

import React from "react";
import { ApolloProvider, ApolloClient, InMemoryCache } from "@apollo/client";

const cache = new InMemoryCache();

const client = new ApolloClient({
    uri: "http://127.0.0.1/query",
    cache,
    credentials: "same-origin",
});

const Connector = ({ children }) => (
    <ApolloProvider client={client}>{children}</ApolloProvider>
);

export default Connector;

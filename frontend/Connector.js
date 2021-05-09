import React from "react";
import fetch from "cross-fetch";
import {
    ApolloProvider,
    ApolloClient,
    ApolloLink,
    InMemoryCache,
    HttpLink,
    split,
} from "@apollo/client";
import { getMainDefinition } from "@apollo/client/utilities";
import { WebSocketLink } from "@apollo/client/link/ws";

const cache = new InMemoryCache();
const graphqlHost =
    process.env.STORYBOOK_GRAPHQL_HOST ||
    process.env.GRAPHQL_HOST ||
    "127.0.0.1";
const graphqlSchema = process.env.GRAPHQL_SCHEMA || "http";
const wsLink =
    process.browser && process.env.NODE_ENV === "production"
        ? new WebSocketLink({
              uri: `ws://${graphqlHost}/query`,
              options: {
                  reconnect: false,
              },
          })
        : null;

const httpLink = new HttpLink({
    uri: `${graphqlSchema}://${graphqlHost}/query`,
    fetch,
});

const splitLink =
    process.browser && process.env.NODE_ENV === "production"
        ? split(
              ({ query }) => {
                  const definition = getMainDefinition(query);
                  return (
                      definition.kind === "OperationDefinition" &&
                      definition.operation === "subscription"
                  );
              },
              wsLink,
              httpLink
          )
        : httpLink;

const client = new ApolloClient({
    ssrMode: false,
    link: ApolloLink.from([splitLink]),
    cache,
    credentials: "same-origin",
});

const Connector = ({ children }) => (
    <ApolloProvider client={client}>{children}</ApolloProvider>
);

export default Connector;

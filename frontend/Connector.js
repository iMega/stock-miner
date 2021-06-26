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

// const wsLink = process.browser
//     ? new WebSocketLink({
//           uri: WS_HOST,
//           options: { reconnect: false },
//       })
//     : null;

const httpLink = new HttpLink({
    uri: process.env.STORYBOOK_GRAPHQL_HOST,
    fetch,
});

const splitLink = process.browser
    ? split(
          ({ query }) => {
              const definition = getMainDefinition(query);
              return (
                  definition.kind === "OperationDefinition" &&
                  definition.operation === "subscription"
              );
          },
          //   wsLink,
          httpLink
      )
    : httpLink;

const client = new ApolloClient({
    ssrMode: false,
    link: ApolloLink.from([httpLink]),
    cache,
    credentials: "same-origin",
});

const Connector = ({ children }) => (
    <ApolloProvider client={client}>{children}</ApolloProvider>
);

export default Connector;

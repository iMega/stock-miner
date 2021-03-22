import React from "react";

export const decorators = [
    (Story) => (
        <React.Fragment>
            <Story />
        </React.Fragment>
    ),
];

export const parameters = { layout: "fullscreen" };

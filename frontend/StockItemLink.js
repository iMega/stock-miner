import React from "react";

const Link = (text) => (
    <a
        target="_blank"
        rel="noopener,nofollow"
        href={"https://finance.yahoo.com/quote/" + text}
    >
        {text}
    </a>
);

export default Link;

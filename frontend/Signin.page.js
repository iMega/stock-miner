import React from "react";
import { Global } from "@emotion/core";
import styled from "@emotion/styled";

import GoogleIcon from "./google.svg";

const Signin = () => (
    <React.Fragment>
        <Global styles={Style} />
        <Main>
            <a href="/google/login">
                <Grid>
                    <IconWrapper>
                        <Icon>
                            <GoogleIcon />
                        </Icon>
                    </IconWrapper>
                    <Text>Sign in with Google</Text>
                </Grid>
            </a>
        </Main>
    </React.Fragment>
);

const Main = styled.main`
    height: 100vh;
    width: 100vm;
    display: grid;
    justify-content: center;
    align-content: center;
`;

const Grid = styled.div`
    background-color: white;
    min-height: 4em;
    min-width: 15em;
    border-radius: 8px;
    box-shadow: 8px 13px 40px 15px rgba(0, 0, 0, 0.73);
    display: grid;
    grid-template-columns: 1fr 3fr;
    transition: all 0.2s;
    &:hover {
        box-shadow: 8px 13px 20px 5px rgba(0, 0, 0, 0.73);
        margin: 2px 0px 0px 2px;
    }
`;

const IconWrapper = styled.div`
    background-color: #4285f4;
    border-radius: 8px 0px 0px 8px;
    display: grid;
    align-content: center;
    justify-content: center;
`;

const Icon = styled.div`
    color: #fff;
    width: 2em;
`;

const Text = styled.div`
    align-content: center;
    display: grid;
    justify-content: center;
    color: #424242;
`;

const Style = `
    html,
    body,
    div {
        margin: 0;
        padding: 0;
        border: 0;
        font-size: 100%;
        font: inherit;
        vertical-align: baseline;
    }
    body {
        font-family: sans-serif;
        color: #fff;
        height: 100vh;
        width: 100vm;
        background-image: radial-gradient(
            ellipse at center,
            Crimson 0%,
            black 100%
        ),
        repeating-linear-gradient(
            45deg,
            transparent,
            transparent 1px,
            rgba(0, 0, 0, 0.2) 1px,
            rgba(255, 255, 255, 0.1) 10px
        ),
        repeating-linear-gradient(
            -45deg,
            transparent,
            transparent 1px,
            rgba(0, 0, 0, 0.2) 1px,
            rgba(255, 255, 255, 0.1) 10px
        );
        background-blend-mode: multiply;
    }
    a {
        text-decoration-line: none;
    }
`;

export default Signin;

const CAN_USE_DOM = Boolean(
    typeof window !== "undefined" &&
        window.document &&
        window.document.createElement
);

export { CAN_USE_DOM };

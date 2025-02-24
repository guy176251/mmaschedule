module.exports = {
    content: ["./components.templ"],
    theme: {
        extend: {},
    },
    plugins: [
        require("@tailwindcss/forms"),
        require("@tailwindcss/typography"),
        require("@tailwindcss/line-clamp"),
        require("@tailwindcss/aspect-ratio"),
        require("daisyui"),
    ],
    daisyui: {
        themes: ["light", "night"],
    },
};

import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./public/index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        bloom: {
          50: "#fdf2f6",
          100: "#fbe4ec",
          200: "#f7c5d6",
          300: "#ef96b4",
          400: "#e36690",
          500: "#d24d7d",
          600: "#c43b6e",
          700: "#a3325a",
          800: "#872c4c",
          900: "#722942",
          DEFAULT: "#C43B6E",
        },
        stem: {
          50: "#f1f8f3",
          100: "#dcefe2",
          200: "#bbdfc6",
          300: "#8cc7a1",
          400: "#5baa78",
          500: "#3b8f5b",
          600: "#2f7d4a",
          700: "#26643c",
          800: "#215033",
          900: "#1c422b",
          DEFAULT: "#2F7D4A",
        },
        pollen: {
          50: "#fef9ec",
          100: "#fbefc9",
          200: "#f6dc8e",
          300: "#f0c454",
          400: "#e8a317",
          500: "#d88b10",
          600: "#ba6a0b",
          700: "#944b0d",
          800: "#7a3c12",
          900: "#663312",
          DEFAULT: "#E8A317",
        },
        ink: {
          50: "#f7f6f5",
          100: "#e8e4e0",
          200: "#d4cdc6",
          300: "#b7ada3",
          400: "#97897d",
          500: "#7c6f64",
          600: "#655a51",
          700: "#534a44",
          800: "#463f3b",
          900: "#1c1917",
          DEFAULT: "#1C1917",
        },
        paper: {
          50: "#fbf7f2",
          100: "#f4ebe1",
          200: "#e8dfd6",
          300: "#d9d0c7",
          400: "#c2b6aa",
          DEFAULT: "#FBF7F2",
        },
      },
      fontFamily: {
        sans: [
          "Inter",
          "system-ui",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "sans-serif",
        ],
        heading: [
          "Fraunces",
          "Georgia",
          "serif",
        ],
      },
    },
  },
  plugins: [],
};

export default config;

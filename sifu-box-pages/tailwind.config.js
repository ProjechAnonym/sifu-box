// tailwind.config.js
import { nextui } from "@nextui-org/react";

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    // ...
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      screens: {
        sm: "340px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
        "2xl": "1536px",
      },
      gridTemplateColumns: {
        "repeat-2": "repeat(2, minmax(0, 1fr))",
        "repeat-1": "repeat(1, minmax(0, 1fr))",
      },
    },
  },
  darkMode: "class",
  plugins: [
    nextui({
      themes: {
        sifulight: {
          extend: "light",
          colors: {
            background: "#F0FCFF",
            foreground: "#18181b",
          },
        },
        sifudark: {
          extend: "dark",
          colors: {
            background: "#18181B",
            foreground: "#e6f1fe",
          },
        },
      },
    }),
  ],
};

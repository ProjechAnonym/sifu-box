import { heroui } from "@heroui/theme";
import { animations } from "framer-motion";
import { transform } from "lodash";

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/layouts/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      keyframes: {
        showIn: {
          "0%": { transform: "translateX(-50%)", opacity: 0 },
          "100%": { transform: "translateX(0)", opacity: 1 },
        },
        showOut: {
          "0%": { transform: "translateX(0)", opacity: 1 },
          "100%": { transform: "translateX(50%)", opacity: 0 },
        },
      },
      animation: {
        showIn_normal: "showIn 0.5s ease",
        showOut_normal: "showOut 0.5s ease",
      },
    },
  },
  darkMode: "class",
  plugins: [
    heroui({
      defaultTheme: "sifulight",
      themes: {
        sifulight: {
          extend: "light",
          layout: {}, // light theme layout tokens
          colors: {
            background: "#fffefb",
            foreground: "#1d1c1c",
            content1: "#f5f4f1",
            content2: "#d4eaf7",
            content3: "#f4f4f5",
          }, // light theme colors
        },
        sifudark: {
          extend: "dark",
          layout: {}, // dark theme layout tokens
          colors: {
            background: "#1E1E1E",
            foreground: "#FFFFFF",
            content1: "#2d2d2d",
            content2: "#2E8B57",
            content3: "#27272A",
          }, // dark theme colors
        },
      },
    }),
  ],
};

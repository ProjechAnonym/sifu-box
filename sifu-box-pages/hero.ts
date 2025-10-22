// hero.ts
import { heroui } from "@heroui/theme";
// or import from theme package if you are using individual packages.
// import { heroui } from "@heroui/theme";
export default heroui(
    {
        addCommonColors: true,
        themes: {
            sifulight: {
            extend: "light", 
            layout: {}, // light theme layout tokens
            colors: {
               default: {
                    50: "#F0FCFF",
                    100: "#FEFEFF",
                    200: "#EDEDEF",
                    300: "#C3F4FD",
                    400: "#338EF7",
                    500: "#7EE7FC",
                    600: "#06B7DB",
                    700: "#09AACD",
                    800: "#0E8AAA",
                    900: "#053B48",
                    DEFAULT:"#E6FAFE"
                },
                background: "#EFFFFF",
                foreground: "#001731",
                content1: "#F9FFFF",
                content2: "#C3F4FD",
                content3: "#3f3f46",
                content4: "#52525b",
            }, // light theme colors
            },
            sifudark: {
            extend: "dark",
            layout: {}, // dark theme layout tokens
            colors: {
                default: {
                    50: "#18181b",
                    100: "#212126",
                    200: "#3f3f46",
                    300: "#52525b",
                    400: "#71717a",
                    500: "#a1a1aa",
                    600: "#d4d4d8",
                    700: "#e4e4e7",
                    800: "#f4f4f5",
                    900: "#fafafa",
                    DEFAULT:"#3f3f46"
                },
                background: "#18181B",
                foreground: "#E6F1FE",
                content1: "#26262a",
                content2: "#27272a",
                content3: "#3f3f46",
                content4: "#52525b",
            }, // dark theme colors
            },
        },
    }
);
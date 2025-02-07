import { ReactNode } from "react";
import { useAppSelector } from "@/redux/hooks";

import Header from "./header";
import Footer from "./footer";
export default function DefaultLayout({ children }: { children: ReactNode }) {
  const theme = useAppSelector((state) => state.theme.theme);
  return (
    <div className={`${theme} bg-background text-foreground w-dvw h-dvh`}>
      <Header />
      <main style={{ height: "calc(100dvh - 3rem - 4rem - 1px)" }}>
        {children}
      </main>
      <Footer theme={theme} />
    </div>
  );
}

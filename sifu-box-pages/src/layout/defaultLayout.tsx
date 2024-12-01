import { ReactNode, useState } from "react";
import { useLocation } from "react-router-dom";
import { useAppSelector } from "@/redux/hooks";
import MenuBar from "./menuBar";
import FooterBar from "./footerBar";
import LoginHead from "./loginHead";
export default function DefaultLayout(props: { children: ReactNode }) {
  const { children } = props;
  const location = useLocation();
  const dark = useAppSelector((state) => state.mode.dark);
  const [menu_height, set_Menu_Height] = useState(0);
  const [footer_height, set_Footer_Height] = useState(0);
  return (
    <div
      className={`${
        dark ? "sifudark" : "sifulight"
      } text-foreground bg-background h-dvh w-full`}
    >
      {location.pathname !== "/login" ? (
        <MenuBar height_callback={(height) => set_Menu_Height(height)} />
      ) : (
        <LoginHead height_callback={(height) => set_Menu_Height(height)} />
      )}
      <div
        style={{
          height: `calc(100dvh - ${menu_height}px - ${footer_height}px - 1px)`,
        }}
      >
        {children}
      </div>
      <FooterBar
        height_callback={(height) => set_Footer_Height(height)}
        dark={dark}
      />
    </div>
  );
}

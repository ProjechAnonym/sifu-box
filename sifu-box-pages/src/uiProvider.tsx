import { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { NextUIProvider } from "@nextui-org/react";
export default function UIProvider({ children }: { children: ReactNode }) {
  const navgate = useNavigate();
  return <NextUIProvider navigate={navgate}>{children}</NextUIProvider>;
}

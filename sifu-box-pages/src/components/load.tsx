import { ReactNode } from "react";
import { Spinner } from "@heroui/spinner";
export default function Load(props: {
  show: boolean;
  children?: ReactNode | undefined;
  fullscreen?: boolean;
  label?: string;
  label_color?:
    | "primary"
    | "secondary"
    | "success"
    | "warning"
    | "danger"
    | undefined;
  color?:
    | "default"
    | "current"
    | "white"
    | "primary"
    | "secondary"
    | "success"
    | "warning"
    | "danger"
    | undefined;
  size?: "sm" | "md" | "lg" | undefined;
}) {
  const {
    show,
    children,
    fullscreen = false,
    size = "md",
    label = "Loading...",
    label_color = "foreground",
    color = "default",
  } = props;

  return (
    <div className="relative">
      {show && (
        <div
          className={`${
            fullscreen ? "fixed" : "absolute"
          } w-full h-full z-50 backdrop-blur-sm flex justify-center items-center translate-x-1/2 -translate-y-1/2 top-1/2 right-1/2`}
        >
          <Spinner
            size={size}
            label={label}
            color={color}
            labelColor={label_color}
          />
        </div>
      )}
      {children !== undefined && children}
    </div>
  );
}

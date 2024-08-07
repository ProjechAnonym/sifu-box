import { useRef, useEffect } from "react";
import { useAppDispatch } from "@/redux/hooks";
import { setDark } from "@/redux/slice";
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  Switch,
} from "@nextui-org/react";

export default function LoginHead(props: {
  height_callback: (height: number) => void;
}) {
  const dispatch = useAppDispatch();
  const ref = useRef<HTMLDivElement>(null);
  const { height_callback } = props;
  useEffect(() => {
    ref.current && height_callback(ref.current.clientHeight);
  }, [ref]);
  return (
    <Navbar
      ref={ref}
      maxWidth="full"
      height={"3rem"}
      isBordered
      classNames={{
        wrapper: ["sm:px-1", "md:px-4", "gap-4"],
      }}
    >
      <NavbarBrand className="font-bold text-lg text-foreground justify-start">
        SifuBox
      </NavbarBrand>
      <NavbarContent className="flex gap-2" justify="center">
        <NavbarItem>
          <Switch
            classNames={{
              wrapper: ["group-data-[selected=true]:bg-orange-400"],
            }}
            size="md"
            onValueChange={() => dispatch(setDark())}
            startContent={
              <span>
                <i className="bi bi-sun-fill text-foreground" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi-moon-fill text-foreground" />
              </span>
            }
          />
        </NavbarItem>
      </NavbarContent>
    </Navbar>
  );
}

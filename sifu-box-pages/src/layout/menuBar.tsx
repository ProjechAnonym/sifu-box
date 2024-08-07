import { useRef, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { setDark, setStatus, setSecret } from "@/redux/slice";
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  NavbarMenuToggle,
  Link,
  Switch,
  Button,
  NavbarMenu,
  NavbarMenuItem,
} from "@nextui-org/react";
import routes from "@/routes";
export default function MenuBar(props: {
  height_callback: (height: number) => void;
}) {
  const dispatch = useAppDispatch();
  const location = useLocation();
  const navigate = useNavigate();
  const dark = useAppSelector((state) => state.mode.dark);
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
        item: [
          "flex",
          "relative",
          "h-full",
          "items-center",
          "data-[active=true]:after:content-['']",
          "data-[active=true]:after:absolute",
          "data-[active=true]:after:bottom-0",
          "data-[active=true]:after:left-0",
          "data-[active=true]:after:right-0",
          "data-[active=true]:after:h-[2px]",
          "data-[active=true]:after:rounded-[2px]",
          "data-[active=true]:after:bg-sky-500",
        ],
      }}
    >
      <NavbarContent className="sm:block md:hidden">
        <NavbarMenuToggle />
      </NavbarContent>
      <NavbarMenu
        className={`bg-gray-600/15 backdrop-blur-md ${
          dark ? "sifudark" : "sifulight"
        }`}
      >
        {routes.map((item, i) => (
          <NavbarMenuItem key={`${item.label}-${i}`}>
            <Link
              href={item.path}
              className={`text-foreground gap-x-1 ${
                location.pathname === item.path && "text-sky-500"
              }`}
            >
              <i className={`${item.icon} text-xl`} />
              <span className="text-xl font-black">{item.label}</span>
            </Link>
          </NavbarMenuItem>
        ))}
      </NavbarMenu>
      <NavbarContent className="sm:hidden md:flex gap-2" justify="center">
        {routes.map((item, i) => (
          <NavbarItem
            key={`${item.label}-${i}`}
            isActive={location.pathname === item.path}
          >
            <Link
              href={item.path}
              className={`text-foreground gap-x-1 ${
                location.pathname === item.path && "text-sky-500"
              }`}
            >
              <i className={`${item.icon} text-xl`} />
              <span className="text-xl font-black">{item.label}</span>
            </Link>
          </NavbarItem>
        ))}
      </NavbarContent>
      <NavbarBrand className="font-bold text-lg text-foreground justify-center">
        SifuBox
      </NavbarBrand>
      <NavbarContent className="flex gap-2" justify="center">
        <NavbarItem>
          <Button
            size="sm"
            radius="sm"
            color="primary"
            onPress={() => {
              navigate("/login");
              dispatch(setStatus(false));
              dispatch(setSecret(""));
            }}
            variant="shadow"
            isIconOnly
          >
            <i className="bi bi-box-arrow-right text-lg" />
          </Button>
        </NavbarItem>
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

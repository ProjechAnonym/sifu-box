import { useLocation } from "react-router-dom";
import { useAppSelector, useAppDispatch } from "@/redux/hooks";
import { useNavigate } from "react-router-dom";
import { setStatus, setAuto, setAdmin } from "@/redux/slice";
import { Navbar, NavbarBrand, NavbarContent, NavbarItem } from "@heroui/navbar";
import { Link } from "@heroui/link";
import { Button } from "@heroui/button";
import { ThemeSwitch } from "@/components/switch";
export default function Header() {
  const status = useAppSelector((state) => state.auth.status);
  const admin = useAppSelector((state) => state.auth.admin);
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const location = useLocation();
  return (
    <Navbar shouldHideOnScroll isBordered maxWidth="full" height={"3rem"}>
      {status ? (
        <NavbarBrand>
          <NavbarItem className="px-2">
            <Link href="/home" isBlock color="foreground">
              <span className="text-xl font-black">首页</span>
            </Link>
          </NavbarItem>
          <NavbarItem className="px-2">
            <Link
              href="/setting"
              isBlock
              color="foreground"
              isDisabled={!admin}
            >
              <span className="text-xl font-black">设置</span>
            </Link>
          </NavbarItem>
        </NavbarBrand>
      ) : (
        <NavbarBrand>
          <span className="text-xl font-black">sifu-box</span>
        </NavbarBrand>
      )}
      <NavbarContent justify="end">
        <NavbarItem className="items-center flex gap-x-2">
          {status && (
            <Button
              color="danger"
              variant="shadow"
              size="sm"
              onPress={() => {
                if (status || location.pathname === "/register") {
                  dispatch(setStatus(false));
                  dispatch(setAuto(false));
                  dispatch(setAdmin(false));
                  navigate("/");
                }
              }}
            >
              <span className="text-xl font-black">退出</span>
            </Button>
          )}
          <ThemeSwitch />
        </NavbarItem>
      </NavbarContent>
    </Navbar>
  );
}

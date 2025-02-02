import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { setTheme } from "@/redux/slice";
import { Switch } from "@heroui/switch";

export function ThemeSwitch() {
  const theme = useAppSelector((state) => state.theme.theme);
  const dispatch = useAppDispatch();
  return (
    <Switch
      endContent={
        <span>
          <i className="bi bi-moon-fill" />
        </span>
      }
      startContent={
        <span>
          <i className="bi bi-sun-fill" />
        </span>
      }
      color="warning"
      isSelected={theme === "sifulight"}
      onValueChange={(isSelect) =>
        isSelect
          ? dispatch(setTheme("sifulight"))
          : dispatch(setTheme("sifudark"))
      }
    />
  );
}

import {
  Button,
  Select,
  Tooltip,
  SelectItem,
  TimeInput,
  TimeInputValue,
} from "@nextui-org/react";
export default function SetInterval(props: {
  dark: boolean;
  interval: TimeInputValue;
  weekday: number;
  time: TimeInputValue | null;
  setTime: (time: TimeInputValue | null) => void;
  setWeekday: (weekday: number) => void;
}) {
  const weekdays = [
    "周日",
    "周一",
    "周二",
    "周三",
    "周四",
    "周五",
    "周六",
    "每天",
    "取消",
  ];
  const { dark, interval, weekday, time, setTime, setWeekday } = props;
  return (
    <>
      <div className="flex flex-row gap-2 justify-center items-center">
        <Select
          selectedKeys={[weekday.toString()]}
          defaultSelectedKeys={["7"]}
          value={weekday.toString()}
          size="sm"
          label="day"
          className="w-20"
          classNames={{
            popoverContent: `${
              dark ? "sifudark" : "sifulight"
            } bg-default-100 text-foreground`,
          }}
          onSelectionChange={(keys) =>
            Object.entries(keys).forEach((key) => setWeekday(parseInt(key[1])))
          }
        >
          {weekdays.map((weekday, i) => (
            <SelectItem key={i}>{weekday}</SelectItem>
          ))}
        </Select>
        <TimeInput
          className="w-20"
          defaultValue={interval}
          value={time}
          hourCycle={24}
          label="interval"
          size="sm"
          errorMessage={(value) => value.isInvalid && "Invalid"}
          onChange={(e) => setTime(e)}
        />
        <Tooltip content="取消自动更新">
          <Button
            size="sm"
            type="button"
            color="danger"
            onPress={() => {
              setTime(null);
              setWeekday(8);
            }}
            startContent={<i className="bi bi-x-lg text-xl" />}
          >
            <span className="text-lg font-black">取消</span>
          </Button>
        </Tooltip>
      </div>
    </>
  );
}

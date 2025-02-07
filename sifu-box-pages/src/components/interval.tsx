import { useState } from "react";
import { Select, SelectItem } from "@heroui/select";
import { Input } from "@heroui/input";
import { Tooltip } from "@heroui/tooltip";
import { Button, ButtonGroup } from "@heroui/button";
import { TimeInput, TimeInputValue } from "@heroui/date-input";
import toast from "react-hot-toast";
import { SetInterval } from "@/utils/application";

export default function Interval(props: { theme: string; token: string }) {
  const { theme, token } = props;
  const [time, setTime] = useState<TimeInputValue | null>(null);
  const [expression, setExpression] = useState<string>("");
  const [cancel, setCancel] = useState(false);
  const [cron, setCron] = useState(false);
  const [day, setDay] = useState<string>("");

  const weekdays = [
    { key: "0", label: "周日" },
    { key: "1", label: "周一" },
    { key: "2", label: "周二" },
    { key: "3", label: "周三" },
    { key: "4", label: "周四" },
    { key: "5", label: "周五" },
    { key: "6", label: "周六" },
    { key: "*", label: "每天" },
  ];
  return (
    <form
      className="w-full flex flex-col gap-2"
      onSubmit={(e) => {
        e.preventDefault();
        if (cancel) {
          toast.promise(SetInterval(token, ""), {
            loading: "取消自动更新中...",
            success: (res) => (res ? "取消自动更新完成" : "取消自动更新失败"),
            error: (e) =>
              e.code === "ERR_NETWORK"
                ? "请检查网络连接"
                : e.response.data.message
                  ? e.response.data.message
                  : e.response.data,
          });
          setCancel(false);
          return;
        }
        if (!time && !cron) {
          toast.error("设置更新时间");
          return;
        }
        const interval = `${time && time.minute} ${time && time.hour} * * ${day}`;
        cron
          ? toast.promise(SetInterval(token, expression), {
              loading: "设置自动更新中...",
              success: (res) => (res ? "设置自动更新完成" : "设置自动更新失败"),
              error: (e) =>
                e.code === "ERR_NETWORK"
                  ? "请检查网络连接"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data,
            })
          : toast.promise(SetInterval(token, interval), {
              loading: "设置自动更新中...",
              success: (res) => (res ? "设置自动更新完成" : "设置自动更新失败"),
              error: (e) =>
                e.code === "ERR_NETWORK"
                  ? "请检查网络连接"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data,
            });
      }}
    >
      <div className="flex flex-row gap-2 items-center">
        {cron ? (
          <Input
            variant="underlined"
            size="sm"
            label="Cron表达式"
            value={expression}
            onValueChange={setExpression}
          />
        ) : (
          <>
            <Select
              size="sm"
              onSelectionChange={(e) => e.currentKey && setDay(e.currentKey)}
              label={"日期"}
              value={day}
              variant="underlined"
              className="w-20"
              classNames={{
                popoverContent: `${theme} bg-content1 text-foreground`,
              }}
              defaultSelectedKeys={["*"]}
            >
              {weekdays.map((value) => (
                <SelectItem key={value.key} textValue={value.label}>
                  <span className="text-xs font-black">{value.label}</span>
                </SelectItem>
              ))}
            </Select>
            <TimeInput
              label="Interval"
              size="sm"
              variant="underlined"
              hourCycle={24}
              value={time}
              onChange={(e) => e && setTime(e)}
              className="w-16"
            />
          </>
        )}
        <ButtonGroup>
          <Button
            size="sm"
            color="primary"
            variant="shadow"
            onPress={() => setCron(!cron)}
            type="button"
          >
            <span className="font-black text-md">
              {cron ? "预设" : "自定义"}
            </span>
          </Button>
        </ButtonGroup>
      </div>
      <div className="flex gap-2 justify-end px-4">
        <ButtonGroup>
          <Tooltip
            content={<span className="fong-black">将取消自动更新配置文件</span>}
            classNames={{ content: `${theme} bg-content1 text-foreground` }}
          >
            <Button
              size="sm"
              color="danger"
              variant="shadow"
              type="submit"
              onPress={() => setCancel(true)}
            >
              <span className="font-black text-lg">手动</span>
            </Button>
          </Tooltip>
        </ButtonGroup>
        <ButtonGroup>
          <Button size="sm" color="primary" variant="shadow" type="submit">
            <span className="font-black text-lg">确认</span>
          </Button>
        </ButtonGroup>
      </div>
    </form>
  );
}

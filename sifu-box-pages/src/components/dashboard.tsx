import { useState, useEffect, useRef } from "react";
import { Button, ButtonGroup } from "@heroui/button";
import { Image } from "@heroui/image";
import { Chip } from "@heroui/chip";
import { Popover, PopoverContent, PopoverTrigger } from "@heroui/popover";
import {
  Autocomplete,
  AutocompleteItem,
  AutocompleteSection,
} from "@heroui/autocomplete";
import toast from "react-hot-toast";
import Interval from "./interval";
import copy from "copy-to-clipboard";
import { Check, Reload, Boot, Stop, Restart, Refresh } from "@/utils/exec";
import { Export, Migrate } from "@/utils/migrate";
import { FetchFiles } from "@/utils/files";
import vpsPic from "@/assets/pictures/vps.svg";
import templatePic from "@/assets/pictures/template.svg";
export function HomeDashBoard(props: {
  provider: string;
  template: string;
  token: string;
  admin: boolean;
  theme: string;
}) {
  const { provider, template, token, admin } = props;
  const [status, setStatus] = useState(false);
  const [check, setCheck] = useState(true);
  useEffect(() => {
    token !== "" &&
      check &&
      Check(token)
        .then((res) => {
          setStatus(res);
          setCheck(false);
        })
        .catch(() => {
          setStatus(false);
          setCheck(false);
        });
  }, [token, check]);

  return (
    <div className="flex flex-wrap gap-2 items-center p-2 h-fit">
      <ButtonGroup>
        <Button
          size="sm"
          variant="shadow"
          color="primary"
          onPress={() =>
            toast.promise(Reload(token), {
              loading: "重载中...",
              success: "重载配置成功",
              error: (e) => {
                setStatus(false);
                return e.code === "ERR_NETWORK"
                  ? "网络连接失败"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data;
              },
            })
          }
        >
          <span className="font-black text-lg">重载</span>
        </Button>
        <Button
          size="sm"
          variant="shadow"
          color="primary"
          onPress={() => setCheck(true)}
        >
          <span className="font-black text-lg">检查</span>
        </Button>
        <Button
          size="sm"
          variant="shadow"
          color="primary"
          isDisabled={!admin}
          onPress={() =>
            toast.promise(Restart(token), {
              loading: "重启中...",
              success: "重启sing-box成功",
              error: (e) => {
                setStatus(false);
                return e.code === "ERR_NETWORK"
                  ? "网络连接失败"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data;
              },
            })
          }
        >
          <span className="font-black text-lg">重启</span>
        </Button>
        <Button
          size="sm"
          variant="shadow"
          color="primary"
          isDisabled={!admin}
          onPress={() =>
            toast.promise(Boot(token), {
              loading: "启动中...",
              success: "启动sing-box成功",
              error: (e) => {
                setStatus(false);
                return e.code === "ERR_NETWORK"
                  ? "网络连接失败"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data;
              },
            })
          }
        >
          <span className="font-black text-lg">启动</span>
        </Button>
        <Button
          size="sm"
          variant="shadow"
          color="primary"
          isDisabled={!admin}
          onPress={() =>
            toast.promise(Stop(token), {
              loading: "关闭中...",
              success: () => {
                setStatus(false);
                return "关闭sing-box成功";
              },
              error: (e) => {
                setCheck(true);
                return e.code === "ERR_NETWORK"
                  ? "网络连接失败"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data;
              },
            })
          }
        >
          <span className="font-black text-lg">关闭</span>
        </Button>
      </ButtonGroup>
      <div className="flex flex-row gap-x-1">
        <Image src={vpsPic} width={40} />
        <div className="flex flex-col justify-center">
          <p className="text-xs font-black select-none">
            {provider !== "" ? provider : "无"}
          </p>
          <p className="text-xs font-black select-none">当前机场</p>
        </div>
      </div>
      <div className="flex flex-row gap-x-1">
        <Image src={templatePic} width={40} />
        <div className="flex flex-col justify-center">
          <p className="text-xs font-black select-none">
            {template !== "" ? template : "无"}
          </p>
          <p className="text-xs font-black select-none">当前模板</p>
        </div>
      </div>
      <Chip size="md" radius="sm" variant="shadow">
        <span className="font-black text-lg select-none">
          sing-box
          {status ? (
            <span className="text-xl text-green-500"> · </span>
          ) : (
            <span className="text-xl text-rose-600"> · </span>
          )}
        </span>
      </Chip>
    </div>
  );
}
export function SettingDashBoard(props: {
  token: string;
  admin: boolean;
  theme: string;
  setUpdate: (update: boolean) => void;
}) {
  const { token, admin, theme, setUpdate } = props;
  const fileInput = useRef<HTMLInputElement>(null);
  const [files, setFiles] = useState<null | {
    [key: string]: Array<{ label: string; path: string }>;
  }>(null);
  useEffect(() => {
    token !== "" &&
      FetchFiles(token)
        .then((res) => res.status && setFiles(res.message))
        .catch((e) =>
          e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? typeof e.response.data.message === "string"
                ? toast.error(e.response.data.message)
                : (e.response.data.message.errors as Array<string>).map((e) =>
                    toast.error(e)
                  ) && setFiles(e.response.data.message.links)
              : toast.error(e.response.data)
        );
  }, [token]);
  return (
    <header className="p-1 h-14 flex flex-row gap-1 items-center">
      <ButtonGroup>
        <Button
          variant="shadow"
          color="primary"
          size="sm"
          onPress={() =>
            toast.promise(Refresh(token), {
              loading: "更新配置文件中...",
              success: (res) => (res ? "更新配置文件成功" : "更新配置文件失败"),
              error: (e) => {
                if (e.code === "ERR_NETWORK") {
                  return "请检查网络连接";
                }
                if (e.response.data.message) {
                  if (typeof e.response.data.message === "string") {
                    return e.response.data.message;
                  }
                  (e.response.data.message as Array<string>).map((m) =>
                    toast.error(m)
                  );
                  return "";
                }
                return e.response.data;
              },
            })
          }
          isDisabled={!admin}
        >
          <span className="font-black text-lg">刷新</span>
        </Button>
        <Popover
          placement="bottom"
          classNames={{ content: `${theme} bg-content1 text-foreground` }}
        >
          <PopoverTrigger>
            <Button
              variant="shadow"
              color="primary"
              size="sm"
              isDisabled={!admin}
            >
              <span className="font-black text-lg">定时</span>
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-64 h-36">
            <Interval theme={theme} token={token} />
          </PopoverContent>
        </Popover>
        <Button
          variant="shadow"
          color="primary"
          size="sm"
          isDisabled={!admin}
          onPress={() =>
            toast.promise(Export(token), {
              loading: "导出配置中...",
              success: (res) => (res ? "配置导出完成" : "配置导出失败"),
              error: (e) =>
                e.code === "ERR_NETWORK"
                  ? "请检查网络连接"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data,
            })
          }
        >
          <span className="font-black text-lg">导出</span>
        </Button>
        <Button
          variant="shadow"
          color="primary"
          size="sm"
          isDisabled={!admin}
          onPress={() => fileInput.current && fileInput.current.click()}
        >
          <span className="font-black text-lg">导入</span>
        </Button>
      </ButtonGroup>
      <input
        className="hidden"
        type="file"
        ref={fileInput}
        onChange={(e) =>
          e.target.files &&
          toast.promise(Migrate(token, e.target.files[0]), {
            loading: "恢复配置中...",
            success: (res) => {
              setUpdate(true);
              return res ? "导入配置成功" : "导入配置失败";
            },
            error: (e) => {
              fileInput.current && (fileInput.current.value = "");
              setUpdate(true);
              return e.code === "ERR_NETWORK"
                ? "请检查网络连接"
                : e.response.data.message
                  ? e.response.data.message
                  : e.response.data;
            },
          })
        }
      />
      {files && (
        <Autocomplete
          variant="underlined"
          label={<span className="text-xs font-black">文件链接</span>}
          size="sm"
          classNames={{
            popoverContent: `${theme} bg-content3 text-foreground`,
          }}
          onSelectionChange={(key) => {
            const url = window.location.origin;
            copy(`${url}/${key}`);
            toast.success("下载链接已复制到剪切板");
          }}
          className="w-36"
        >
          {Object.entries(files!).map(([key, file]) => (
            <AutocompleteSection title={key} key={key} showDivider>
              {file.map((value) => (
                <AutocompleteItem key={value.path}>
                  {value.label}
                </AutocompleteItem>
              ))}
            </AutocompleteSection>
          ))}
        </Autocomplete>
      )}
    </header>
  );
}

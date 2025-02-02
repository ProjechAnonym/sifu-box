import { useState, useEffect } from "react";
import { Button, ButtonGroup } from "@heroui/button";
import { Image } from "@heroui/image";
import { Chip } from "@heroui/chip";
import {
  Autocomplete,
  AutocompleteItem,
  AutocompleteSection,
} from "@heroui/autocomplete";
import toast from "react-hot-toast";
import copy from "copy-to-clipboard";
import { Check, Refresh, Reload, Boot, Stop, Restart } from "@/utils/exec";
import { FetchFiles } from "@/utils/files";

import vpsPic from "@/assets/pictures/vps.svg";
import templatePic from "@/assets/pictures/template.svg";
export default function DashBoard(props: {
  provider: string;
  template: string;
  token: string;
  admin: boolean;
  theme: string;
}) {
  const { provider, template, token, admin, theme } = props;
  const [files, setFiles] = useState<null | {
    [key: string]: Array<{ label: string; path: string }>;
  }>(null);
  const [status, setStatus] = useState(false);
  const [check, setCheck] = useState(true);
  useEffect(() => {
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
      <Chip size="md" color="primary" radius="sm" variant="shadow">
        <span className="font-black text-lg select-none">
          sing-box
          {status ? (
            <span className="text-xl text-green-500"> · </span>
          ) : (
            <span className="text-xl text-rose-600"> · </span>
          )}
        </span>
      </Chip>
      <div className="flex flex-row gap-4 items-center">
        <Button
          variant="shadow"
          color="primary"
          size="sm"
          onPress={() =>
            Refresh(token)
              .then((res) => res && toast.success("更新配置文件成功"))
              .catch((e) =>
                e.code === "ERR_NETWORK"
                  ? "请检查网络连接"
                  : e.response.data.message
                    ? typeof e.response.data.message === "string"
                      ? toast.error(e.response.data.message)
                      : (e.response.data.message as Array<string>).map((m) =>
                          toast.error(m)
                        )
                    : e.response.data
              )
          }
          isDisabled={!admin}
        >
          <span className="font-black text-lg">刷新</span>
        </Button>
        {files && (
          <Autocomplete
            variant="underlined"
            label="配置文件链接"
            placeholder="配置文件链接"
            size="sm"
            classNames={{
              popoverContent: `${theme} bg-content3 text-foreground`,
            }}
            onSelectionChange={(key) => {
              const url = window.location.origin;
              copy(`${url}/${key}`);
              toast.success("下载链接已复制到剪切板");
            }}
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
      </div>
    </div>
  );
}

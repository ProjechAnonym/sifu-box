import { useState } from "react";
import { Card, CardBody, Button, Tooltip } from "@nextui-org/react";
import toast from "react-hot-toast";
import {
  GetServiceStatus,
  BootServiceUp,
  StopService,
} from "@/utils/service/ServiceExec";
import { HostValue } from "@/types/host";
export default function ServiceMonitor(props: {
  server: HostValue;
  dark: boolean;
  secret: string;
}) {
  const { server, dark, secret } = props;
  const [singbox, setSingbox] = useState(false);
  GetServiceStatus(secret, "sing-box", server.url)
    .then((res) => setSingbox(res))
    .catch(() => setSingbox(false));
  return (
    <Card
      shadow="none"
      className="h-14 w-56 flex justify-center items-center"
      radius="sm"
    >
      <CardBody
        className={`${
          dark ? "bg-zinc-800" : "bg-slate-100"
        } flex flex-row gap-x-2 items-center`}
      >
        <div className="font-black text-md space-x-2 flex justify-center">
          <span className="flex items-center">singbox</span>
          {singbox ? (
            <span className="text-green-600 text-2xl">·</span>
          ) : (
            <span className="text-red-600 text-2xl">·</span>
          )}
        </div>
        <div className="flex flex-row gap-x-1">
          <Tooltip
            content={
              <span
                className={`${dark ? "sifudark" : "sifulight"} text-foreground`}
              >
                boot up
              </span>
            }
            classNames={{
              content: [`${dark ? "bg-zinc-800" : "bg-slate-100"}`],
            }}
          >
            <Button
              isIconOnly
              size="sm"
              onPress={() =>
                toast.promise(BootServiceUp(secret, "sing-box", server.url), {
                  loading: "loading",
                  success: (res) => {
                    setSingbox(res);
                    return `启动sing-box完成`;
                  },
                  error: (err) => {
                    setSingbox(false);
                    return `${err.response.data.message}`;
                  },
                })
              }
            >
              <i className="bi bi-capslock-fill" />
            </Button>
          </Tooltip>
          <Tooltip
            content={
              <span
                className={`${dark ? "sifudark" : "sifulight"} text-foreground`}
              >
                check status
              </span>
            }
            classNames={{
              content: [`${dark ? "bg-zinc-800" : "bg-slate-100"}`],
            }}
          >
            <Button
              isIconOnly
              size="sm"
              onPress={() =>
                toast.promise(
                  GetServiceStatus(secret, "sing-box", server.url),
                  {
                    loading: "loading",
                    success: (res) => {
                      setSingbox(res);
                      return `sing-box is ${res ? "active" : "dead"} now`;
                    },
                    error: (err) => {
                      setSingbox(false);
                      return `${err.response.data.message}`;
                    },
                  }
                )
              }
            >
              <i className="bi bi-radar" />
            </Button>
          </Tooltip>
          <Tooltip
            content={
              <span
                className={`${dark ? "sifudark" : "sifulight"} text-foreground`}
              >
                stop sing-box
              </span>
            }
            classNames={{
              content: [`${dark ? "bg-zinc-800" : "bg-slate-100"}`],
            }}
          >
            <Button
              isIconOnly
              size="sm"
              onPress={() =>
                toast.promise(StopService(secret, "sing-box", server.url), {
                  loading: "loading",
                  success: () => {
                    setSingbox(false);
                    return `关闭sing-box完成`;
                  },
                  error: (err) => {
                    setSingbox(false);
                    return `${err.response.data.message}`;
                  },
                })
              }
            >
              <i className="bi bi-x-lg" />
            </Button>
          </Tooltip>
        </div>
      </CardBody>
    </Card>
  );
}

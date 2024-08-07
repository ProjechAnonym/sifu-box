import { useEffect, useState, useRef } from "react";
import { NavigateFunction, useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import {
  Divider,
  Skeleton,
  Autocomplete,
  AutocompleteSection,
  AutocompleteItem,
  Button,
  Modal,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalBody,
  useDisclosure,
  TimeInputValue,
} from "@nextui-org/react";
import toast from "react-hot-toast";
import Load from "@/components/load";
import ServiceMonitor from "@/components/home/serviceMonitor";
import SwitchConfig from "@/components/home/swichConfig";
import SetInterval from "@/components/home/setInterval";
import ProxyServers from "@/layout/proxyServers";
import { cloneDeep } from "lodash";
import copy from "copy-to-clipboard";
import { ClienAuth } from "@/utils/ClientAuth";
import { FetchHosts } from "@/utils/host/FetchHost";
import { DeleteHost } from "@/utils/host/DelHost";
import { FetchConfigs } from "@/utils/configs/FetchConfig";
import { SelectTemplate } from "@/utils/configs/SwitchConfig";
import { SetIntervalTime } from "@/utils/SetIntervalTime";
import { HostValue } from "@/types/host";
function redirectLogin(navigate: NavigateFunction) {
  navigate("/login");
  toast.error("Please login", { duration: 2000 });
}
function setCurrentHost(key: string, hosts: Array<HostValue>) {
  const keySlice = key.split("-");
  const index = parseInt(keySlice[keySlice.length - 1]);
  return index ? hosts[index - 1] : null;
}

export default function Home() {
  const interval = { hour: 4, minute: 30 } as TimeInputValue;
  const headerRef = useRef<HTMLHeadElement>(null);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const dark = useAppSelector((state) => state.mode.dark);
  const status = useAppSelector((state) => state.auth.status);
  const login = useAppSelector((state) => state.auth.login);
  const load = useAppSelector((state) => state.auth.load);
  const secret = useAppSelector((state) => state.auth.secret);
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const [updateHosts, setUpdateHosts] = useState(true);
  const [hosts, setHosts] = useState<Array<HostValue> | null>(null);
  const [selectHost, setSelectHost] = useState<HostValue | null>(null);
  const [configs, setConfigs] = useState<Array<{
    template: string;
    key: string;
    files: Array<{ path: string; label: string }>;
  }> | null>(null);
  const [time, setTime] = useState<TimeInputValue | null>(interval);
  const [weekday, setWeekday] = useState(7);
  const [headerHeight, setHeaderHeight] = useState(0);
  const [updateServers, setUpdateServers] = useState(true);
  useEffect(() => {
    !login && !status && dispatch(ClienAuth({ auto: true }));
    login && !status && !load && redirectLogin(navigate);
    if (secret !== "" && updateHosts) {
      setUpdateHosts(false);
      FetchHosts(secret)
        .then((res) =>
          res ? setHosts(res) : toast.error("Fetch servers failed")
        )
        .catch(() => toast.error("Fetch servers failed"));
    }
    secret !== "" && FetchConfigs(secret).then((res) => setConfigs(res));
    headerRef.current && setHeaderHeight(headerRef.current.clientHeight);
  }, [
    status,
    login,
    load,
    secret,
    updateHosts,
    headerRef.current?.clientHeight,
  ]);
  return (
    <main className="h-full">
      <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
        <ModalContent className="sifudark bg-background text-foreground">
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">
                设置自动更新时间
              </ModalHeader>
              <form
                className="flex flex-wrap gap-x-1"
                onSubmit={(e) => {
                  e.preventDefault();
                  toast.promise(SetIntervalTime(secret, time, weekday), {
                    loading: "设置中",
                    success: () => {
                      onClose();
                      return weekday === 8
                        ? "取消自动更新成功"
                        : "设置新的自动更新时间成功";
                    },
                    error: (err) => `${err.response.data.message}`,
                  });
                }}
              >
                <ModalBody>
                  <SetInterval
                    dark={dark}
                    time={time}
                    weekday={weekday}
                    setTime={(e) => setTime(e)}
                    setWeekday={(e) => setWeekday(e)}
                    interval={interval}
                  />
                </ModalBody>
                <ModalFooter className="w-full">
                  <Button color="danger" onPress={onClose} type="button">
                    <span className="text-lg font-black">关闭</span>
                  </Button>
                  <Button color="primary" type="submit">
                    <span className="text-lg font-black">提交</span>
                  </Button>
                </ModalFooter>
              </form>
            </>
          )}
        </ModalContent>
      </Modal>
      <Load show={load} fullscreen={true} />
      <header
        className="w-full flex flex-wrap gap-x-4 gap-y-2 items-center px-3 py-1"
        ref={headerRef}
      >
        <div className="w-fit flex flex-row gap-x-2 items-center justify-center">
          {hosts ? (
            <Autocomplete
              selectedKey={selectHost ? selectHost.key : null}
              label="Select a server"
              size="sm"
              radius="sm"
              className="w-60 h-12"
              classNames={{
                popoverContent: `${
                  dark ? "sifudark" : "sifulight"
                } bg-default-100 text-foreground`,
              }}
              onSelectionChange={(e) =>
                e && setSelectHost(setCurrentHost(e.toString(), hosts))
              }
            >
              {hosts.map((host) => (
                <AutocompleteItem
                  key={host.key}
                  startContent={
                    host.localhost && (
                      <i className="bi bi-pc-display-horizontal" />
                    )
                  }
                >
                  {host.url}
                </AutocompleteItem>
              ))}
            </Autocomplete>
          ) : (
            <Skeleton className="w-60 h-12 bg-zinc-400 rounded-lg" />
          )}
          <Button
            className="h-10"
            color="danger"
            size="sm"
            onPress={() =>
              selectHost &&
              toast.promise(
                DeleteHost(secret, selectHost.url).then((res) => {
                  setUpdateHosts(true);
                  res && setSelectHost(null);
                }),
                {
                  loading: "loading",
                  success: "删除成功",
                  error: (err) => `${err.response.data.message}`,
                }
              )
            }
          >
            <i className="bi bi-trash text-lg" />
            <span className="text-lg font-black">删除</span>
          </Button>
        </div>
        <div>
          {configs ? (
            <Autocomplete
              onSelectionChange={(e) => {
                e && copy(e.toString());
                e && toast.success("复制到剪贴板");
              }}
              label="files"
              size="sm"
              radius="sm"
              className="w-36 h-12"
              classNames={{
                popoverContent: `${
                  dark ? "sifudark" : "sifulight"
                } bg-default-100 text-foreground`,
              }}
            >
              {configs.map((template) => (
                <AutocompleteSection
                  showDivider
                  title={template.template}
                  key={template.key}
                >
                  {template.files.map((file) => (
                    <AutocompleteItem key={file.path}>
                      {file.label}
                    </AutocompleteItem>
                  ))}
                </AutocompleteSection>
              ))}
            </Autocomplete>
          ) : (
            <Skeleton className="bg-zinc-400 w-36 h-12 rounded-lg" />
          )}
        </div>
        {configs && selectHost ? (
          <SwitchConfig
            key={`config-${selectHost.key}`}
            secret={secret}
            dark={dark}
            config={selectHost.config}
            url={selectHost.url}
            files={SelectTemplate("default", configs)}
            onSwitch={(config) => {
              selectHost.config = config;
              const newCurrentHost = cloneDeep(selectHost);
              setSelectHost(newCurrentHost);
              setUpdateServers(true);
            }}
          />
        ) : (
          <div
            className={`${
              dark ? "bg-zinc-800" : "bg-slate-100"
            } w-44 h-12 text-center rounded-lg font-black p-3`}
          >
            未选择服务器
          </div>
        )}
        <Button
          color="primary"
          startContent={<i className="bi bi-clock text-lg" />}
          size="sm"
          onPress={onOpen}
          className="p-2 h-10"
        >
          <span className="font-black text-lg">定时器</span>
        </Button>
        {selectHost ? (
          <ServiceMonitor
            key={selectHost.key}
            server={selectHost}
            secret={secret}
            dark={dark}
          />
        ) : (
          <div
            className={`${
              dark ? "bg-zinc-800" : "bg-slate-100"
            } w-56 h-14 text-center rounded-lg font-black text-xl py-3`}
          >
            尚未选择服务器
          </div>
        )}
      </header>
      <Divider />
      {selectHost ? (
        <ProxyServers
          updateServers={updateServers}
          setUpdateServers={(e) => setUpdateServers(e)}
          host={selectHost}
          key={selectHost.key}
          dark={dark}
          headerHeight={headerHeight}
        />
      ) : (
        <></>
      )}
    </main>
  );
}

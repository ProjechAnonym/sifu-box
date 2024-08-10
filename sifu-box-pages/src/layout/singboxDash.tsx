import { useState, useEffect, useMemo } from "react";
import {
  Accordion,
  AccordionItem,
  ScrollShadow,
  Button,
  Spinner,
  Badge,
} from "@nextui-org/react";
import toast from "react-hot-toast";
import SingboxData from "@/components/home/singboxData";
import { cloneDeep } from "lodash";
import { FetchProxyServers } from "@/utils/singbox/FetchProxy";
import { TestGroup, TestHost } from "@/utils/singbox/TestDelay";
import { SwitchServer } from "@/utils/singbox/SwitchServer";
import { SocketUrl } from "@/utils/singbox/SocketConn";
import { FetchRules } from "@/utils/singbox/FetchRules";
import { ParseConnection, ParseLog } from "@/utils/singbox/ParseData";
import { ClearFakeip } from "@/utils/singbox/ClearFakeip";
import { HostValue } from "@/types/host";
import {
  ServerValue,
  ServerGroupValue,
  ConnectionData,
  ConnectionLog,
} from "@/types/singbox";
export default function SingboxDash(props: {
  host: HostValue;
  dark: boolean;
  headerHeight: number;
  updateServers: boolean;
  setUpdateServers: (updateServers: boolean) => void;
}) {
  const { host, dark, headerHeight, updateServers, setUpdateServers } = props;
  const memoryUrl = SocketUrl(host.url, host.port, "memory", [
    { token: host.secret },
  ]);
  const networkUrl = SocketUrl(host.url, host.port, "traffic", [
    { token: host.secret },
  ]);
  const logUrl = SocketUrl(host.url, host.port, "logs", [
    { token: host.secret },
    { level: "info" },
  ]);
  const connectUrl = SocketUrl(host.url, host.port, "connections", [
    { token: host.secret },
  ]);
  const [servers, setServers] = useState<{
    [key: string]: ServerGroupValue | ServerValue;
  } | null>(null);
  const [groups, setGroups] = useState<Array<string> | null>(null);
  const [testGroups, setTestGroups] = useState<Array<string>>([]);
  const [testServers, setTestServers] = useState<Array<string>>([]);
  const [network, setNetwork] = useState({ up: 0, down: 0 });
  const [memory, setMemory] = useState({ inuse: 0, oslimit: 0 });
  const [connections, setConnections] = useState<{
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<ConnectionLog>;
  }>({
    labels: [
      { key: "id", label: "ID", allowSort: true, initShow: false },
      { key: "rule", label: "命中规则", allowSort: true, initShow: true },
      {
        key: "rulePayload",
        label: "规则信息",
        allowSort: true,
        initShow: false,
      },
      { key: "start", label: "时间", allowSort: true, initShow: true },
      { key: "upload", label: "上传", allowSort: true, initShow: true },
      { key: "download", label: "下载", allowSort: true, initShow: true },
      { key: "chains", label: "节点链", allowSort: true, initShow: true },
      {
        key: "destinationIP",
        label: "目标IP",
        allowSort: true,
        initShow: true,
      },
      {
        key: "destinationPort",
        label: "目标端口",
        allowSort: true,
        initShow: false,
      },
      {
        key: "dnsMode",
        label: "dns模式",
        allowSort: false,
        initShow: false,
      },
      {
        key: "network",
        label: "网络",
        allowSort: false,
        initShow: false,
      },
      {
        key: "processPath",
        label: "进程",
        allowSort: true,
        initShow: false,
      },
      {
        key: "sourceIP",
        label: "源IP",
        allowSort: true,
        initShow: false,
      },
      {
        key: "sourcePort",
        label: "源端口",
        allowSort: true,
        initShow: false,
      },
      {
        key: "type",
        label: "类型",
        allowSort: false,
        initShow: false,
      },
      { key: "host", label: "嗅探域名", allowSort: true, initShow: true },
    ],
    values: [],
  });
  const [disConnections, setDisConnections] = useState<{
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<ConnectionLog>;
  }>({
    labels: [
      { key: "id", label: "ID", allowSort: true, initShow: false },
      { key: "rule", label: "命中规则", allowSort: true, initShow: true },
      {
        key: "rulePayload",
        label: "规则信息",
        allowSort: true,
        initShow: false,
      },
      { key: "start", label: "时间", allowSort: true, initShow: true },
      { key: "upload", label: "上传", allowSort: true, initShow: true },
      { key: "download", label: "下载", allowSort: true, initShow: true },
      { key: "chains", label: "节点链", allowSort: true, initShow: true },
      {
        key: "destinationIP",
        label: "目标IP",
        allowSort: true,
        initShow: true,
      },
      {
        key: "destinationPort",
        label: "目标端口",
        allowSort: true,
        initShow: false,
      },
      {
        key: "dnsMode",
        label: "dns模式",
        allowSort: false,
        initShow: false,
      },
      {
        key: "network",
        label: "网络",
        allowSort: false,
        initShow: false,
      },
      {
        key: "processPath",
        label: "进程",
        allowSort: true,
        initShow: false,
      },
      {
        key: "sourceIP",
        label: "源IP",
        allowSort: true,
        initShow: false,
      },
      {
        key: "sourcePort",
        label: "源端口",
        allowSort: true,
        initShow: false,
      },
      {
        key: "type",
        label: "类型",
        allowSort: false,
        initShow: false,
      },
      { key: "host", label: "目标域名", allowSort: true, initShow: true },
    ],
    values: [],
  });
  const [statistics, setStatistics] = useState<{
    upload: number;
    download: number;
  }>({ upload: 0, download: 0 });
  const [logs, setLogs] = useState<{
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{
      type: string;
      payload: string;
      time: string;
      key: string;
    }>;
  }>({
    labels: [
      { label: "时间", key: "time", allowSort: true, initShow: true },
      { label: "等级", key: "type", allowSort: true, initShow: true },
      { label: "信息", key: "payload", allowSort: true, initShow: true },
    ],
    values: [],
  });
  const [rules, setRules] = useState<{
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{
      type: string;
      payload: string;
      proxy: string;
      key: string;
    }>;
  }>({
    labels: [
      { label: "类型", key: "type", allowSort: true, initShow: true },
      { label: "规则", key: "payload", allowSort: true, initShow: true },
      { label: "出站", key: "proxy", allowSort: true, initShow: true },
    ],
    values: [],
  });
  useEffect(() => {
    updateServers &&
      FetchProxyServers(host.url + ":" + host.port, host.secret)
        .then((res) => {
          setUpdateServers(false);
          res ? setServers(res.proxies) : toast.error("获取代理失败");
        })
        .catch((e) => {
          setUpdateServers(false);
          toast.error(e.message);
        });
  }, [updateServers]);
  useEffect(() => {
    FetchRules(host.url + ":" + host.port, host.secret)
      .then((res) => setRules(res))
      .catch(() => toast.error("获取规则列表失败"));
    const networkSocket = new WebSocket(networkUrl);
    networkSocket.addEventListener("open", () => {
      console.log("连接已建立");
      // 这里可以发送消息给服务器
    });
    networkSocket.addEventListener("error", (event) => {
      toast.error("网络连接失败");
      console.error("发生错误:", event);
    });
    networkSocket.addEventListener("message", (event) => {
      setNetwork(JSON.parse(event.data));
    });
    const memorySocket = new WebSocket(memoryUrl);
    memorySocket.addEventListener("open", () => {
      console.log("连接已建立");
      // 这里可以发送消息给服务器
    });
    memorySocket.addEventListener("error", (event) => {
      toast.error("网络连接失败");
      console.error("发生错误:", event);
    });
    memorySocket.addEventListener("message", (event) => {
      setMemory(JSON.parse(event.data));
    });
    const connectionSocket = new WebSocket(connectUrl);
    connectionSocket.addEventListener("open", () => {
      console.log("连接已建立");
      // 这里可以发送消息给服务器
    });
    connectionSocket.addEventListener("error", (event) => {
      toast.error("网络连接失败");
      console.error("发生错误:", event);
    });
    connectionSocket.addEventListener("message", (event) => {
      const connectionData = JSON.parse(event.data) as ConnectionData;
      const { labels, aliveConnections, deadConnections } = ParseConnection(
        connectionData.connections,
        connections.values
      );
      setDisConnections({ labels: labels, values: deadConnections });
      setConnections({ labels: labels, values: aliveConnections });
      setStatistics({
        upload: connectionData.uploadTotal,
        download: connectionData.downloadTotal,
      });
    });
    const logSocket = new WebSocket(logUrl);
    logSocket.addEventListener("open", () => {
      console.log("连接已建立");
      // 这里可以发送消息给服务器
    });
    logSocket.addEventListener("error", (event) => {
      toast.error("网络连接失败");
      console.error("发生错误:", event);
    });
    logSocket.addEventListener("message", (event) =>
      setLogs(ParseLog(JSON.parse(event.data), logs.values))
    );
    return () => {
      networkSocket.close();
      memorySocket.close();
      connectionSocket.close();
      logSocket.close();
    };
  }, []);
  useMemo(() => {
    servers &&
      setGroups(
        Object.values(servers)
          .filter(
            (server) =>
              server.type === "Selector" ||
              server.type === "Fallback" ||
              server.type === "URLTest"
          )
          .map((server) => server.name)
          .reverse()
      );
  }, [servers]);
  return (
    <ScrollShadow style={{ height: `calc(100% - ${headerHeight}px)` }}>
      <SingboxData
        disconnections={disConnections}
        connections={connections}
        rules={rules}
        logs={logs}
        dark={dark}
        network={network}
        memory={memory}
        networkTotal={statistics}
        clearFakeip={() =>
          toast.promise(ClearFakeip(host.url, host.secret), {
            loading: "loading",
            success: "清空fakeip成功",
            error: "清空fakeip失败",
          })
        }
      />
      {groups && (
        <ScrollShadow className="h-fit">
          <Accordion selectionMode="multiple">
            {groups.map((group, i) => (
              <AccordionItem
                aria-label={group}
                title={<span className="text-xl font-black">{group}</span>}
                key={`${group}-${i}`}
              >
                {servers![group].type !== "Fallback" && (
                  <div className="flex flex-row gap-2">
                    <span className="text-xl font-black">测试延迟</span>
                    <Button
                      isLoading={testGroups?.includes(group)}
                      onPress={() => {
                        const newTestGroup = cloneDeep(testGroups);
                        newTestGroup.push(group);
                        setTestGroups(newTestGroup);
                        toast.promise(
                          TestGroup(
                            group,
                            host.url + ":" + host.port,
                            host.secret,
                            servers!
                          ),
                          {
                            loading: "loading",
                            success: (res) => {
                              const newTestGroups = cloneDeep(testGroups);
                              newTestGroups
                                .filter((item) => item !== group)
                                .map((item) => item);
                              setTestGroups(newTestGroups);
                              setServers(res);
                              return `测试${group}完成`;
                            },
                            error: () => {
                              const newTestGroups = cloneDeep(testGroups);
                              newTestGroups
                                .filter((item) => item !== group)
                                .map((item) => item);
                              setTestGroups(newTestGroups);
                              return `测试${group}失败`;
                            },
                          }
                        );
                      }}
                      isIconOnly
                      size="sm"
                      variant="light"
                    >
                      <i className="bi bi-stopwatch-fill text-lg" />
                    </Button>
                  </div>
                )}
                <ScrollShadow className="flex flex-wrap gap-x-3 gap-y-2 max-h-56 py-2">
                  {(servers![group] as ServerGroupValue).all.map(
                    (server: string, j: number) => {
                      return (
                        <Badge
                          showOutline={false}
                          placement="top-right"
                          content={
                            servers![server].udp ? (
                              <span className="text-teal-50 text-sm font-black">
                                udp
                              </span>
                            ) : (
                              <del className="text-teal-50 text-sm font-black">
                                udp
                              </del>
                            )
                          }
                          color={servers![server].udp ? "primary" : "danger"}
                          key={`${group}-${host}-${j}`}
                        >
                          <div
                            className={`${
                              dark
                                ? (servers![group] as ServerGroupValue).now ===
                                  server
                                  ? "bg-primary"
                                  : "bg-zinc-700"
                                : "bg-slate-200"
                            } p-2 sm:w-40 md:w-44 cursor-pointer rounded-md space-x-2`}
                            onClick={(e) => {
                              e.stopPropagation();
                              SwitchServer(
                                group,
                                server,
                                host.url + ":" + host.port,
                                host.secret,
                                servers!
                              )
                                .then((res) => setServers(res))
                                .catch((e) => toast.error(e.message));
                            }}
                          >
                            <span className="text-md font-black">{server}</span>
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                const newTestServers = cloneDeep(testServers);
                                newTestServers.push(server);
                                setTestServers(newTestServers);
                                TestHost(
                                  server,
                                  host.url + ":" + host.port,
                                  host.secret,
                                  servers!
                                )
                                  .then((res) => {
                                    const newTestServers =
                                      cloneDeep(testServers);
                                    newTestServers
                                      .filter((item) => item !== server)
                                      .map((item) => item);
                                    setTestServers(newTestServers);
                                    setServers(res);
                                  })
                                  .catch((e) => {
                                    const newTestServers =
                                      cloneDeep(testServers);
                                    newTestServers
                                      .filter((item) => item !== server)
                                      .map((item) => item);
                                    setTestServers(newTestServers);
                                    toast.error(e.message);
                                  });
                              }}
                              className="font-black text-lime-500 rounded-sm hover:bg-teal-50/30 w-12 h-6 text-center"
                            >
                              <span className="text-xs">
                                {testServers.includes(server) ? (
                                  <Spinner size="sm"></Spinner>
                                ) : servers![server]["history"][0] ? (
                                  servers![server]["history"][0].delay + "ms"
                                ) : (
                                  "--- ms"
                                )}
                              </span>
                            </button>
                          </div>
                        </Badge>
                      );
                    }
                  )}
                </ScrollShadow>
              </AccordionItem>
            ))}
          </Accordion>
        </ScrollShadow>
      )}
    </ScrollShadow>
  );
}

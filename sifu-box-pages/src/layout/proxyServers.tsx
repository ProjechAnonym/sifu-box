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
import { cloneDeep } from "lodash";
import { FetchProxyServers } from "@/utils/servers/FetchProxy";
import { TestGroup, TestHost } from "@/utils/servers/TestDelay";
import { SwitchServer } from "@/utils/servers/SwitchServer";
import { HostValue } from "@/types/host";
import { ServerValue, ServerGroupValue } from "@/types/servers";
export default function ProxyServers(props: {
  host: HostValue;
  dark: boolean;
  headerHeight: number;
  updateServers: boolean;
  setUpdateServers: (updateServers: boolean) => void;
}) {
  const { host, dark, headerHeight, updateServers, setUpdateServers } = props;
  const [servers, setServers] = useState<{
    [key: string]: ServerGroupValue | ServerValue;
  } | null>(null);
  const [groups, setGroups] = useState<Array<string> | null>(null);
  const [testGroups, setTestGroups] = useState<Array<string>>([]);
  const [testServers, setTestServers] = useState<Array<string>>([]);
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
    <div style={{ height: `calc(100% - ${headerHeight}px)` }}>
      {groups && (
        <ScrollShadow className="h-full">
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
    </div>
  );
}

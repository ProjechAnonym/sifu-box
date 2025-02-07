import { useState, useEffect, useMemo } from "react";
import { Accordion, AccordionItem } from "@heroui/accordion";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Button } from "@heroui/button";
import { Badge } from "@heroui/badge";
import { Spinner } from "@heroui/spinner";
import toast from "react-hot-toast";
import NavBlock from "./navBlock";

import { cloneDeep } from "lodash";
import { SwitchOutbound } from "@/utils/singbox/switch";
import { FetchOutbounds } from "@/utils/singbox/fetch";
import { GroupDelay, OutboundDelay } from "@/utils/singbox/delay";
import { Outbound, OutboundGroup } from "@/types/singbox/outbound";
export default function SingBox(props: {
  listen: string;
  secret: string;
  height: number;
  theme: string;
}) {
  const { listen, secret, height, theme } = props;
  const [updateOutbounds, setUpdateOutbounds] = useState(true);
  const [outbounds, setOutbounds] = useState<{
    [key: string]: OutboundGroup | Outbound;
  } | null>(null);
  const [groups, setGroups] = useState<Array<string> | null>(null);
  const [testGroups, setTestGroups] = useState<Array<string>>([]);
  const [testServers, setTestServers] = useState<Array<string>>([]);
  useEffect(() => {
    updateOutbounds &&
      listen !== "" &&
      secret !== "" &&
      FetchOutbounds(listen, secret)
        .then((res) => {
          setUpdateOutbounds(false);
          res ? setOutbounds(res.proxies) : toast.error("获取代理失败");
        })
        .catch((e) => {
          setUpdateOutbounds(false);
          toast.error(e.code === "ERR_NETWORK" ? "请检查网络连接" : e.message);
        });
  }, [listen, secret, updateOutbounds]);
  useMemo(() => {
    outbounds &&
      setGroups(
        Object.values(outbounds)
          .filter(
            (outbound) =>
              outbound.type === "Selector" ||
              outbound.type === "Fallback" ||
              outbound.type === "URLTest"
          )
          .map((outbound) => outbound.name)
          .reverse()
      );
  }, [outbounds]);
  return (
    <div className="p-2" style={{ height: `calc(100% - ${height}px)` }}>
      {groups && <NavBlock groups={groups} />}
      {groups && (
        <ScrollShadow className="w-full h-full">
          <Accordion selectionMode="multiple">
            {groups.map((group, i) => (
              <AccordionItem
                aria-label={group}
                title={<span className="text-xl font-black">{group}</span>}
                key={`${group}-${i}`}
                id={group}
              >
                {outbounds![group].type !== "Fallback" && (
                  <div className="flex flex-row gap-2">
                    <span className="text-xl font-black">测试延迟</span>
                    <Button
                      isLoading={testGroups?.includes(group)}
                      onPress={() => {
                        const newTestGroup = cloneDeep(testGroups);
                        newTestGroup.push(group);
                        setTestGroups(newTestGroup);
                        toast.promise(
                          GroupDelay(group, listen, secret, outbounds!),
                          {
                            loading: "loading",
                            success: (res) => {
                              const newTestGroups = cloneDeep(testGroups);
                              newTestGroups
                                .filter((item) => item !== group)
                                .map((item) => item);
                              setTestGroups(newTestGroups);
                              setOutbounds(res);
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
                  {(outbounds![group] as OutboundGroup).all.map(
                    (outbound: string, j: number) => {
                      return (
                        <Badge
                          showOutline={false}
                          placement="top-right"
                          content={
                            outbounds![outbound].udp ? (
                              <span className="text-teal-50 text-sm font-black">
                                UDP
                              </span>
                            ) : (
                              <del className="text-teal-50 text-sm font-black">
                                UDP
                              </del>
                            )
                          }
                          color={
                            outbounds![outbound].udp ? "primary" : "danger"
                          }
                          key={`${group}-${outbound}-${j}`}
                        >
                          <div
                            className={`${theme} bg-content1 ${
                              (outbounds![group] as OutboundGroup).now ===
                                outbound && "bg-content2"
                            }  p-2 w-40 cursor-pointer rounded-xl flex flex-col`}
                            onClick={(e) => {
                              e.stopPropagation();
                              SwitchOutbound(
                                group,
                                outbound,
                                listen,
                                secret,
                                outbounds!
                              )
                                .then((res) => setOutbounds(res))
                                .catch((e) =>
                                  toast.error(
                                    e.code === "ERR_NETWORK"
                                      ? "请检查网络连接"
                                      : e.response.data.message
                                        ? e.response.data.message
                                        : e.response.data
                                  )
                                );
                            }}
                          >
                            <span className="text-md font-black">
                              {outbound}
                            </span>
                            <span
                              onClick={(e) => {
                                e.stopPropagation();
                                const newTestServers = cloneDeep(testServers);
                                newTestServers.push(outbound);
                                setTestServers(newTestServers);
                                OutboundDelay(
                                  outbound,
                                  listen,
                                  secret,
                                  outbounds!
                                )
                                  .then((res) => {
                                    const newTestServers =
                                      cloneDeep(testServers);
                                    newTestServers
                                      .filter((item) => item !== outbound)
                                      .map((item) => item);
                                    setTestServers(newTestServers);
                                    setOutbounds(res);
                                  })
                                  .catch((e) => {
                                    const newTestServers =
                                      cloneDeep(testServers);
                                    newTestServers
                                      .filter((item) => item !== outbound)
                                      .map((item) => item);
                                    setTestServers(newTestServers);
                                    toast.error(e.message);
                                  });
                              }}
                              className="font-black text-lime-500 rounded-sm hover:bg-teal-50/30 w-12 h-6 text-center"
                            >
                              <span className="text-xs w-full">
                                {testServers.includes(outbound) ? (
                                  <Spinner size="sm" />
                                ) : outbounds![outbound]["history"][0] ? (
                                  outbounds![outbound]["history"][0].delay +
                                  "ms"
                                ) : (
                                  "-- ms"
                                )}
                              </span>
                            </span>
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

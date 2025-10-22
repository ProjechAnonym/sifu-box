import { useState,useEffect,useMemo } from "react";
import { Accordion, AccordionItem } from "@heroui/accordion";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Button } from "@heroui/button";
import { Badge } from "@heroui/badge";
import { Spinner } from "@heroui/spinner";
import NavBar from "@/components/navbar";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { FetchOutbounds } from "@/utils/outbound/fetch";
import { Outbound,OutboundGroup } from "@/types/outbound";
export default function OutboundsTable(props: {
  listen: string;
  secret: string;
  height: number;
  theme: string;
}){
    const { listen, secret, height, theme } = props;
    const [updateOutbounds, setUpdateOutbounds] = useState(true);
    const [outbounds, setOutbounds] = useState<{[key: string]: OutboundGroup | Outbound;} | null>(null);
    const [groups, setGroups] = useState<Array<string> | null>(null);
    const [test_groups, setTestGroups] = useState<Array<string>>([]);
    const [test_servers, setTestServers] = useState<Array<string>>([]);
    useEffect(() => {
      updateOutbounds && listen !== "" && secret !== "" && FetchOutbounds(listen, secret).then((res) => {
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
            .filter((outbound) => outbound.type === "Selector" || outbound.type === "Fallback" || outbound.type === "URLTest")
            .map((outbound) => outbound.name)
            .reverse()
        );
    }, [outbounds]);
    return (
    <div className="p-2" style={{ height: `calc(100% - ${height}px)` }}>
      {groups && groups.length > 0 && <NavBar groups={groups} />}
      {groups && 
        <ScrollShadow className="w-full h-full">
          <Accordion selectionMode="multiple" variant="bordered" isCompact>
          {groups.map((group, i) => (
            <AccordionItem
                aria-label={group}
                title={<span className="text-xl font-black">{group}</span>}
                key={`${group}-${i}`}
                id={group}
                subtitle={<span className="text-xs font-black select-none">点击展开</span>}
            >
              {outbounds![group].type !== "Fallback" && (
                <div className="flex flex-row gap-2">
                  <span className="text-xl font-black select-none">测试延迟</span>
                    <Button
                      isLoading={test_groups?.includes(group)}
                      onPress={() => {
                        const newTestGroup = cloneDeep(test_groups);
                        newTestGroup.push(group);
                        setTestGroups(newTestGroup);
                        // toast.promise(
                        //   GroupDelay(group, listen, secret, outbounds!),
                        //   {
                        //     loading: "loading",
                        //     success: (res) => {
                        //       const newTestGroups = cloneDeep(testGroups);
                        //       newTestGroups
                        //         .filter((item) => item !== group)
                        //         .map((item) => item);
                        //       setTestGroups(newTestGroups);
                        //       setOutbounds(res);
                        //       return `测试${group}完成`;
                        //     },
                        //     error: () => {
                        //       const newTestGroups = cloneDeep(testGroups);
                        //       newTestGroups
                        //         .filter((item) => item !== group)
                        //         .map((item) => item);
                        //       setTestGroups(newTestGroups);
                        //       return `测试${group}失败`;
                        //     },
                        //   }
                        // );
                      }}
                      isIconOnly
                      size="sm"
                      variant="light"
                  >
                    <i className="bi bi-stopwatch-fill text-lg" />
                  </Button>
                </div>
              )}
            </AccordionItem>
            ))
          }
          </Accordion>
        </ScrollShadow>}
    </div>
    )
}
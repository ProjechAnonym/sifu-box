import { useState,useEffect,useMemo,useCallback } from "react";
import { Accordion, AccordionItem } from "@heroui/accordion";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Button } from "@heroui/button";
import { Badge } from "@heroui/badge";
import { Spinner } from "@heroui/spinner";
import NavBar from "@/components/navbar";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { GroupDelay } from "@/utils/outbound/delay";
import { FetchOutbounds } from "@/utils/outbound/fetch";
import { Outbound,OutboundGroup } from "@/types/outbound";
export default function OutboundsTable(props: {
  listen: string;
  secret: string;
  height: number;
  theme: string;
}){
    const { listen, secret, height, theme } = props;
    const [update_outbounds, setUpdateOutbounds] = useState(true);
    const [outbounds, setOutbounds] = useState<{[key: string]: OutboundGroup | Outbound;} | null>(null);
    const [groups, setGroups] = useState<Array<string> | null>(null);
    const [test_groups, setTestGroups] = useState<Array<string>>([]);
    const [test_servers, setTestServers] = useState<Array<string>>([]);
    useEffect(() => {
      update_outbounds && listen !== "" && secret !== "" && FetchOutbounds(listen, secret).then((res) => {
        setUpdateOutbounds(false);
        res ? setOutbounds(res.proxies) : toast.error("获取代理失败");
      })
      .catch((e) => {
        setUpdateOutbounds(false);
        toast.error(e.code === "ERR_NETWORK" ? "请检查网络连接" : e.message);
      });
    }, [listen, secret, update_outbounds]);
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
                      isLoading={test_groups && test_groups.includes(group)}
                      onPress={() => {
                        const new_test_group = cloneDeep(test_groups);
                        new_test_group.push(group);
                        setTestGroups(new_test_group);
                        // toast.promise(
                        //   GroupDelay(group, listen, secret, outbounds!),
                        //   {
                        //     loading: "loading",
                        //     success: (res) => {
                        //       const new_test_groups = cloneDeep(test_groups);
                        //       setTestGroups(new_test_groups.filter((item) => item !== group).map((item) => item));
                        //       setOutbounds(res);
                        //       return `测试${group}完成`;
                        //     },
                        //     error: () => {
                        //       const new_test_groups = cloneDeep(test_groups);
                        //       setTestGroups(new_test_groups.filter((item) => item !== group).map((item) => item));
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
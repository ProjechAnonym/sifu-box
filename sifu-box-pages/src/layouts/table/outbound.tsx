import { useState,useEffect,useMemo,useCallback } from "react";
import { Accordion, AccordionItem } from "@heroui/accordion";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Button } from "@heroui/button";
import { Spinner } from "@heroui/spinner";
import OutboundsCard from "@/components/card/outbound";
import NavBar from "@/components/navbar";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { GroupDelay, OutboundDelay } from "@/utils/outbound/delay";
import { FetchOutbounds } from "@/utils/outbound/fetch";
import { Outbound,OutboundGroup } from "@/types/outbound";
import { SwitchOutbound } from "@/utils/outbound/switch";
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
    const switchOutbound = useCallback((tag: string, group: string) => {
      console.log(tag, group);
      outbounds && SwitchOutbound(group, tag, listen, secret, outbounds).then((res) => setOutbounds(res)).
        catch((e) => toast.error(
          e.code === "ERR_NETWORK"
          ? "请检查网络连接"
          : e.response.data.message
            ? e.response.data.message
            : e.response.data));
    }, [outbounds]);
    const testOutboundDelay = useCallback((tag: string) => {
      setTestServers(test_servers.concat(tag));
      outbounds && OutboundDelay(tag, listen, secret, outbounds).then((res) => {
        setTestServers(test_servers.filter((item) => item !== tag).map((item) => item))
        setOutbounds(res)
      }).catch((e) => {
        setTestServers(test_servers.filter((item) => item !== tag).map((item) => item))
        toast.error(
          e.code === "ERR_NETWORK"
          ? "请检查网络连接"
          : e.response.data.message
            ? e.response.data.message
            : e.response.data)}
      );
    }, [outbounds]);
    const handleGroupRes = (res: { [key: string]: OutboundGroup | Outbound; } | null, group: string) => {
      res && setOutbounds(res);
      return `测试"${group}"组${res ? "完成" : "失败"}`;
    }
   
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
                      onPress={() => toast.promise(
                          GroupDelay(group, listen, secret, outbounds!),
                          {
                            loading: "loading",
                            success: (res) => handleGroupRes(res, group),
                            error: () => handleGroupRes(null, group),
                          }
                        )
                      }
                      isIconOnly
                      size="sm"
                      variant="light"
                  >
                    <i className="bi bi-stopwatch-fill text-lg" />
                  </Button>
                </div>
              )}
              <ScrollShadow className="flex flex-wrap gap-x-3 gap-y-2 max-h-56 py-2">
                {(outbounds![group] as OutboundGroup).all.map((tag: string, j: number)=><OutboundsCard load={test_servers && test_servers.includes(tag)} theme={theme} group={group} tag={tag} outbound_msg={outbounds![tag] as Outbound} key={`${group}-${tag}-${j}`} select={tag === (outbounds![group] as OutboundGroup).now} switchOutbound={switchOutbound} testDelay={testOutboundDelay}/>)}
              </ScrollShadow>
            </AccordionItem>
            ))
          }
          </Accordion>
        </ScrollShadow>}
    </div>
    )
}
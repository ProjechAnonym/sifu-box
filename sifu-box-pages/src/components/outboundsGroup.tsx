import { useState, useEffect } from "react";
import { Select, SelectItem } from "@heroui/select"
import { Input } from "@heroui/input"
import { Button } from "@heroui/button";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Divider } from "@heroui/divider";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { SharedSelection } from "@heroui/system";

export default function OutboundsGroup(props: {theme: string, providers: string[], outbounds_group: {tag: string, type:string, providers: string[], tag_groups: string[]}[],setOutboundsGroup: (outboundsGroup: {tag: string, type:string, providers: string[], tag_groups: string[]}[]) => void}) {
    const { theme, providers, outbounds_group, setOutboundsGroup } = props;
    const [outbounds_groups, setOutboundsGroups] = useState<Array<{type: string, tag: string, providers: SharedSelection, tag_groups: SharedSelection}>>(outbounds_group ? outbounds_group.map((group: {tag: string, type:string, providers: string[], tag_groups: string[]}) => {return {tag: group.tag, type: group.type, providers: new Set(group.providers), tag_groups: new Set(group.tag_groups)}}) : [{type: "direct", tag: "direct", providers: new Set<string>(), tag_groups: new Set<string>()}]);
    useEffect(()=>{
        setOutboundsGroup(outbounds_groups.map((group)=>{
            return {
                tag: group.tag,
                type: group.type,
                tag_groups: Array.from(new Set(Object.entries(group.tag_groups).map(([_, value]) => (value)))),
                providers: Array.from(new Set(Object.entries(group.providers).map(([_, value]) => (value))))}
        }))
    },[outbounds_groups])
    return (
        <ScrollShadow className="gap-2 flex flex-col h-44">
            <header className="flex flex-row gap-2">
                <Button size="sm" isIconOnly 
                    onPress={() => {
                    if (outbounds_groups.length == 1) {
                      toast.error("无法继续删除");
                      return;
                    }
                    const temp_outbounds_groups = cloneDeep(outbounds_groups);
                    temp_outbounds_groups.pop();
                    setOutboundsGroups(temp_outbounds_groups);
                }}>
                  <i className="bi bi-dash text-3xl" />
                </Button>
                <Button size="sm" isIconOnly 
                  onPress={() => {
                  const temp_outbounds_groups = cloneDeep(outbounds_groups);
                    temp_outbounds_groups.push({
                      tag: "",
                      type: "",
                      providers: new Set<string>(),
                      tag_groups: new Set<string>(),
                    });
                    setOutboundsGroups(temp_outbounds_groups);
                  }}
                >
                  <i className="bi bi-plus text-3xl" />
                </Button></header>
            {outbounds_groups.map((group, i) => (
                <div key={`${group.tag}-${i}`} className="flex flex-wrap gap-1">
                    <Select size="sm" className="w-5/12" label="type" classNames={{popoverContent: `${theme} bg-content1 text-foreground`}} selectedKeys={[group.type]} onChange={e => setOutboundsGroups(outbounds_groups.map((outbounds_group, j) => i === j ? {...outbounds_group, type: e.target.value} : outbounds_group))}>
                        <SelectItem key="direct">direct</SelectItem>
                        <SelectItem key="selector">selector</SelectItem>
                        <SelectItem key="urltest">urltest</SelectItem>
                    </Select>
                    <Input size="sm" className="w-5/12" label="tag" value={group.tag} onValueChange={(value) => setOutboundsGroups(outbounds_groups.map((outbounds_group, j) => i === j ? {...outbounds_group, tag: value} : outbounds_group))}/>
                    {group.type !== "direct" && <Select size="sm" className="w-5/6" label="providers" classNames={{popoverContent: `${theme} bg-content1 text-foreground`}} selectedKeys={group.providers} selectionMode="multiple" onSelectionChange={keys=>setOutboundsGroups(outbounds_groups.map((outbounds_group, j) => i === j ? {...outbounds_group, providers: keys} : outbounds_group))}>
                        {providers.map((provider) => (<SelectItem key={provider}>{provider}</SelectItem>))}
                    </Select>}
                    {group.type !== "direct" && <Select size="sm" className="w-5/6" label="build-in" classNames={{popoverContent: `${theme} bg-content1 text-foreground`}} selectedKeys={group.tag_groups} selectionMode="multiple" onSelectionChange={keys=>setOutboundsGroups(outbounds_groups.map((outbounds_group, j) => i === j ? {...outbounds_group, tag_groups: keys} : outbounds_group))}>
                        {outbounds_groups.map((outbounds_group) => (<SelectItem key={`${outbounds_group.tag}`}>{outbounds_group.tag}</SelectItem>))}
                    </Select>}
                    <Divider />
                </div>
            ))}
            
        </ScrollShadow>
    )
}
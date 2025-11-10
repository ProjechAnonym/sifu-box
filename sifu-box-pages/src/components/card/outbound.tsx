import { Badge } from "@heroui/badge";
import { Spinner } from "@heroui/spinner";
import { Outbound } from "@/types/singbox/outbound";

export default function OutboundsCard(props: { theme: string, tag: string, group: string, outbound_msg: Outbound, select: boolean, load: boolean, SwitchOutbound: (tag:string, group:string) => void, TestDelay: (tag: string) => void}) {
    const {theme, tag, group ,outbound_msg, select, SwitchOutbound, TestDelay, load} = props;
    return (
        <Badge showOutline={false}
            placement="top-right"
            content={outbound_msg.udp ? 
                (<span className="text-teal-50 text-sm font-black">
                    UDP
                </span>) : (
                <del className="text-teal-50 text-sm font-black">
                    UDP
                </del>
                )
            }
            color={outbound_msg.udp ? "primary" : "danger"}
        >
            <div className={`${theme} ${select? `bg-content2` : `bg-content1`} p-2 w-40 cursor-pointer rounded-xl flex flex-col`}
                onClick={(e) => {
                    e.stopPropagation();
                    SwitchOutbound(tag, group);
                }}
            >
                 <span className="text-md font-black">
                    {tag}
                </span>
                <span className="font-black text-lime-500 rounded-sm hover:bg-gray-500/30 w-12 h-6 text-center"
                    onClick={(e) => {
                        e.stopPropagation();
                        TestDelay(tag);
                    }}>
                    <span className="text-xs w-full">
                        {load ? <Spinner size="sm" /> : outbound_msg.history[0] ? `${outbound_msg.history[0].delay}ms` : `--ms`}
                    </span>
                </span>
            </div>
        </Badge>
    )
}
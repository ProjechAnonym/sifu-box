
import { SizeCalc } from "@/utils/singbox/general";
export default function MonitorBar(props: {memory: { inuse: number; oslimit: number }, network: { up: number; down: number }, statistics: { upload: number; download: number }}) { 
    const { memory, statistics, network } = props;
    return (
    <div className="flex flex-row justify-center items-center gap-x-1">
        <div className="bg-default flex flex-col justify-center items-center rounded-md px-1">
            <span className="font-black text-xs"><i className="bi bi-memory" /> 内存</span>
            <span className="font-black text-xs">{SizeCalc(memory.inuse)}</span>
        </div>
        <div className="bg-default flex flex-col justify-center items-center rounded-md px-1">
            <span className="font-black text-xs">下载总量</span>
            <span className="font-black text-xs">
            {SizeCalc(statistics.download)}
            </span>
        </div>
        <div className="bg-default flex flex-col justify-center items-center rounded-md px-1">
            <span className="font-black text-xs">上传总量</span>
            <span className="font-black text-xs">
            {SizeCalc(statistics.upload)}
            </span>
        </div>
        <div className="bg-default flex flex-col justify-center items-center rounded-md px-1">
            <span className="font-black text-xs">当前上传</span>
            <span className="font-black text-xs">{SizeCalc(network.up)}</span>
        </div>
        <div className="bg-default flex flex-col justify-center items-center rounded-md px-1">
            <span className="font-black text-xs">当前下载</span>
            <span className="font-black text-xs">{SizeCalc(network.down)}</span>
        </div>
    </div>
    )
}
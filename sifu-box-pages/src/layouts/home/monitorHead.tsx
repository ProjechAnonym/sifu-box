import { useRef, useEffect, useState, useCallback } from "react"
import toast from "react-hot-toast";
import { LogSocket } from "@/components/socket/log";
import { ConnectionSocket } from "@/components/socket/connection";
import { MemorySocket } from "@/components/socket/memory";
import { TrafficSocket } from "@/components/socket/traffic";
import ControlPanel from "@/components/singboxPanel";
import MonitorBar from "@/components/monitorBar";
import DataBoard from "@/components/dataBoard";
import { cloneDeep } from "lodash";
import { FetchRules } from "@/utils/singbox/fetch";
import { ParseConnection,ParseLog } from "@/utils/singbox/general";
import { ConnectionData , ConnectionColumns} from "@/types/singbox/connection";
import { ruleColumns } from "@/types/singbox/rules";
import { logsColumns } from "@/types/singbox/log";

export default function MonitorHead(props: {template: string, theme: string, admin: boolean, token: string, listen: string, secret: string, log: boolean, fetchHeight: (height: number) => void }) {
    const {template, theme, admin, token, listen, secret, log, fetchHeight} = props;
    const header_container = useRef<HTMLElement>(null);
    const [network, setNetwork] = useState({ up: 0, down: 0 });
    const [memory, setMemory] = useState({ inuse: 0, oslimit: 0 });
    const [statistics, setStatistics] = useState<{
        upload: number;
        download: number;
    }>({ upload: 0, download: 0 });
    const [connetcions, setConnetcions] = useState<Array<ConnectionColumns>>([]);
    const [disConnetcions, setDisConnetcions] = useState<Array<ConnectionColumns>>([]);
    const [rules, setRules] = useState<Array<ruleColumns>>([]);
    const [logs, setLogs] = useState<Array<logsColumns>>([]);
    
    useEffect(() => {
        listen !== "" && secret !== "" && FetchRules(listen, secret).then((res) => setRules(res)).catch(() => toast.error("获取规则列表失败"));
        header_container.current && fetchHeight(header_container.current.clientHeight);
    }, [header_container.current && header_container.current.clientHeight]);
    const memoryReceiver = useCallback((data: { inuse: number; oslimit: number }) => setMemory(data), []);
    const trafficReceiver = useCallback((data: { up: number; down: number; }) => setNetwork({ up: data.up, down: data.down }), []);
    const stasticsReceiver = useCallback((data: ConnectionData) => {
        const { aliveConnections, deadConnections } = ParseConnection(
            data.connections,
            connetcions
          );
          setConnetcions(cloneDeep(aliveConnections));
          setDisConnetcions(deadConnections);
          setStatistics({
            upload: data.uploadTotal,
            download: data.downloadTotal,
          });
    }, []);
    const logReceiver = useCallback((data: { type: string; payload: string }) => setLogs(ParseLog(data, logs)), [logs]);
    return (
        <header className="flex flex-wrap gap-2 p-3" ref={header_container}>
            {log && <LogSocket secret={secret} listen={listen} receiver={logReceiver}/>}
            <ConnectionSocket secret={secret} listen={listen} receiver={stasticsReceiver}/>
            <MemorySocket secret={secret} listen={listen} receiver={memoryReceiver}/>
            <TrafficSocket secret={secret} listen={listen} receiver={trafficReceiver}/>
            <MonitorBar memory={memory} statistics={statistics} network={network}/>
            <ControlPanel admin={admin} token={token} theme={theme}/>
            <DataBoard theme={theme} rules={rules} logs={logs} connection={connetcions} disconnection={disConnetcions} log={log} listen={listen} secret={secret} />
            <div className="flex flex-col px-2 bg-default rounded-md select-none" >
                <span className="text-sm font-black p-0">当前模板</span>
                <span className="text-xs font-black p-0">{template? template : "无" }</span>
            </div>
        </header>
    )
}
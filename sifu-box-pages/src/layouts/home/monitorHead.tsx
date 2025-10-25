import { useRef, useEffect, useState } from "react"
import ControlPanel from "@/components/singboxPanel";
export default function MonitorHead(props: {template: string, theme: string, admin: boolean, token: string, listen: string, secret: string, log: boolean, fetchHeight: (height: number) => void }) {
    const {template, theme, admin, token, listen, secret, log, fetchHeight} = props;
    const header_container = useRef<HTMLElement>(null);
    useEffect(() => { 
        header_container.current && fetchHeight(header_container.current.clientHeight);
    }, [header_container.current && header_container.current.clientHeight]);
    return (
        <header className="flex flex-wrap gap-2 p-3" ref={header_container}>
            <ControlPanel admin={admin} token={token} theme={theme}/>
        </header>
    )
}
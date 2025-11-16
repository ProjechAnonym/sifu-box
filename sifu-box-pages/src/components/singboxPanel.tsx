import { useState,useEffect } from "react";
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem } from "@heroui/dropdown";
import { Button } from "@heroui/button";
import { ControlSignal } from "@/utils/singbox/control";
import { Key } from "@react-types/shared";
import toast from "react-hot-toast";
export default function ControlPanel(props: {admin: boolean, token: string, theme: string, }) {
    const { admin, token, theme } = props;
    const [status, setStatus] = useState(false);
    const [check, setCheck] = useState(true);
    useEffect(() => {
        token !== "" && check &&
        ControlSignal(token, "check").then((res) => {
            setStatus(res);
            setCheck(false);
            })
            .catch(() => {
            setStatus(false);
            setCheck(false);
            });
    }, [token, check]);
    const Control = (action: Key) =>toast.promise(ControlSignal(token, (action as string)!), 
        {
            loading: "loading",
            success: (res) => {
                switch (action) {
                  case "check":
                    setStatus(res);
                    return `检查操作完成`;
                  case "boot":
                    setStatus(res);
                    return `启动操作完成`;
                  case "reload":
                    setStatus(res);
                    return `重载操作完成`;
                  case "stop":
                    setStatus(res);
                    return `关闭操作完成`;
                  default:
                    break;
                }
                return `操作完成`;
            },
            error: (e) => {
                setCheck(true);
                return e.code === "ERR_NETWORK"
                  ? "网络连接失败"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data;
            },
        }
    )
        
    
    return (
    <Dropdown classNames={{content: [theme, `bg-content1`, `text-foreground`]}}>
      <DropdownTrigger>
        <Button variant="shadow" size="md" radius="sm" color="default"><span className="text-md font-black">Sing-Box <span className={`${status ? "text-green-500" : "text-rose-600"}`}> ·</span></span></Button>
      </DropdownTrigger>
      <DropdownMenu aria-label="Sing-box Control Panel" onAction={Control} variant="shadow" disabledKeys={admin? [] : ["stop", "boot"]}>
        <DropdownItem key="boot" description="启动Sing-box"><span className="font-black">启动</span></DropdownItem>
        <DropdownItem key="reload" description="重载Sing-box配置文件"><span className="font-black">重载</span></DropdownItem>
        <DropdownItem key="check" description="检查Sing-box运行状态"><span className="font-black">检查</span></DropdownItem>
        <DropdownItem key="stop" description="关闭Sing-box" className="text-danger" color="danger">
          <span className="font-black">关闭</span>
        </DropdownItem>
      </DropdownMenu>
    </Dropdown>
   
  );
}
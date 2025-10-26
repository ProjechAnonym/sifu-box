import { useState } from "react";
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem } from "@heroui/dropdown";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/modal";
import { Button } from "@heroui/button";
import Rules from "./table/rule";
import Logs from "./table/logs";
import Connections from "./table/connection";
import { ruleColumns } from "@/types/singbox/rules";
import { logsColumns } from "@/types/singbox/log";
import { ConnectionColumns } from "@/types/singbox/connection";
export default function DataBoard(props: {theme: string, rules: Array<ruleColumns>, logs: Array<logsColumns>, connection: Array<ConnectionColumns>, disconnection:Array<ConnectionColumns>,log: boolean}) {
    const { theme, rules, logs, connection, disconnection, log} = props;
    const [data, setData] = useState("")
    const { isOpen, onClose, onOpen } = useDisclosure();
    return (
        <div>
            <Modal
                isOpen={isOpen}
                onClose={onClose}
                classNames={{ base: `${theme} bg-background text-foreground` }}
                size="xl"
            >
                <ModalContent>
                {(onClose) => (
                    <>
                    <ModalHeader></ModalHeader>
                    <ModalBody>
                        {data === "rules" && <Rules rules={rules} theme={theme} />}
                        {data === "logs" && <Logs theme={theme} logs={logs} />}
                        {data === "connections" && (<Connections theme={theme} connection={connection} disConnection={disconnection}/>
                )}
                    </ModalBody>
                    <ModalFooter>
                        <Button
                        variant="shadow"
                        size="sm"
                        color="danger"
                        onPress={onClose}
                        >
                        <span className="font-black text-lg">关闭</span>
                        </Button>
                    </ModalFooter>
                    </>
                )}
                </ModalContent>
            </Modal>
            <Dropdown classNames={{content: [theme, `bg-content1`, `text-foreground`]}}>   
            <DropdownTrigger>
                <Button variant="shadow" size="md" radius="sm" color="default"><span className="text-md font-black">数据列表</span></Button>
            </DropdownTrigger>
            <DropdownMenu aria-label="Databoard Panel" onAction={(key) => {setData(key as string);key !== "clear" && onOpen()}} variant="shadow" disabledKeys={log? [] : ["logs"]}>
                <DropdownItem key="clear" description="清空DNS缓存"><span className="font-black">清空</span></DropdownItem>
                <DropdownItem key="rules" description="查看规则"><span className="font-black">规则</span></DropdownItem>
                <DropdownItem key="logs" description="查看日志"><span className="font-black">日志</span></DropdownItem>
                <DropdownItem key="connections" description="查看连接"><span className="font-black">连接</span></DropdownItem>
            </DropdownMenu>
            </Dropdown>
    </div>
    )
}
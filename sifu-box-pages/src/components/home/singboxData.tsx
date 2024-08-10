import { useMemo, useState } from "react";
import {
  Card,
  CardHeader,
  CardBody,
  Divider,
  Button,
  Modal,
  ModalHeader,
  ModalContent,
  ModalBody,
  ModalFooter,
  useDisclosure,
  Switch,
} from "@nextui-org/react";
import MyTable from "../table";
import { ConnectionLog } from "@/types/singbox";
function SizeCalc(size: number) {
  return size < 1024
    ? `${size} B`
    : size < 1024 * 1024
    ? `${(size / 1024).toFixed(2)} KB`
    : size < 1024 * 1024 * 1024
    ? `${(size / (1024 * 1024)).toFixed(2)} MB`
    : `${(size / (1024 * 1024 * 1024)).toFixed()}`;
}
export default function SingboxData(props: {
  clearFakeip: () => void;
  network: { up: number; down: number };
  memory: { inuse: number; oslimit: number };
  networkTotal: { upload: number; download: number };
  dark: boolean;
  logs: {
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{
      type: string;
      payload: string;
      time: string;
      key: string;
    }>;
  };
  rules: {
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{
      type: string;
      payload: string;
      proxy: string;
      key: string;
    }>;
  };
  connections: {
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<ConnectionLog>;
  };
  disconnections: {
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<ConnectionLog>;
  };
}) {
  const {
    clearFakeip,
    network,
    memory,
    dark,
    networkTotal,
    rules,
    logs,
    connections,
    disconnections,
  } = props;
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const [modal, setModal] = useState("");
  const [connecting, setConnecting] = useState(true);
  const datas: {
    labels: Array<{
      label: string;
      key: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{ key: string; [addtionalProp: string]: any }>;
  } = useMemo(() => {
    switch (modal) {
      case "规则列表":
        const ruleValues = rules.values.map((value) => {
          return { ...value };
        });
        return { labels: rules.labels, values: ruleValues };
      case "日志列表":
        const logValues = logs.values.map((value) => {
          return { ...value };
        });
        return { labels: logs.labels, values: logValues };
      case "连接列表":
        const connectionValues = connecting
          ? connections.values.map((value) => {
              return { ...value };
            })
          : disconnections.values.map((value) => {
              return { ...value };
            });
        return { labels: connections.labels, values: connectionValues };
      default:
        break;
    }
    return {
      labels: [
        { label: "none", key: "none", allowSort: false, initShow: true },
      ],
      values: [{ key: "none", none: "none" }],
    };
  }, [modal, rules, logs, connections]);
  const searchField = useMemo(() => {
    switch (modal) {
      case "规则列表":
        return "payload";
      case "日志列表":
        return "payload";
      case "连接列表":
        return "host";
      default:
        break;
    }
    return "";
  }, [modal]);
  return (
    <header className="w-full flex flex-wrap gap-2 p-2 items-center">
      <Modal isOpen={isOpen} onOpenChange={onOpenChange} size="5xl">
        <ModalContent
          className={`${
            dark ? "sifudark" : "sifulight"
          } bg-background text-foreground`}
        >
          {(onClose) => (
            <>
              <ModalHeader>
                <span className="text-lg font-black">{modal}</span>
              </ModalHeader>
              <ModalBody>
                <div className="flex flex-row gap-2 items-center">
                  {modal === "连接列表" && (
                    <Switch
                      isSelected={connecting}
                      onValueChange={setConnecting}
                      size="sm"
                      color="success"
                      startContent={
                        <span>
                          <i className="bi bi-wifi" />
                        </span>
                      }
                      endContent={
                        <span>
                          <i className="bi bi-wifi-off" />
                        </span>
                      }
                    >
                      <span className="font-black">
                        连接状态:{`${connecting ? "活动" : "断开"}`}
                      </span>
                    </Switch>
                  )}
                </div>
                <MyTable
                  dark={dark}
                  data={datas}
                  defaultSearchField={searchField}
                  rowsPerPage={10}
                />
              </ModalBody>
              <ModalFooter>
                <Button size="sm" color="danger" onPress={onClose}>
                  <span className="text-lg font-black">关闭</span>
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
      <Button size="sm" color="primary" onPress={clearFakeip}>
        <span className="font-black text-lg">清空fakeip</span>
      </Button>
      <Button
        size="sm"
        color="primary"
        onPress={() => {
          onOpen();
          setModal("规则列表");
        }}
      >
        <span className="font-black text-lg">规则</span>
      </Button>
      <Button
        size="sm"
        color="primary"
        onPress={() => {
          onOpen();
          setModal("连接列表");
        }}
      >
        <span className="font-black text-lg">连接</span>
      </Button>
      <Button
        size="sm"
        color="primary"
        onPress={() => {
          onOpen();
          setModal("日志列表");
        }}
      >
        <span className="font-black text-lg">日志</span>
      </Button>
      <Card radius="sm">
        <CardHeader className={`${dark ? "bg-zinc-800" : "bg-slate-100"} p-1`}>
          <i className="bi bi-arrow-up-short text-xs" />
          <span className="text-xs font-black">总量</span>
        </CardHeader>
        <Divider />
        <CardBody
          className={`${dark ? "bg-zinc-800" : "bg-slate-100"} px-2 py-0`}
        >
          <span className="text-foreground text-xs">
            {SizeCalc(networkTotal.upload)}
          </span>
        </CardBody>
      </Card>
      <Card radius="sm">
        <CardHeader className={`${dark ? "bg-zinc-800" : "bg-slate-100"} p-1`}>
          <i className="bi bi-arrow-down-short text-xs" />
          <span className="text-xs font-black">总量</span>
        </CardHeader>
        <Divider />
        <CardBody
          className={`${dark ? "bg-zinc-800" : "bg-slate-100"} px-2 py-0`}
        >
          <span className="text-foreground text-xs">
            {SizeCalc(networkTotal.download)}
          </span>
        </CardBody>
      </Card>
      <Card radius="sm">
        <CardHeader className={`${dark ? "bg-zinc-800" : "bg-slate-100"} p-1`}>
          <i className="bi bi-memory text-xs" />
        </CardHeader>
        <Divider />
        <CardBody
          className={`${dark ? "bg-zinc-800" : "bg-slate-100"} px-2 py-0`}
        >
          <span className="text-foreground text-xs">
            {SizeCalc(memory.inuse)}
          </span>
        </CardBody>
      </Card>
      <Card radius="sm">
        <CardHeader className={`${dark ? "bg-zinc-800" : "bg-slate-100"} p-1`}>
          <i className="bi bi-arrow-up-short text-xs" />
        </CardHeader>
        <Divider />
        <CardBody
          className={`${dark ? "bg-zinc-800" : "bg-slate-100"} px-2 py-0`}
        >
          <span className="text-foreground text-xs">
            {SizeCalc(network.up)}/s
          </span>
        </CardBody>
      </Card>
      <Card radius="sm">
        <CardHeader className={`${dark ? "bg-zinc-800" : "bg-slate-100"} p-1`}>
          <i className="bi bi-arrow-down-short text-xs" />
        </CardHeader>
        <Divider />
        <CardBody
          className={`${dark ? "bg-zinc-800" : "bg-slate-100"} px-2 py-0`}
        >
          <span className="text-foreground text-xs">
            {SizeCalc(network.down)}/s
          </span>
        </CardBody>
      </Card>
    </header>
  );
}

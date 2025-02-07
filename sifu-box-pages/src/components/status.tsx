import { useState, useEffect } from "react";
import { Button, ButtonGroup } from "@heroui/button";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/modal";
import Rules from "./monitor/rules";
import toast from "react-hot-toast";
import { ConnectionSocket, Connections } from "./monitor/connections";
import { MemorySocket } from "./monitor/memory";
import { TrafficSocket } from "./monitor/traffic";
import { Logs, LogSocket } from "./monitor/logs";
import { cloneDeep } from "lodash";
import { ParseLog } from "@/utils/singbox/parse";
import { ClearFakeip } from "@/utils/singbox/clear";
import { FetchRules } from "@/utils/singbox/fetch";
import { ParseConnection } from "@/utils/singbox/parse";
import { sizeCalc } from "@/utils/singbox/calc";
import { ruleColumns } from "@/types/singbox/rules";
import { logsColumns } from "@/types/singbox/logs";
import { ConnectionColumns } from "@/types/singbox/connections";

export default function Status(props: {
  listen: string;
  secret: string;
  log: boolean;
  theme: string;
}) {
  const { listen, secret, log, theme } = props;
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [network, setNetwork] = useState({ up: 0, down: 0 });
  const [memory, setMemory] = useState({ inuse: 0, oslimit: 0 });
  const [statistics, setStatistics] = useState<{
    upload: number;
    download: number;
  }>({ upload: 0, download: 0 });
  const [info, setInfo] = useState("");
  const [rules, setRules] = useState<Array<ruleColumns>>([]);
  const [logs, setLogs] = useState<Array<logsColumns>>([]);
  const [connetcions, setConnetcions] = useState<Array<ConnectionColumns>>([]);
  const [disConnetcions, setDisConnetcions] = useState<
    Array<ConnectionColumns>
  >([]);

  useEffect(() => {
    listen !== "" &&
      secret !== "" &&
      FetchRules(listen, secret)
        .then((res) => setRules(res))
        .catch(() => toast.error("获取规则列表失败"));
  }, [listen, secret]);
  return (
    <header className="flex items-center flex-wrap gap-2 px-2">
      <MemorySocket
        secret={secret}
        listen={listen}
        receiver={(data) => setMemory(data)}
      />
      <TrafficSocket
        secret={secret}
        listen={listen}
        receiver={(data) => setNetwork(data)}
      />
      <ConnectionSocket
        listen={listen}
        secret={secret}
        receiver={(data) => {
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
        }}
      />
      {log && listen !== "" && secret !== "" && (
        <LogSocket
          listen={listen}
          secret={secret}
          receiver={(data) => setLogs(ParseLog(data, logs))}
        />
      )}
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
                {info === "rules" && <Rules rules={rules} theme={theme} />}
                {info === "logs" && <Logs theme={theme} logs={logs} />}
                {info === "connections" && (
                  <Connections
                    theme={theme}
                    connection={connetcions}
                    disConnection={disConnetcions}
                  />
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
      <ButtonGroup>
        <Button
          color="primary"
          size="sm"
          variant="shadow"
          onPress={() =>
            toast.promise(ClearFakeip(listen, secret), {
              loading: "清空FakeIP中...",
              success: "清除FakeIP完成",
              error: (e) =>
                e.code === "ERR_NETWORK"
                  ? "请检查网络连接"
                  : e.response.data.message
                    ? e.response.data.message
                    : e.response.data,
            })
          }
        >
          <span className="font-black text-lg">清空</span>
        </Button>
        <Button
          color="primary"
          size="sm"
          variant="shadow"
          onPress={() => {
            setInfo("rules");
            onOpen();
          }}
        >
          <span className="font-black text-lg">规则</span>
        </Button>
        <Button
          color="primary"
          size="sm"
          variant="shadow"
          onPress={() => {
            setInfo("connections");
            onOpen();
          }}
        >
          <span className="font-black text-lg">连接</span>
        </Button>
        {log && (
          <Button
            color="primary"
            size="sm"
            variant="shadow"
            onPress={() => {
              setInfo("logs");
              onOpen();
            }}
          >
            <span className="font-black text-lg">日志</span>
          </Button>
        )}
      </ButtonGroup>
      <div className="bg-content1 flex flex-col justify-center items-center px-2 rounded-md">
        <p className="flex flex-row justify-center items-center">
          <i className="bi bi-memory" />
          <span className="font-black text-sm">内存</span>
        </p>
        <span className="font-black text-xs">{sizeCalc(memory.inuse)}</span>
      </div>
      <div className="bg-content1 flex flex-col justify-center items-center px-2 rounded-md">
        <span className="font-black text-sm">下载总量</span>
        <span className="font-black text-sm">
          {sizeCalc(statistics.download)}
        </span>
      </div>
      <div className="bg-content1 flex flex-col justify-center items-center px-2 rounded-md">
        <span className="font-black text-sm">上传总量</span>
        <span className="font-black text-sm">
          {sizeCalc(statistics.upload)}
        </span>
      </div>
      <div className="bg-content1 flex flex-col justify-center items-center px-2 rounded-md">
        <span className="font-black text-sm">当前上传</span>
        <span className="font-black text-sm">{sizeCalc(network.up)}</span>
      </div>
      <div className="bg-content1 flex flex-col justify-center items-center px-2 rounded-md">
        <span className="font-black text-sm">当前下载</span>
        <span className="font-black text-sm">{sizeCalc(network.down)}</span>
      </div>
    </header>
  );
}

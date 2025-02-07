import { useCallback, useEffect, useState } from "react";
import { Switch } from "@heroui/switch";
import useWebSocket, { ReadyState } from "react-use-websocket";
import toast from "react-hot-toast";
import CustomTable from "../table";
import { SocketUrl } from "@/utils/singbox/socket";
import {
  ConnectionData,
  MetaData,
  ConnectionColumns,
  initConnectionLabels,
} from "@/types/singbox/connections";
export function ConnectionSocket(props: {
  secret: string;
  listen: string;
  receiver: (data: ConnectionData) => void;
}) {
  const { secret, listen, receiver } = props;
  const connectionURL = useCallback(
    () =>
      listen !== "" && secret !== ""
        ? SocketUrl(listen, "connections", [{ token: secret }])
        : "",
    [listen, secret]
  );
  const { lastJsonMessage, readyState } = useWebSocket(connectionURL());
  useEffect(() => {
    isConnectionMessage(lastJsonMessage) &&
      receiver(lastJsonMessage as ConnectionData);
    readyState === ReadyState.CLOSED ||
      (readyState === ReadyState.CLOSING && toast.error("连接Socket关闭"));
  }, [readyState, lastJsonMessage]);
  return <></>;
}
function isConnectionMessage(data: any): data is ConnectionData {
  return (
    typeof data === "object" &&
    data !== null &&
    Array.isArray(data.connections) &&
    data.connections.every(
      (connection: any) =>
        typeof connection === "object" &&
        connection !== null &&
        Array.isArray(connection.chains) &&
        connection.chains.every((chain: any) => typeof chain === "string") &&
        typeof connection.download === "number" &&
        typeof connection.upload === "number" &&
        typeof connection.id === "string" &&
        typeof connection.rule === "string" &&
        typeof connection.rulePayload === "string" &&
        typeof connection.start === "string" &&
        isMetaData(connection.metadata)
    ) &&
    typeof data.downloadTotal === "number" &&
    typeof data.uploadTotal === "number" &&
    typeof data.memory === "number"
  );
}
function isMetaData(data: any): data is MetaData {
  return (
    typeof data === "object" &&
    data !== null &&
    typeof data.destinationIP === "string" &&
    typeof data.destinationPort === "string" &&
    typeof data.dnsMode === "string" &&
    typeof data.host === "string" &&
    typeof data.network === "string" &&
    typeof data.processPath === "string" &&
    typeof data.sourceIP === "string" &&
    typeof data.sourcePort === "string" &&
    typeof data.type === "string"
  );
}
export function Connections(props: {
  theme: string;
  connection: Array<ConnectionColumns>;
  disConnection: Array<ConnectionColumns>;
}) {
  const { theme, connection, disConnection } = props;
  useEffect(() => {}, [connection, disConnection]);
  const [status, setStatus] = useState(true);
  return (
    <div className="flex flex-col gap-2">
      <Switch
        isSelected={status}
        onValueChange={setStatus}
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
        {status ? (
          <span className="font-black text-md">活动连接</span>
        ) : (
          <span className="font-black text-md">断开连接</span>
        )}
      </Switch>
      <CustomTable
        theme={theme}
        data={{
          labels: initConnectionLabels,
          values: status ? connection : disConnection,
        }}
      />
    </div>
  );
}

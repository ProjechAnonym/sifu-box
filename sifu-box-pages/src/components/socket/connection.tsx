import { useEffect, useCallback } from "react";
import toast from "react-hot-toast";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { SocketUrl } from "@/utils/singbox/general";
import { ConnectionData, MetaData } from "@/types/singbox/connection";
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
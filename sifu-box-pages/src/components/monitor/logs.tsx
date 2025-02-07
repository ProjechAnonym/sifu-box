import { useEffect, useCallback } from "react";
import CustomTable from "../table";
import toast from "react-hot-toast";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { SocketUrl } from "@/utils/singbox/socket";
import { initLogsLabels, logsColumns } from "@/types/singbox/logs";

export function LogSocket(props: {
  secret: string;
  listen: string;
  receiver: (data: { type: string; payload: string }) => void;
}) {
  const { secret, listen, receiver } = props;
  const logURL = useCallback(
    () =>
      listen !== "" && secret !== ""
        ? SocketUrl(listen, "logs", [{ level: "info", token: secret }])
        : "",
    [listen, secret]
  );

  const { lastJsonMessage, readyState } = useWebSocket(logURL());
  useEffect(() => {
    isLogMessage(lastJsonMessage) &&
      receiver(lastJsonMessage as { type: string; payload: string });
    readyState === ReadyState.CLOSED ||
      (readyState === ReadyState.CLOSING && toast.error("日志Socket关闭"));
  }, [readyState, lastJsonMessage]);
  return <></>;
}
const isLogMessage = (
  message: any
): message is { type: string; payload: string } => {
  return (
    typeof message === "object" &&
    message !== null &&
    typeof message.type === "string" &&
    typeof message.payload === "string"
  );
};
export function Logs(props: { theme: string; logs: Array<logsColumns> }) {
  const { theme, logs } = props;
  return (
    <div>
      <CustomTable
        theme={theme}
        data={{ labels: initLogsLabels, values: logs }}
        defaultSearchField="payload"
      />
    </div>
  );
}

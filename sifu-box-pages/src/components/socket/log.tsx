import { useEffect, useCallback } from "react";
import toast from "react-hot-toast";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { SocketUrl } from "@/utils/singbox/general";
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
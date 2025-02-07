import { useEffect, useCallback } from "react";
import toast from "react-hot-toast";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { SocketUrl } from "@/utils/singbox/socket";

export function MemorySocket(props: {
  secret: string;
  listen: string;
  receiver: (data: { inuse: number; oslimit: number }) => void;
}) {
  const { secret, listen, receiver } = props;
  const memoryURL = useCallback(
    () =>
      listen !== "" && secret !== ""
        ? SocketUrl(listen, "memory", [{ token: secret }])
        : "",
    [listen, secret]
  );

  const { lastJsonMessage, readyState } = useWebSocket(memoryURL());
  useEffect(() => {
    isMemoryMessage(lastJsonMessage) &&
      receiver(lastJsonMessage as { inuse: number; oslimit: number });
    readyState === ReadyState.CLOSED ||
      (readyState === ReadyState.CLOSING && toast.error("内存Socket关闭"));
  }, [readyState, lastJsonMessage]);
  return <></>;
}

function isMemoryMessage(
  message: any
): message is { inuse: number; oslimit: number } {
  return (
    typeof message === "object" &&
    message !== null &&
    typeof message.inuse === "number" &&
    typeof message.oslimit === "number"
  );
}

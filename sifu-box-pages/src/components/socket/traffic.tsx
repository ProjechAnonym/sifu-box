import { useEffect, useCallback } from "react";
import toast from "react-hot-toast";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { SocketUrl } from "@/utils/singbox/general";

export function TrafficSocket(props: {
  secret: string;
  listen: string;
  receiver: (data: { up: 0; down: 0 }) => void;
}) {
  const { secret, listen, receiver } = props;
  const trafficURL = useCallback(
    () =>
      listen !== "" && secret !== ""
        ? SocketUrl(listen, "traffic", [{ token: secret }])
        : "",
    [listen, secret]
  );

  const { lastJsonMessage, readyState } = useWebSocket(trafficURL());
  useEffect(() => {
    isTrafficMessage(lastJsonMessage) &&
      receiver(lastJsonMessage as { up: 0; down: 0 });
    readyState === ReadyState.CLOSED ||
      (readyState === ReadyState.CLOSING && toast.error("网络Socket关闭"));
  }, [readyState, lastJsonMessage]);
  return <></>;
}
function isTrafficMessage(message: any): message is { up: 0; down: 0 } {
  return (
    typeof message === "object" &&
    message !== null &&
    typeof message.up === "number" &&
    typeof message.down === "number"
  );
}
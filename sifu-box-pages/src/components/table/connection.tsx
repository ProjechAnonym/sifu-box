import { useEffect, useState } from "react";
import { Switch } from "@heroui/switch";
import MyTable from "@/layouts/mytable";
import { ConnectionColumns, init_connection_labels } from "@/types/singbox/connection";
export default function Connections(props: {
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
      <MyTable
        theme={theme}
        data={{
          labels: init_connection_labels,
          values: status ? connection : disConnection,
        }}
      />
    </div>
  );
}
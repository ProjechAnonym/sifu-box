import MyTable from "@/layouts/mytable";
import { logsColumns, init_logs_labels } from "@/types/singbox/log";
export default function Logs(props: { theme: string; logs: Array<logsColumns> }) {
  const { theme, logs } = props;
  return (
    <div>
      <MyTable
        theme={theme}
        data={{ labels: init_logs_labels, values: logs }}
        defaultSearchField="payload"
      />
    </div>
  );
}
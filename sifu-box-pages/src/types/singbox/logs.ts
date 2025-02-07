export interface logsColumns {
  type: string;
  payload: string;
  time: string;
  key: string;
}
export const initLogsLabels = [
  { label: "时间", key: "time", allowSort: true, initShow: true },
  { label: "等级", key: "type", allowSort: true, initShow: true },
  { label: "信息", key: "payload", allowSort: true, initShow: true },
];

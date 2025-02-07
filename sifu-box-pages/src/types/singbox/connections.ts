export const initConnectionLabels = [
  { key: "id", label: "ID", allowSort: true, initShow: false },
  { key: "rule", label: "命中规则", allowSort: true, initShow: true },
  {
    key: "rulePayload",
    label: "规则信息",
    allowSort: true,
    initShow: false,
  },
  { key: "start", label: "时间", allowSort: true, initShow: true },
  { key: "upload", label: "上传", allowSort: true, initShow: true },
  { key: "download", label: "下载", allowSort: true, initShow: true },
  { key: "chains", label: "节点链", allowSort: true, initShow: true },
  {
    key: "destinationIP",
    label: "目标IP",
    allowSort: true,
    initShow: true,
  },
  {
    key: "destinationPort",
    label: "目标端口",
    allowSort: true,
    initShow: false,
  },
  {
    key: "dnsMode",
    label: "dns模式",
    allowSort: false,
    initShow: false,
  },
  {
    key: "network",
    label: "网络",
    allowSort: false,
    initShow: false,
  },
  {
    key: "processPath",
    label: "进程",
    allowSort: true,
    initShow: false,
  },
  {
    key: "sourceIP",
    label: "源IP",
    allowSort: true,
    initShow: false,
  },
  {
    key: "sourcePort",
    label: "源端口",
    allowSort: true,
    initShow: false,
  },
  {
    key: "type",
    label: "类型",
    allowSort: false,
    initShow: false,
  },
  { key: "host", label: "嗅探域名", allowSort: true, initShow: true },
];
export interface MetaData {
  destinationIP: string;
  destinationPort: string;
  dnsMode: string;
  host: string;
  network: string;
  processPath: string;
  sourceIP: string;
  sourcePort: string;
  type: string;
}
export interface ConnectionColumns extends MetaData {
  chains: string;
  download: string;
  upload: string;
  id: string;
  rule: string;
  rulePayload: string;
  start: string;
  key: string;
}
export interface ConnectionData {
  connections: Array<{
    chains: Array<string>;
    download: number;
    upload: number;
    id: string;
    rule: string;
    rulePayload: string;
    start: string;
    metadata: MetaData;
  }>;
  downloadTotal: number;
  uploadTotal: number;
  memory: number;
}

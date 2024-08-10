export interface ServerValue {
  history: [{ time: string; delay: number }];
  name: string;
  type: string;
  udp: boolean;
}
export interface ServerGroupValue {
  all: Array<string>;
  history: [{ time: string; delay: number }];
  name: string;
  type: string;
  udp: boolean;
  now: string;
}
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
export interface ConnectionLog extends MetaData {
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

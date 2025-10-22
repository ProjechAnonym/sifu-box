export interface Outbound {
  history: [{ time: string; delay: number }];
  name: string;
  type: string;
  udp: boolean;
}
export interface OutboundGroup {
  all: Array<string>;
  history: [{ time: string; delay: number }];
  name: string;
  type: string;
  udp: boolean;
  now: string;
}
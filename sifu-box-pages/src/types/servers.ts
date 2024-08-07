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

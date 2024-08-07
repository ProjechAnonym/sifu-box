export interface HostValue {
  key: string;
  url: string;
  username: string;
  password: string;
  localhost: boolean;
  config: string;
  secret: string;
  port: number;
}
export type HostSendData = Omit<HostValue, "key" | "localhost" | "config">;
export type HostRevData = Omit<HostValue, "key">;

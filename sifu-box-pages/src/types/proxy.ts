export interface ProviderValue {
  id: number;
  path: string;
  proxy: boolean;
  name: string;
  remote: boolean;
}
export interface RulesetValue {
  id: number;
  type: string;
  url: string;
  path: string;
  format: string;
  tag: string;
  download_detour: string;
  update_interval: string;
  dnsRule: string;
  china: boolean;
  label: string;
}
export type ProviderData = Omit<ProviderValue, "id">;
export type RulesetRemoteData = Omit<RulesetValue, "id" | "path">;
export type RulesetLocalData = Omit<
  RulesetValue,
  "id" | "url" | "download_detour" | "update_interval"
>;

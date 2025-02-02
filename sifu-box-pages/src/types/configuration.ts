export interface provider {
  name: string;
  path: string;
  detour: string;
  remote: boolean;
}
export interface ruleset {
  tag: string;
  type: string;
  path: string;
  format: string;
  china: boolean;
  name_server: string;
  label: string;
  download_detour: string;
  update_interval: string;
}

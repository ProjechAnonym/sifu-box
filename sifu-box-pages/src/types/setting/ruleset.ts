export interface RuleSet {
    id: number,
    name: string,
    path: string,
    remote: boolean,
    update_interval: string,
    binary: boolean,
    download_detour: string
}
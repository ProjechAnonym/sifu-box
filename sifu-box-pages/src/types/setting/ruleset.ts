export interface RuleSet {
    id: number,
    name: string,
    path: string,
    remote: boolean,
    update_interval: string,
    binary: boolean,
    download_detour: string
}

export const DEFAULT_RULESET = {
    name: '',
    path: '',
    remote: false,
    update_interval: '',
    binary: false,
    download_detour: ''
}
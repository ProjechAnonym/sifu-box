export interface Provider {
    id: number;
    name: string;
    path: string;
    remote: boolean;
}

export const DEFAULT_PROVIDER = {
    name: "",
    path: "",
    remote: false
}
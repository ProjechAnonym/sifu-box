import { useRef } from "react";
import { Autocomplete, AutocompleteItem } from "@heroui/autocomplete";
import { Button,ButtonGroup } from "@heroui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@heroui/popover";
import Interval from "@/components/card/interval";
import toast from "react-hot-toast";
import copy from "copy-to-clipboard";
import { Refresh } from "@/utils/configuration/refresh";
import { Export } from "@/utils/migrate/export";
import { Import } from "@/utils/migrate/import";
export default function SettingHead(props: {
  token: string;
  admin: boolean;
  theme: string;
  files: Array<{label: string, path: string}>;
  setUpdate: (update: boolean) => void;
}) {
    const { token, admin, theme, files ,setUpdate } = props;
    const file_input = useRef<HTMLInputElement>(null);
    console.log(files)
    const refresh = () => toast.promise(Refresh(token), {
            loading: "更新配置文件中...",
            success: (res) => {          
                Array.isArray(res) && res.every(item => typeof item === 'object') && res.forEach(item => 'message' in item && typeof item.message === 'string' ? toast.error(item.message) : toast.error("未知错误"))
                return "更新配置文件完成"
            },
            error: (e) => {
            if (e.code === "ERR_NETWORK") {
                return "请检查网络连接";
            }
            if (e.response.data.message) {
                if (typeof e.response.data.message === "string") {
                    return e.response.data.message;
                }
                (e.response.data.message as Array<string>).map((m) =>
                    toast.error(m)
                );
                return "";
            }
            return e.response.data;
            },
        })
    return (
        <header className="flex flex-wrap gap-1 p-1"> 
            <ButtonGroup>
                <Button
                    variant="shadow"
                    color="primary"
                    size="sm"
                    onPress = {refresh}
                    isDisabled={!admin}
                >
                    <span className="font-black text-lg">刷新</span>
                </Button>
                <Popover
                    placement="bottom"
                    classNames={{ content: `${theme} bg-content1 text-foreground` }}
                >
                    <PopoverTrigger>
                        <Button
                            variant="shadow"
                            color="primary"
                            size="sm"
                            isDisabled={!admin}
                        >
                            <span className="font-black text-lg">定时</span>
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-64 h-36">
                        <Interval theme={theme} token={token} />
                    </PopoverContent>
                </Popover>
                <Button
                variant="shadow"
                color="primary"
                size="sm"
                isDisabled={!admin}
                onPress={() =>
                    toast.promise(Export(token), {
                    loading: "导出配置中...",
                    success: (res) => (res ? "配置导出完成" : "配置导出失败"),
                    error: (e) =>
                        e.code === "ERR_NETWORK"
                        ? "请检查网络连接"
                        : e.response.data.message
                            ? e.response.data.message
                            : e.response.data,
                    })
                }
                >
                <span className="font-black text-lg">导出</span>
                </Button>
                <Button
                    variant="shadow"
                    color="primary"
                    size="sm"
                    isDisabled={!admin}
                    onPress={() => file_input.current && file_input.current.click()}
                    >
                    <span className="font-black text-lg">导入</span>
                </Button>
            </ButtonGroup>
             <input
                className="hidden"
                type="file"
                ref={file_input}
                onChange={(e) =>
                e.target.files &&
                toast.promise(Import(token, e.target.files[0]), {
                    loading: "恢复配置中...",
                    success: (res) => {
                        setUpdate(true);
                        Array.isArray(res) && 
                        res.every(item => typeof item === 'object') && 
                        res.forEach(
                            item => 
                                'message' in item && 
                                "status" in item && 
                                typeof item.message === 'string' && 
                                typeof item.status === "boolean" ? 
                                (item.status ? toast.success(item.message) : toast.error(item.message)) : toast.error("未知错误"))
                        return "导入配置完成";
                    },
                    error: (e) => {
                    file_input.current && (file_input.current.value = "");
                    setUpdate(true);
                    return e.code === "ERR_NETWORK"
                        ? "请检查网络连接"
                        : e.response.data.message
                        ? e.response.data.message
                        : e.response.data;
                    },
                })
                }
            />
            {files && (
                <Autocomplete
                    variant="underlined"
                    label={<span className="text-xs font-black">文件链接</span>}
                    size="sm"
                    classNames={{
                        popoverContent: `${theme} bg-content1 text-foreground`,
                    }}
                    onSelectionChange={(key) => {
                        copy(key as string);
                        toast.success("下载链接已复制到剪切板");
                    }}
                    className="w-32"
                >
                {
                    files.map((file) => (
                        <AutocompleteItem key={file.path}>
                            {file.label}
                        </AutocompleteItem>
                    ))
                }
                </Autocomplete>
            )}
        </header>
    )
}
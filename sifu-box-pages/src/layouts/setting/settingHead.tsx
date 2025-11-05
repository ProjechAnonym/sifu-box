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
import { FileData } from "@/types/hosting/file";
export default function SettingHead(props: {
  token: string;
  admin: boolean;
  theme: string;
  files: Array<FileData>;
  
  setUpdate: (update: boolean) => void;

}) {
    const { token, admin, theme, files, setUpdate } = props;
    const file_input = useRef<HTMLInputElement>(null);
    const refresh = () => toast.promise(Refresh(token), {
            loading: "更新配置文件中...",
            success: (res) => {
                res !== true && res !== false ? res.message.map(item=>toast.error(item.message)) : toast.error("未知错误");
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
                return "";
            }
            return e.response.data;
            },
        })
    return (
        <header className="flex flex-wrap gap-2 p-2"> 
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
                {/* <Button
                    variant="shadow"
                    color="primary"
                    size="sm"
                    onPress={() => setMode(!mode)}>
                    <span className="font-black text-lg">{mode ? `模板` : `` }</span>
                </Button> */}
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
                        res !== true && res !== false ? res.map(item => !item.status && toast.error(item.message)) : toast.error("未知错误")
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
                    classNames={{ popoverContent: `${theme} bg-content1 text-foreground` }}
                    onSelectionChange={(key) => {
                        copy(key as string);
                        toast.success("下载链接已复制到剪切板");
                    }}
                    className="w-32"
                >
                    {
                        files.map((file) => (
                            <AutocompleteItem key={`${window.location.origin}/api/files/download/${file.name}/${file.expire_time}/${file.signature}/${file.path}`}>
                                {file.name}
                            </AutocompleteItem>
                        ))
                    }
                </Autocomplete>
            )}
        </header>
    )
}
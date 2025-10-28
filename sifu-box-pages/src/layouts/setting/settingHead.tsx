import { Button,ButtonGroup } from "@heroui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@heroui/popover";
import Interval from "@/components/card/interval";
import toast from "react-hot-toast";
import { Export } from "@/utils/migrate/export";
export default function SettingHead(props: {
  token: string;
  admin: boolean;
  theme: string;
  setUpdate: (update: boolean) => void;
}) {
    const { token, admin, theme } = props;
    return (
        <header className="flex flex-wrap gap-1 p-1"> 
            <ButtonGroup>
                <Button
                    variant="shadow"
                    color="primary"
                    size="sm"
                    
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
                    // onPress={() => fileInput.current && fileInput.current.click()}
                    >
                    <span className="font-black text-lg">导入</span>
                </Button>
            </ButtonGroup>
        </header>
    )
}
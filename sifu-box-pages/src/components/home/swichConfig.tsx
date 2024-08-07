import { Button, ScrollShadow } from "@nextui-org/react";
import toast from "react-hot-toast";
import { SwitchFile } from "@/utils/configs/SwitchConfig";

export default function SwitchConfig(props: {
  url: string;
  config: string;
  dark: boolean;
  secret: string;
  files: Array<{ label: string; path: string }> | null;
  onSwitch: (config: string) => void;
}) {
  const { url, config, secret, files, onSwitch } = props;
  return (
    <ScrollShadow className="h-12 w-44" orientation="horizontal">
      {files ? (
        <div className="grid grid-cols-repeat-2 gap-x-3 gap-y-2">
          {files
            ? files.map((file, i) => (
                <Button
                  key={`${file.label}-${i}`}
                  size="sm"
                  color={file.label === config ? "primary" : "default"}
                  className="gap-1 h-12"
                  onPress={() =>
                    toast.promise(SwitchFile(secret, url, file.label), {
                      loading: "loading",
                      success: (res) => {
                        res && onSwitch(file.label);
                        return `更新配置文件为${file.label}`;
                      },
                      error: (err) => `${err.response.data.message}`,
                    })
                  }
                >
                  <span className="text-md font-black p-0 text-wrap">
                    {file.label}
                  </span>
                  {file.label === config && (
                    <span className="p-0">
                      <i className="bi bi-check text-md" />
                    </span>
                  )}
                </Button>
              ))
            : "尚无配置文件"}
        </div>
      ) : (
        <div>无配置文件信息</div>
      )}
    </ScrollShadow>
  );
}

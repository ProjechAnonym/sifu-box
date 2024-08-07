import { useState, useEffect } from "react";
import { Input, Switch, Tooltip } from "@nextui-org/react";
import { ProviderValue } from "@/types/proxy";
export default function AddUrl(props: {
  onSubmit: (provider: ProviderValue) => void;
}) {
  const { onSubmit } = props;
  const [path, setPath] = useState<string>("");
  const [name, setName] = useState<string>("");
  const [proxy, setProxy] = useState<boolean>(false);
  const [remote, setRemote] = useState<boolean>(true);
  useEffect(() => {
    path !== "" &&
      name !== "" &&
      onSubmit({ path, name, proxy, remote, id: 0 });
  }, [path, proxy, name, remote]);
  return (
    <div className="flex flex-col gap-y-2">
      <Input
        type="url"
        label="path"
        className="w-full"
        size="sm"
        isRequired
        validationBehavior="native"
        value={path}
        onValueChange={(e) => setPath(e)}
      />
      <div className="w-full flex flex-row gap-2">
        <Input
          size="sm"
          label="label"
          isRequired
          validationBehavior="native"
          onValueChange={(e) => setName(e)}
          value={name}
        />
        <Tooltip content="是否使用代理下载该机场的链接">
          <Switch
            size="md"
            onValueChange={(e) => setProxy(e)}
            startContent={
              <span>
                <i className="bi bi-send-check-fill" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi-send-slash-fill" />
              </span>
            }
          />
        </Tooltip>
        <Tooltip content="是否远程链接,服务器模式一般为远程">
          <Switch
            size="md"
            onValueChange={(e) => setRemote(e)}
            defaultSelected
            startContent={
              <span>
                <i className="bi bi bi-cloud" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi bi-cloud-slash" />
              </span>
            }
          />
        </Tooltip>
      </div>
    </div>
  );
}

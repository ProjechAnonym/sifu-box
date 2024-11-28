import { useState, useEffect } from "react";
import { Input, Button } from "@nextui-org/react";
import toast from "react-hot-toast";
import { AddHost } from "@/utils/host/Addhost";
export default function HostAdd(props: {
  secret: string;
  submit: boolean;
  setUpdateHosts: (updateHosts: boolean) => void;
  setSubmit: (submit: boolean) => void;
  onClose: () => void;
}) {
  const { submit, setUpdateHosts, secret, setSubmit, onClose } = props;
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [url, setUrl] = useState("");
  const [port, setPort] = useState(0);
  useEffect(() => {
    submit &&
      toast.promise(
        AddHost(secret, {
          username,
          password,
          secret: token,
          url,
          port,
        }),
        {
          loading: "loading",
          success: (res) => {
            res && onClose();
            setUpdateHosts(true);
            setSubmit(false);
            return `添加${url}成功`;
          },
          error: (err) => {
            setSubmit(false);
            return err.code === "ERR_NETWORK"
              ? "网络错误"
              : `${err.response.data.message}`;
          },
        }
      );
  }, [submit]);
  return (
    <>
      <div className="flex flex-row gap-2">
        <Input
          size="sm"
          className="w-24"
          label="用户"
          isRequired
          validationBehavior="native"
          onValueChange={(e) => setUsername(e)}
          placeholder="root"
          value={username}
        />
        <Input
          size="sm"
          className="w-40"
          label="密码"
          type="password"
          isClearable
          isRequired
          validationBehavior="native"
          onValueChange={(e) => setPassword(e)}
          placeholder="password"
          value={password}
        />
      </div>
      <Input
        size="sm"
        label="链接"
        type="url"
        isClearable
        isRequired
        validationBehavior="native"
        onValueChange={(e) => setUrl(e)}
        placeholder="http://hostip"
        value={url}
      />
      <div className="flex flex-row gap-2 items-center">
        <Input
          size="sm"
          className="w-20"
          label="端口"
          isRequired
          validationBehavior="native"
          type="number"
          onValueChange={(e) => setPort(parseInt(e))}
          placeholder="9090"
          value={port.toString()}
        />
        <Input
          size="sm"
          className="w-24"
          label="token"
          isRequired
          onValueChange={(e) => setToken(e)}
          placeholder="123456"
          value={token}
        />
        <Button
          size="sm"
          color="danger"
          startContent={<i className="bi bi-x-lg text-xl" />}
          onPress={() => {
            setUsername("");
            setPassword("");
            setUrl("");
            setToken("");
            setPort(0);
          }}
          type="button"
        >
          <span className="font-black text-lg">清空</span>
        </Button>
      </div>
    </>
  );
}

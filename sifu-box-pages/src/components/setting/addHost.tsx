import { Input, Button } from "@nextui-org/react";
export default function HostAdd(props: {
  username: string;
  password: string;
  token: string;
  url: string;
  port: number;
  setUsername: (username: string) => void;
  setPassword: (password: string) => void;
  setUrl: (url: string) => void;
  setToken: (token: string) => void;
  setPort: (port: number) => void;
}) {
  const {
    username,
    password,
    token,
    url,
    port,
    setPassword,
    setPort,
    setToken,
    setUrl,
    setUsername,
  } = props;

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
        >
          <span className="font-black text-lg">清空</span>
        </Button>
      </div>
    </>
  );
}

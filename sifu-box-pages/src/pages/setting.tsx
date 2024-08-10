import { useEffect, useState } from "react";
import { NavigateFunction, useNavigate } from "react-router-dom";
import { useAppSelector, useAppDispatch } from "@/redux/hooks";
import {
  Modal,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalBody,
  useDisclosure,
  Button,
  Divider,
} from "@nextui-org/react";
import Load from "@/components/load";
import toast from "react-hot-toast";
import HostAdd from "@/components/setting/addHost";
import { ClienAuth } from "@/utils/ClientAuth";
import { AddHost } from "@/utils/host/Addhost";
function redirectLogin(navigate: NavigateFunction) {
  navigate("/login");
  toast.error("Please login", { duration: 2000 });
}
export default function Setting() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const secret = useAppSelector((state) => state.auth.secret);
  const status = useAppSelector((state) => state.auth.status);
  const login = useAppSelector((state) => state.auth.login);
  const load = useAppSelector((state) => state.auth.load);
  const dark = useAppSelector((state) => state.mode.dark);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [url, setUrl] = useState("");
  const [port, setPort] = useState(0);
  useEffect(() => {
    !login && !status && dispatch(ClienAuth({ auto: true }));
    login && !status && !load && redirectLogin(navigate);
  }, [status, login, load]);
  return (
    <div className="h-full">
      <Load show={load} fullscreen={true} />
      <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
        <ModalContent
          className={`${
            dark ? "sifudark" : "sifulight"
          } bg-background text-foreground`}
        >
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">
                添加主机
              </ModalHeader>
              <form
                onSubmit={(e) => {
                  e.preventDefault();
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
                        return `添加${url}成功`;
                      },
                      error: (err) => `${err.response.data.message}`,
                    }
                  );
                }}
              >
                <ModalBody>
                  <HostAdd
                    setPassword={setPassword}
                    setPort={setPort}
                    setToken={setToken}
                    setUrl={setUrl}
                    setUsername={setUsername}
                    url={url}
                    password={password}
                    token={token}
                    username={username}
                    port={port}
                  />
                </ModalBody>
                <ModalFooter className="w-full">
                  <Button color="danger" onPress={onClose} type="button">
                    <span className="text-lg font-black">关闭</span>
                  </Button>
                  <Button color="primary" type="submit">
                    <span className="text-lg font-black">提交</span>
                  </Button>
                </ModalFooter>
              </form>
            </>
          )}
        </ModalContent>
      </Modal>
      <header className="p-2 flex flex-wrap gap-x-2 gap-y-1 items-center">
        <Button onPress={onOpen} color="primary">
          <span className="font-black text-lg">添加主机</span>
        </Button>
      </header>
      <Divider />
    </div>
  );
}

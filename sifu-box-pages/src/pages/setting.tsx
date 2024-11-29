import { useEffect, useState, useRef } from "react";
import { NavigateFunction, useNavigate } from "react-router-dom";
import { useAppSelector, useAppDispatch } from "@/redux/hooks";
import { Button, Divider, Tooltip } from "@nextui-org/react";
import TemplateDash from "@/layout/templateDash";
import Load from "@/components/load";
import toast from "react-hot-toast";
import SetHost from "@/components/setting/setHost";
import { ClienAuth } from "@/utils/ClientAuth";
import { FetchHosts } from "@/utils/host/FetchHost";
import { FetchTemplate, RecoverTemplate } from "@/utils/template/FetchTemplate";
import { GetConfig, ImportConfig } from "@/utils/migrate/Migration";
import { HostValue } from "@/types/host";

function redirectLogin(navigate: NavigateFunction) {
  navigate("/login");
  toast.error("Please login", { duration: 2000 });
}
export default function Setting() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const recoverInput = useRef<HTMLInputElement>(null);
  const sifuboxInput = useRef<HTMLInputElement>(null);
  const headRef = useRef<HTMLHeadElement>(null);
  const [headHeight, setHeadHeight] = useState(0);
  const [updateTemplate, setUpdateTemplate] = useState(true);
  const [updateHosts, setUpdateHosts] = useState(true);
  const [hosts, setHosts] = useState<Array<HostValue> | null>(null);
  const [templates, setTemplates] = useState<Array<{
    Name: string;
    Template: Object;
  }> | null>(null);
  const [baseTemplate, setBaseTemplate] = useState<{
    Name: string;
    Template: Object;
  } | null>(null);
  const [updateBaseTemplate, setUpdateBaseTemplate] = useState(true);
  const secret = useAppSelector((state) => state.auth.secret);
  const status = useAppSelector((state) => state.auth.status);
  const login = useAppSelector((state) => state.auth.login);
  const load = useAppSelector((state) => state.auth.load);
  const dark = useAppSelector((state) => state.mode.dark);

  useEffect(() => {
    !login && !status && dispatch(ClienAuth({ auto: true }));
    login && !status && !load && redirectLogin(navigate);
    headRef.current && setHeadHeight(headRef.current.clientHeight);
    status &&
      updateHosts &&
      FetchHosts(secret)
        .then((res) => {
          res ? setHosts(res) : toast.error("Fetch servers failed");
          setUpdateHosts(false);
        })
        .catch((err) =>
          err.code === "ERR_NETWORK"
            ? toast.error("网络错误")
            : toast.error("Fetch servers failed")
        );
    status &&
      updateTemplate &&
      FetchTemplate(secret)
        .then((res) => {
          setTemplates(res);
          setUpdateTemplate(false);
        })
        .catch((e) => {
          setUpdateTemplate(false);
          e.code === "ERR_NETWORK"
            ? toast.error("网络错误")
            : toast.error(e.response.data.message);
        });
    status &&
      updateBaseTemplate &&
      RecoverTemplate(secret)
        .then((res) => {
          setBaseTemplate(res);
          setUpdateBaseTemplate(false);
        })
        .catch((e) => {
          setUpdateBaseTemplate(false);
          e.code === "ERR_NETWORK"
            ? toast.error("网络错误")
            : toast.error(e.response.data.message);
        });
  }, [
    status,
    login,
    load,
    headRef.current?.clientHeight,
    updateTemplate,
    updateHosts,
    updateBaseTemplate,
  ]);

  return (
    <div className="h-full">
      <Load show={load} fullscreen={true} />
      <header ref={headRef} className="p-2">
        <SetHost
          secret={secret}
          dark={dark}
          hosts={hosts}
          templates={templates}
          setUpdateHosts={setUpdateHosts}
        />
        <div className="flex flex-wrap gap-2 items-center my-2">
          <Button onPress={() => GetConfig(secret)} color="primary" size="sm">
            <span className="font-black text-lg">备份</span>
          </Button>
          <Button color="danger" size="sm">
            <label htmlFor="recoverFile-upload" className="font-black text-lg">
              恢复
            </label>
          </Button>
          <input
            type="file"
            id="recoverFile-upload"
            className="hidden"
            onChange={(e) =>
              e.target.files &&
              toast.promise(ImportConfig(secret, e.target.files), {
                loading: "loading",
                success: (res) => {
                  if (recoverInput.current) {
                    recoverInput.current.value = "";
                  }
                  return res ? "恢复成功" : "恢复失败";
                },
                error: (err) => {
                  if (recoverInput.current) {
                    recoverInput.current.value = "";
                  }
                  return err.code === "ERR_NETWORK"
                    ? "网络错误"
                    : `${err.response.data.message}`;
                },
              })
            }
            ref={recoverInput}
          />
          <Tooltip content="上传sifu-box新版本,务必确保文件名为sifu-box">
            <Button color="primary" size="sm">
              <label
                htmlFor="sifuboxFile-upload"
                className="font-black text-md"
              >
                上传sifu-box
              </label>
            </Button>
          </Tooltip>
          <input
            type="file"
            id="sifuboxFile-upload"
            className="hidden"
            ref={sifuboxInput}
          />
        </div>
      </header>
      <Divider />
      <TemplateDash
        baseTemplate={baseTemplate}
        secret={secret}
        dark={dark}
        headHeight={headHeight}
        template={templates}
        setUpdateTemplate={setUpdateTemplate}
      />
    </div>
  );
}

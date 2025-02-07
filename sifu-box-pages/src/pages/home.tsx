import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import DefaultLayout from "@/layouts/default";
import SingBox from "@/components/singbox";
import { HomeDashBoard } from "@/components/dashboard";
import Status from "@/components/status";
import { toast } from "react-hot-toast";
import { AutoLogin } from "@/utils/auth";
import { FetchSingboxApplication } from "@/utils/application";

export default function HomePage() {
  const admin = useAppSelector((state) => state.auth.admin);
  const auto = useAppSelector((state) => state.auth.auto);
  const status = useAppSelector((state) => state.auth.status);
  const theme = useAppSelector((state) => state.theme.theme);
  const navigate = useNavigate();
  const token = useAppSelector((state) => state.auth.jwt);
  const dispatch = useAppDispatch();
  const headerContainer = useRef<HTMLHeadElement>(null);
  const [height, setHeight] = useState(0);
  const [secret, setSecret] = useState("");
  const [listen, setListen] = useState("");
  const [log, setLog] = useState(false);
  const [provider, setProvider] = useState("");
  const [template, setTemplate] = useState("");
  useEffect(() => {
    !auto && !status && navigate("/");
    auto && dispatch(AutoLogin({}));
    token !== "" &&
      FetchSingboxApplication(token)
        .then((res) => {
          res.status && setListen(res.msg.listen);
          res.status && setSecret(res.msg.secret);
          res.status && setProvider(res.msg.current_provider);
          res.status && setTemplate(res.msg.current_template);
          res.status && setLog(res.msg.log);
        })
        .catch((e) => {
          if (e.code === "ERR_NETWORK") {
            toast.error("请检查网络连接");
            return;
          }
          if (e.response.data.message) {
            setListen(e.response.data.message.listen);
            setSecret(e.response.data.message.secret);
            setProvider(e.response.data.message.current_provider);
            setTemplate(e.response.data.message.current_template);
            setLog(e.response.data.message.log);
            toast.error(e.response.data.message.error);
            return;
          }
        });
    headerContainer.current && setHeight(headerContainer.current.clientHeight);
  }, [
    admin,
    token,
    headerContainer.current && headerContainer.current.clientHeight,
  ]);
  return (
    <DefaultLayout>
      <header className="flex flex-wrap gap-2" ref={headerContainer}>
        <HomeDashBoard
          provider={provider}
          template={template}
          admin={admin}
          token={token}
          theme={theme}
        />
        <Status listen={listen} secret={secret} log={log} theme={theme} />
      </header>
      <SingBox listen={listen} secret={secret} height={height} theme={theme} />
    </DefaultLayout>
  );
}

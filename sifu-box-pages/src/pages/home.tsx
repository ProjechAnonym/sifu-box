import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hook";
import DefaultLayout from "@/layouts/default";
import OutboundsTable from "@/layouts/table/outbound";
// import SingBox from "@/components/singbox";
// import { HomeDashBoard } from "@/components/dashboard";
// import Status from "@/components/status";
import { toast } from "react-hot-toast";
import { Verify } from "@/utils/auth";
import { FetchYacd } from "@/utils/configuration/fetch";

export default function HomePage() {
  const admin = useAppSelector((state) => state.auth.admin);
  const auto = useAppSelector((state) => state.auth.auto);
  const status = useAppSelector((state) => state.auth.status);
  const theme = useAppSelector((state) => state.theme.theme);
  const navigate = useNavigate();
  const token = useAppSelector((state) => state.auth.jwt);
  const dispatch = useAppDispatch();
  const header_container = useRef<HTMLHeadElement>(null);
  const [height, setHeight] = useState(0);
  const [secret, setSecret] = useState("");
  const [listen, setListen] = useState("");
  const [log, setLog] = useState(false);
  const [template, setTemplate] = useState("");
  useEffect(() => {
    !auto && !status && navigate("/");
    auto && dispatch(Verify({}));
    token !== "" &&
      FetchYacd(token)
        .then((res) => {
          res.status && setListen(res.msg.url);
          res.status && setSecret(res.msg.secret);
          res.status && setTemplate(res.msg.template ? res.msg.template : "");
          res.status && setLog(res.msg.log);
        })
        .catch((e) => {
          if (e.code === "ERR_NETWORK") {
            toast.error("请检查网络连接");
            return;
          }
          if (e.response.data.message) {
            setListen(e.response.data.message.listen ? e.response.data.message.listen : "");
            setSecret(e.response.data.message.secret ? e.response.data.message.secret : "");
            setTemplate(e.response.data.message.current_template ? e.response.data.message.current_template : "");
            setLog(e.response.data.message.log ? e.response.data.message.log : false);
            toast.error(e.response.data.message);
            return;
          }
        });
    header_container.current && setHeight(header_container.current.clientHeight);
  }, [
    admin,
    token,
    header_container.current && header_container.current.clientHeight,
  ]);
  return (
    <DefaultLayout>
      <header className="flex flex-wrap gap-2" ref={header_container}>
        {/* <HomeDashBoard
          provider={provider}
          template={template}
          admin={admin}
          token={token}
          theme={theme}
        />
        <Status listen={listen} secret={secret} log={log} theme={theme} /> */}
      </header>
      <OutboundsTable listen={listen} secret={secret} height={height} theme={theme} />
    </DefaultLayout>
  );
}
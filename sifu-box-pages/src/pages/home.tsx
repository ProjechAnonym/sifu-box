import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";

import DefaultLayout from "@/layouts/default";
import DashBoard from "@/components/dashboard";
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
  const [secret, setSecret] = useState("");
  const [listen, setListen] = useState("");
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
        })
        .catch((e) => toast.error(e));
  }, [admin, token]);
  return (
    <DefaultLayout>
      <header className="h-32 flex gap-2">
        {admin && (
          <DashBoard
            provider={provider}
            template={template}
            admin={admin}
            token={token}
            theme={theme}
          />
        )}
      </header>
    </DefaultLayout>
  );
}

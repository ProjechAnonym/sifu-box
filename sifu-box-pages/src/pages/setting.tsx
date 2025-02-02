import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import toast from "react-hot-toast";
import DefaultLayout from "@/layouts/default";
import Template from "@/components/template";
import {
  FetchConfiguration,
  FetchDefaultTemplate,
} from "@/utils/configuration";
import { AutoLogin } from "@/utils/auth";
import { provider, ruleset } from "@/types/configuration";
export default function SettingPages() {
  const admin = useAppSelector((state) => state.auth.admin);
  const auto = useAppSelector((state) => state.auth.auto);
  const status = useAppSelector((state) => state.auth.status);
  const theme = useAppSelector((state) => state.theme.theme);
  const navigate = useNavigate();
  const token = useAppSelector((state) => state.auth.jwt);
  const dispatch = useAppDispatch();
  const [providers, setProviders] = useState<Array<provider>>([]);
  const [ruleset, setRuleset] = useState<Array<ruleset>>([]);
  const [template, setTemplate] = useState<{ [key: string]: Object } | null>(
    null
  );
  const [defaultTemplate, setDefaultTemplate] = useState<Object | null>(null);

  const [update, setUpdate] = useState(true);
  useEffect(() => {
    !admin && navigate("/");
    !auto && !status && navigate("/");
    auto && dispatch(AutoLogin({}));
    token !== "" &&
      update &&
      FetchConfiguration(token)
        .then((res) => {
          setUpdate(false);
          setProviders(res.message.providers);
          setRuleset(res.message.rulesets);
          setTemplate(res.message.templates);
        })
        .catch((e) => {
          setUpdate(false);
          return e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data);
        });
    token !== "" &&
      !defaultTemplate &&
      FetchDefaultTemplate(token)
        .then((res) =>
          res.status
            ? setDefaultTemplate(res.message)
            : toast.error("获取默认模板失败")
        )
        .catch((e) =>
          e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data)
        );
  }, [admin, auto, status, token, update, defaultTemplate]);
  return (
    <DefaultLayout>
      <Template
        template={template}
        token={token}
        setUpdate={setUpdate}
        defaultTemplate={defaultTemplate}
        theme={theme}
      />
    </DefaultLayout>
  );
}

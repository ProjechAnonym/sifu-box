import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import toast from "react-hot-toast";
import DefaultLayout from "@/layouts/default";
import { SettingDashBoard } from "@/components/dashboard";
import Template from "@/components/template";
import Provider from "@/components/providers";
import Ruleset from "@/components/ruletset";
import {
  FetchConfiguration,
  FetchDefaultTemplate,
} from "@/utils/configuration";
import { FetchSingboxApplication } from "@/utils/application";
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
  const [currentProvider, setCurrentProvider] = useState("");
  const [currentTemplate, setCurrentTemplate] = useState("");
  const [updateCurrentApplication, setUpdateCurrentApplication] =
    useState(true);
  const [height, setHeight] = useState<number>(0);
  const [providers, setProviders] = useState<Array<provider>>([]);
  const [rulesets, setRulesets] = useState<Array<ruleset>>([]);
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
          setRulesets(res.message.rulesets);
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
    updateCurrentApplication &&
      token !== "" &&
      FetchSingboxApplication(token)
        .then((res) => {
          res.status && setCurrentProvider(res.msg.current_provider);
          res.status && setCurrentTemplate(res.msg.current_template);
        })
        .catch((e) => {
          if (e.code === "ERR_NETWORK") {
            toast.error("请检查网络连接");
            return;
          }
          if (e.response.data.message) {
            setCurrentProvider(e.response.data.message.current_provider);
            setCurrentTemplate(e.response.data.message.current_template);
            toast.error(e.response.data.message.error);
            return;
          }
        });
    updateCurrentApplication && setUpdateCurrentApplication(false);
  }, [
    admin,
    auto,
    status,
    token,
    update,
    defaultTemplate,
    updateCurrentApplication,
  ]);
  return (
    <DefaultLayout>
      <SettingDashBoard
        token={token}
        theme={theme}
        admin={admin}
        setUpdate={setUpdate}
      />
      <Template
        currentTemplate={currentTemplate}
        onHeight={setHeight}
        template={template}
        token={token}
        setUpdate={setUpdate}
        defaultTemplate={defaultTemplate}
        theme={theme}
        setUpdateCurrentSetting={setUpdateCurrentApplication}
      />
      <div style={{ height: `calc(100% - ${height}px - 3.5rem)` }}>
        <Provider
          setUpdateProviders={setUpdate}
          providers={providers}
          currentProvider={currentProvider}
          token={token}
          setCurrentUpdate={setUpdateCurrentApplication}
          theme={theme}
        />
        <Ruleset
          token={token}
          rulesets={rulesets}
          setUpdateRuleset={setUpdate}
          theme={theme}
        />
      </div>
    </DefaultLayout>
  );
}

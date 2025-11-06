import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hook";
import toast from "react-hot-toast";
import DefaultLayout from "@/layouts/default";
import SettingHead from "@/layouts/setting/settingHead";
import ProviderLayout from "@/layouts/setting/provider";
import { FetchFile } from "@/utils/hosting/fetch";
import { FetchConfiguration } from "@/utils/configuration/fetch";
import { Verify } from "@/utils/auth";
import { FileData } from "@/types/hosting/file";
import { Provider } from "@/types/setting/provider";
import { RuleSet } from "@/types/setting/ruleset";
export default function SettingPage() {
  const admin = useAppSelector((state) => state.auth.admin);
  const auto = useAppSelector((state) => state.auth.auto);
  const status = useAppSelector((state) => state.auth.status);
  const theme = useAppSelector((state) => state.theme.theme);
  const navigate = useNavigate();
  const token = useAppSelector((state) => state.auth.jwt);
  const dispatch = useAppDispatch();
  const [files, setFiles] = useState<Array<FileData>>([])
  const [providers, setProviders] = useState<Array<Provider>>([]);
  const [rulesets, setRulesets] = useState<Array<RuleSet>>([]);
  const [height, setHeight] = useState(0);
  const [template_mode, setTemplateMode] = useState(false);
  const [update, setUpdate] = useState(true);
  useEffect(() => {
    !admin && navigate("/");
    !auto && !status && navigate("/");
    auto && dispatch(Verify({}));
    token !== "" && update && FetchConfiguration(token)
      .then((res) => {
          res && res.map(item => {
            switch (item.type) {
              case "provider":
                item.message.every(
                  (provider: any): provider is Provider => 
                    typeof provider === "object" && "id" in provider && "name" in provider && "path" in provider && "remote" in provider && 
                    typeof provider.id === "number" && typeof provider.name === "string" && typeof provider.path === "string" && typeof provider.remote === "boolean") && 
                    setProviders(item.message)
                break;
              case "ruleset":
                item.message.every(
                  (ruleset: any): ruleset is RuleSet => 
                    typeof ruleset === "object" && "id" in ruleset && "name" in ruleset && "path" in ruleset && "remote" in ruleset && "binary" in ruleset &&
                    typeof ruleset.id === "number" && typeof ruleset.name === "string" && typeof ruleset.path === "string" && typeof ruleset.remote === "boolean" && typeof ruleset.binary === "boolean") &&  
                    setRulesets(item.message)
                break
              default:
                break;
            }
            
          })
          setUpdate(false);
      }).catch((e) => {
          setUpdate(false);
          return e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data);
      });
    token !== "" && update && FetchFile(token)
      .then((res) => 
        res ? setFiles(res.message.map(item=>item)) : toast.error("获取模板文件失败")).catch((e) => {
          setUpdate(false);
          return e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data);
      });
    }, [admin, auto, status, token, update]);
  return (
    <DefaultLayout>
      <SettingHead
        template_mode={template_mode}
        setTemplateMode={setTemplateMode}
        token={token}
        admin={admin}
        theme={theme}
        files={files}
        setUpdate={setUpdate}
        setHeight={setHeight}
      />
      <div style={{height: `calc(100% - ${height}px)`}}>
        {template_mode ? <>1</>: <ProviderLayout providers={providers} />}
      </div>
      
    </DefaultLayout>
  );
}
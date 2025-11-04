import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hook";
import toast from "react-hot-toast";
import DefaultLayout from "@/layouts/default";
import SettingHead from "@/layouts/setting/settingHead";
import { FetchFile } from "@/utils/hosting/fetch";
import { FetchConfiguration } from "@/utils/configuration/fetch";
import { Verify } from "@/utils/auth";
import { FileData } from "@/types/hosting/file";
export default function SettingPage() {
  const admin = useAppSelector((state) => state.auth.admin);
  const auto = useAppSelector((state) => state.auth.auto);
  const status = useAppSelector((state) => state.auth.status);
  const theme = useAppSelector((state) => state.theme.theme);
  const navigate = useNavigate();
  const token = useAppSelector((state) => state.auth.jwt);
  const dispatch = useAppDispatch();
  const [current_template, setCurrentTemplate] = useState("");
  const [files, setFiles] = useState<Array<FileData>>([])
  const [updateCurrentApplication, setUpdateCurrentApplication] =
    useState(true);
  
  const [update, setUpdate] = useState(true);
  useEffect(() => {
    !admin && navigate("/");
    !auto && !status && navigate("/");
    auto && dispatch(Verify({}));
    token !== "" && update &&
      FetchConfiguration(token)
        .then((res) => {
          setUpdate(false);
        })
        .catch((e) => {
          setUpdate(false);
          return e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data);
        });
    token !== "" && update && FetchFile(token)
      .then((res) => res ? setFiles(res.message.map(item=>item)) : toast.error("获取模板文件失败")).catch((e) => {
          setUpdate(false);
          return e.code === "ERR_NETWORK"
            ? toast.error("请检查网络连接")
            : e.response.data.message
              ? toast.error(e.response.data.message)
              : toast.error(e.response.data);
      });
    }, [admin, auto, status, token, update, updateCurrentApplication]);
  return (
    <DefaultLayout>
      <SettingHead
        token={token}
        admin={admin}
        theme={theme}
        setUpdate={setUpdate}
        files={files}
      />
      
    </DefaultLayout>
  );
}
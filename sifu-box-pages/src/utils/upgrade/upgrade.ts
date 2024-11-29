import axios from "axios";
import { SharedSelection } from "@nextui-org/react";
export async function UpgradeApp(
  secret: string,
  file: FileList,
  hosts: SharedSelection
) {
  const application = file[0];
  if (application.name != "sing-box") {
    console.error("文件名必须为sing-box");
    throw { response: { data: { message: ["文件名必须为sing-box"] } } };
  }
  const formData = new FormData();
  Array.from(hosts).forEach((value) =>
    formData.append("addr", value.toString().split("-")[0])
  );
  formData.append("file", application);
  try {
    const res = await axios.post("/api/upgrade/singbox", formData, {
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

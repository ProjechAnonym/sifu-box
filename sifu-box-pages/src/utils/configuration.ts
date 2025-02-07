import axios from "axios";
import { provider, ruleset } from "@/types/configuration";
export async function FetchConfiguration(token: string) {
  try {
    const res = await axios.get("/api/configuration/fetch", {
      headers: { Authorization: token },
    });
    return res.status === 200
      ? { status: true, message: res.data.message }
      : { status: false, message: "获取配置失败" };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function DeleteItems(
  token: string,
  value: Array<string>,
  category: string
) {
  const data = new FormData();
  switch (category) {
    case "providers":
      value.forEach((item) => {
        data.append("providers", item);
      });
      break;
    case "rulesets":
      value.forEach((item) => {
        data.append("rulesets", item);
      });
      break;
    case "templates":
      value.forEach((item) => {
        data.append("templates", item);
      });
      break;
    default:
      throw { response: { data: { message: `不存在${value}键值` } } };
  }
  try {
    const res = await axios.delete(
      "/api/configuration/items",

      { data: data, headers: { Authorization: token } }
    );
    return res.status === 200
      ? { status: true, message: res.data.message }
      : { status: false, message: "删除配置失败" };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function FetchDefaultTemplate(token: string) {
  try {
    const res = await axios.get("/api/configuration/recover", {
      headers: { Authorization: token },
    });
    return res.status === 200
      ? { status: true, message: res.data.message }
      : { status: false, message: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function ModifyTemplate(
  token: string,
  name: string,
  string: string
) {
  const content = JSON.parse(string);
  try {
    const res = await axios.post("/api/configuration/template", content, {
      headers: { Authorization: token },
      params: { name: name },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function AddProviderFiles(token: string, files: FileList) {
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    files.item(i) &&
      formData.append("files", files.item(i)!, files.item(i)!.name);
  }

  try {
    const res = await axios.post("/api/configuration/files", formData, {
      headers: { Authorization: token },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function AddItems(
  token: string,
  items: Array<provider | ruleset>,
  kind: string
) {
  let data;
  switch (kind) {
    case "providers":
      data = { providers: items };
      break;
    case "rulesets":
      data = {
        rulesets: items.map((item) => {
          (item as ruleset).update_interval += "d";
          return item;
        }),
      };
      break;
    default:
      throw { response: { data: { message: [`没有"${kind}"字段`] } } };
  }
  try {
    const res = await axios.post("/api/configuration/add", data, {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

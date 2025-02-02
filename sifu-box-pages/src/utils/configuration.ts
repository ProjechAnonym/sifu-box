import axios from "axios";
export async function FetchConfiguration(token: string) {
  try {
    const res = await axios.get(
      "http://192.168.1.2:8080/api/configuration/fetch",
      { headers: { Authorization: token } }
    );
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
      "http://192.168.1.2:8080/api/configuration/items",

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
    const res = await axios.get(
      "http://192.168.1.2:8080/api/configuration/recover",
      { headers: { Authorization: token } }
    );
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
    const res = await axios.post(
      "http://192.168.1.2:8080/api/configuration/template",
      content,
      {
        headers: { Authorization: token },
        params: { name: name },
      }
    );
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

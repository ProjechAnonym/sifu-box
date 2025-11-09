import axios from "axios";

export async function SetTemplate(token: string, name: string) {
  const data = new FormData();
  data.append("name", name);
  try {
    const res = await axios.post("http://192.168.10.6:9090/api/application/template", data, {
      headers: { Authorization: token },
    });
    return res.status === 200 ? `设置"${name}"模板成功` : "遭遇未知错误";
  } catch (e) {
    console.error(e);
    throw e;
  }
}
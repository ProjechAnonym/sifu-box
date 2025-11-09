import axios from "axios";

export async function DeleteTemplate(token: string, value: Array<string>) {
  const data = new FormData();
  value.forEach((item) => data.append("name", item));
  try {
    const res = await axios.delete("http://192.168.10.6:9090/api/configuration/delete/template",
        { data: data, headers: { Authorization: token } }
    );
    return res.status === 207
      ? (IsValidRes(res.data) ? res.data : false)
      : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function AddTemplateMsg(token: string, template: {name: string, [key: string]: any}) {
  try {
    const res = await axios.post("http://192.168.10.6:9090/api/configuration/add/template",
        template,
        { headers: { Authorization: token } }
    );
    return res.status === 200 ? res.data.message : "遭遇未知错误"
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function EditTemplateMsg(token: string, template: {name: string, [key: string]: any}) {
  try {
    const res = await axios.patch("http://192.168.10.6:9090/api/configuration/edit/template",
        template,
        { headers: { Authorization: token } }
    );
    return res.status === 200 ? res.data.message : "遭遇未知错误"
  } catch (e) {
    console.error(e);
    throw e;
  }
}
function IsValidRes(res: any): res is {status: boolean, message: string}[] {
    return Array.isArray(res) && res.every((item: unknown): item is {status: boolean, message: string} => typeof item === 'object' && item !== null && 'status' in item && 'message' in item && typeof item.status === 'boolean' && typeof item.message === 'string')
}
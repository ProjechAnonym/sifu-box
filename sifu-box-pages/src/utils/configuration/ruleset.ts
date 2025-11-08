import axios from "axios";
export async function DeleteRuleset(token: string, value: Array<string>) {
  const data = new FormData();
  value.forEach((item) => data.append("name", item));
  try {
    const res = await axios.delete("http://192.168.10.6:9090/api/configuration/delete/ruleset",
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
export async function AddRulesetFiles(token: string, files: FileList) {
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    files.item(i) && formData.append("file", files.item(i)!, files.item(i)!.name);
  }
  try {
    const res = await axios.post("http://192.168.10.6:9090/api/configuration/add/ruleset/local",
        formData,
        { headers: { Authorization: token } }
    );
    return res.status === 207
      ? IsValidRes(res.data) ? res.data : false
      : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function AddRulesetMsg(token: string, providers: Array<{name: string, path: string, remote: boolean}>) {
  try {
    const res = await axios.post("http://192.168.10.6:9090/api/configuration/add/ruleset/remote",
        providers,
        { headers: { Authorization: token } }
    );
    return res.status === 207
      ? IsValidRes(res.data) ? res.data : false
      : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function EditRuleset(token: string, provider: {name: string, path: string, remote: boolean}) {
  const data = new FormData();
  data.append("name", provider.name);
  data.append("path", provider.path);
  data.append("remote", provider.remote.toString());
  try {
    const res = await axios.patch("http://192.168.10.6:9090/api/configuration/edit/ruleset",
        data,
        { headers: { Authorization: token } }
    );
    return res.status === 200 ? res.data : false 
  } catch (e) {
    console.error(e);
    throw e;
  }
}
function IsValidRes(res: any): res is {status: boolean, message: string}[] {
    return Array.isArray(res) && res.every((item: unknown): item is {status: boolean, message: string} => typeof item === 'object' && item !== null && 'status' in item && 'message' in item && typeof item.status === 'boolean' && typeof item.message === 'string')
}
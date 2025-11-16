import axios from "axios";
export async function DeleteProvider(token: string, value: Array<string>) {
  const data = new FormData();
  value.forEach((item) => data.append("name", item));
  try {
    const res = await axios.delete("/api/configuration/delete/provider",
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
export async function AddProviderFiles(token: string, files: FileList) {
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    files.item(i) && formData.append("file", files.item(i)!, files.item(i)!.name);
  }
  try {
    const res = await axios.post("/api/configuration/add/provider/local",
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

export async function AddProviderMsg(token: string, providers: Array<{name: string, path: string, remote: boolean}>) {
  try {
    const res = await axios.post("/api/configuration/add/provider/remote",
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

export async function EditProvider(token: string, provider: {name: string, path: string, remote: boolean}) {
  const data = new FormData();
  data.append("name", provider.name);
  data.append("path", provider.path);
  data.append("remote", provider.remote.toString());
  try {
    const res = await axios.patch("/api/configuration/edit/provider",
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
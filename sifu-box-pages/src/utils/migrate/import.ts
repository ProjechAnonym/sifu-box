import axios from "axios";
export async function Import(token: string, file: File) {
  const data = new FormData();
  data.append("file", file);

  try {
    const res = await axios.post("http://192.168.10.6:9090/api/migrate/import", data, {
      headers: { Authorization: token },
    });
    return res.status === 200 ? true : res.data;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
import axios from "axios";

export async function FetchFile(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/files/list", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200 ? res.data : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
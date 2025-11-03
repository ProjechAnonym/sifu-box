import axios from "axios";
export async function Refresh(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/refresh", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200 ? true : res.data;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
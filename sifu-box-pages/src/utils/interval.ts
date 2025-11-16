import axios from "axios";
export async function SetInterval(token: string, interval: string) {
  const data = new FormData();
  data.append("interval", interval);
  try {
    const res = await axios.post("/api/application/interval", data, {
      headers: { Authorization: token },
    });
    return res.status === 200 && res.data.message;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
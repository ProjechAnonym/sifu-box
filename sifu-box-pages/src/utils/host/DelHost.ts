import axios from "axios";
export async function DeleteHost(secret: string, addr: string) {
  const data = new FormData();
  data.append("url", addr);
  try {
    const res = await axios.delete("/api/host/delete", {
      data: data,
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

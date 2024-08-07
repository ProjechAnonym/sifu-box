import axios from "axios";
export async function RefreshConfig(secret: string) {
  try {
    const res = await axios.get("/api/exec/refresh", {
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

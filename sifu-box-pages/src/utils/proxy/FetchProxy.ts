import axios from "axios";
export async function FetchProxy(secret: string) {
  try {
    const res = await axios.get("/api/proxy/fetch", {
      headers: {
        Authorization: secret,
      },
    });
    return res.status === 200 ? res.data : null;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

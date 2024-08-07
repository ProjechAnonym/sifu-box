import axios from "axios";

export async function FetchProxyServers(url: string, secret: string) {
  try {
    const res = await axios.get(url + "/proxies", {
      headers: {
        Authorization: `Bearer ${secret}`,
      },
    });
    return res.status ? res.data : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

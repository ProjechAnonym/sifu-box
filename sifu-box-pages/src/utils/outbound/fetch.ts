import axios from "axios";
export async function FetchOutbounds(url: string, secret: string) {
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
export async function FetchRules(url: string, secret: string) {
  try {
    const res = await axios.get(url + "/rules", {
      headers: {
        Authorization: `Bearer ${secret}`,
      },
    });
    const values = res.status
      ? res.data.rules.map((item: any, i: number) => {
          return { key: `${item.payload}-${i}`, ...item };
        })
      : [];
    return values;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
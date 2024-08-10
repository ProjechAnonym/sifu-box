import axios from "axios";
export async function FetchRules(url: string, secret: string) {
  try {
    const res = await axios.get(url + "/rules", {
      headers: {
        Authorization: `Bearer ${secret}`,
      },
    });
    const labels = [
      { label: "类型", key: "type", allowSort: true, initShow: true },
      { label: "规则", key: "payload", allowSort: true, initShow: true },
      { label: "出站", key: "proxy", allowSort: true, initShow: true },
    ];
    const values = res.status
      ? res.data.rules.map((item: any, i: number) => {
          return { key: `${item.payload}-${i}`, ...item };
        })
      : [];
    return { labels, values };
  } catch (e) {
    console.error(e);
    throw e;
  }
}

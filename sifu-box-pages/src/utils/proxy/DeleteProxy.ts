import axios from "axios";

export async function DeleteItems(
  secret: string,
  urls: number[],
  rulesets: number[]
) {
  try {
    const res = await axios.delete("/api/proxy/delete", {
      data: { providers: urls, rulesets: rulesets },
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export function Key2Index(urls: string[], rulesets: string[]) {
  const urlsIndex =
    urls.length !== 0
      ? urls.map((url) => {
          const urlSlice = url.split("-");
          return parseInt(urlSlice[urlSlice.length - 1]);
        })
      : [];
  const rulesetsIndex =
    rulesets.length !== 0
      ? rulesets.map((ruleset) => {
          const rulesetSlice = ruleset.split("-");
          return parseInt(rulesetSlice[rulesetSlice.length - 1]);
        })
      : [];
  return { urlsIndex, rulesetsIndex };
}

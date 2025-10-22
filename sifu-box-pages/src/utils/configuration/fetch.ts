import axios from "axios";
export async function FetchYacd(token: string) {
    try {
    const res = await axios.get("http://192.168.10.6:9090/api/configuration/yacd", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
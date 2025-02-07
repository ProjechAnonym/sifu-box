import axios from "axios";
export async function ClearFakeip(url: string, secret: string) {
  try {
    const res = await axios.post(url + "/cache/fakeip/flush", null, {
      headers: { Authorization: `Bearer ${secret}` },
    });
    return res.data;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

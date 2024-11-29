import axios from "axios";

export async function FetchTemplate(secret: string) {
  try {
    const res = await axios.get(
      "http://192.168.213.128:8080/api/templates/fetch",
      {
        headers: { Authorization: secret },
      }
    );
    return res.status === 200 ? res.data.message : null;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function RecoverTemplate(secret: string) {
  try {
    const res = await axios.get(
      "http://192.168.213.128:8080/api/templates/recover",
      {
        headers: { Authorization: secret },
      }
    );
    return res.status === 200
      ? { Name: "基本模板", Template: res.data.message }
      : null;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

import axios from "axios";
export async function FetchSingboxApplication(token: string) {
  try {
    const res = await axios.get(
      "http://192.168.1.2:8080/api/application/fetch",
      {
        headers: {
          Authorization: token,
        },
      }
    );
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function SetTemplate(token: string, name: string) {
  const data = new FormData();
  data.append("value", name);

  try {
    const res = await axios.post(
      "http://192.168.1.2:8080/api/application/set/template",
      data,
      { headers: { Authorization: token } }
    );
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}

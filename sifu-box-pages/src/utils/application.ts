import axios from "axios";
export async function FetchSingboxApplication(token: string) {
  try {
    const res = await axios.get("/api/application/fetch", {
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

export async function SetTemplate(token: string, name: string) {
  const data = new FormData();
  data.append("value", name);

  try {
    const res = await axios.post("/api/application/set/template", data, {
      headers: { Authorization: token },
    });
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function SetProvider(token: string, name: string) {
  const data = new FormData();
  data.append("value", name);

  try {
    const res = await axios.post("/api/application/set/provider", data, {
      headers: { Authorization: token },
    });
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function SetInterval(token: string, interval: string) {
  const data = new FormData();
  data.append("interval", interval);
  try {
    const res = await axios.post("/api/application/interval", data, {
      headers: { Authorization: token },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

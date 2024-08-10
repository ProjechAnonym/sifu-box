import axios from "axios";
export async function GetServiceStatus(
  secret: string,
  service: string,
  url: string
) {
  const data = new FormData();
  data.append("service", service);
  data.append("url", url);
  try {
    const res = await axios.post("/api/exec/check", data, {
      headers: {
        Authorization: secret,
      },
    });
    return res.status === 200 && res.data.message;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function BootServiceUp(
  secret: string,
  service: string,
  url: string
) {
  const data = new FormData();
  data.append("service", service);
  data.append("url", url);
  try {
    const res = await axios.post("/api/exec/boot", data, {
      headers: {
        Authorization: secret,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function StopService(
  secret: string,
  service: string,
  url: string
) {
  const data = new FormData();
  data.append("service", service);
  data.append("url", url);
  try {
    const res = await axios.post("/api/exec/stop", data, {
      headers: {
        Authorization: secret,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

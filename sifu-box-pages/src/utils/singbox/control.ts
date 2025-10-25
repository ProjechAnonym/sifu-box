import axios from "axios";

export async function ControlSignal(token: string, action: string) {
try {
    const res = await axios.get(`http://192.168.10.6:9090/api/execute/${action}`, {
      headers: {
        Authorization: token,
      },
    });
    console.log(res.data.message);
    return res.status === 200 ? {message: res.data.message, status: res.data.message} : {message: res.data.message, status: false};
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function Check(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/check", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200 ? res.data.message : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function Reload(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/reload", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}


export async function Boot(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/boot", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function Stop(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/stop", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

import axios from "axios";
export async function Check(token: string) {
  try {
    const res = await axios.get("/api/exec/status", {
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
    const res = await axios.get("/api/exec/reload", {
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

export async function Restart(token: string) {
  try {
    const res = await axios.get("/api/exec/restart", {
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
    const res = await axios.get("/api/exec/boot", {
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
    const res = await axios.get("/api/exec/stop", {
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
export async function Refresh(token: string) {
  try {
    const res = await axios.get("/api/exec/refresh", {
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

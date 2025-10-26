import axios from "axios";

export async function ControlSignal(token: string, action: string) {
try {
    const res = await axios.get(`http://192.168.10.6:9090/api/execute/${action}`, {
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
import axios from "axios";
import { HostSendData } from "@/types/host";
export async function AddHost(secret: string, host: HostSendData) {
  try {
    const res = await axios.post("/api/host/add", host, {
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

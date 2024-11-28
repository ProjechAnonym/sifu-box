import axios from "axios";
import { HostRevData } from "@/types/host";
export async function FetchHosts(secret: string) {
  try {
    const response = await axios.get(
      "http://192.168.213.128:8080/api/host/fetch",
      {
        headers: { Authorization: secret },
      }
    );
    return response.status === 200
      ? response.data.map((server: HostRevData, i: number) => {
          return {
            key: `${server.url}-${i + 1}`,
            url: server.url,
            localhost: server.localhost,
            config: server.config,
            secret: server.secret,
            port: server.port,
          };
        })
      : null;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

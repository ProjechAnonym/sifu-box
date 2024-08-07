import axios from "axios";
import { cloneDeep } from "lodash";
import { ServerGroupValue, ServerValue } from "@/types/servers";
export async function SwitchServer(
  group: string,
  name: string,
  url: string,
  secret: string,
  servers: {
    [key: string]: ServerGroupValue | ServerValue;
  }
) {
  try {
    const res = await axios.put(
      `${url}/proxies/${encodeURI(group)}`,
      { name: name },
      { headers: { Authorization: `Bearer ${secret}` } }
    );
    if (res.status === 204) {
      const newServers = cloneDeep(servers);
      (newServers[group] as ServerGroupValue).now = name;
      return newServers;
    }
    return servers;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

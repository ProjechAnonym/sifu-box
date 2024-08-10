import axios from "axios";
import { cloneDeep } from "lodash";
import { ServerGroupValue, ServerValue } from "@/types/singbox";
export async function TestGroup(
  group: string,
  url: string,
  secret: string,
  servers: {
    [key: string]: ServerGroupValue | ServerValue;
  }
) {
  try {
    const res = await axios.get(`${url}/group/${encodeURI(group)}/delay`, {
      params: {
        timeout: 5000,
        url: "https://www.gstatic.com?generate_204",
      },
      headers: { Authorization: `Bearer ${secret}` },
    });
    const newServers = cloneDeep(servers);
    res.status === 200
      ? Object.keys(res.data).forEach((server) => {
          newServers[server].history = [
            { delay: res.data[server], time: new Date().toLocaleString() },
          ];
        })
      : (newServers[group].history = [
          { delay: 0, time: new Date().toLocaleString() },
        ]);
    return newServers;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function TestHost(
  name: string,
  url: string,
  secret: string,
  servers: {
    [key: string]: ServerGroupValue | ServerValue;
  }
) {
  try {
    const res = await axios.get(`${url}/proxies/${encodeURI(name)}/delay`, {
      params: {
        timeout: 5000,
        url: "https://www.gstatic.com?generate_204",
      },
      headers: { Authorization: `Bearer ${secret}` },
    });
    const newServers = cloneDeep(servers);
    res.status === 200
      ? (newServers[name].history = [
          { delay: res.data.delay, time: new Date().toLocaleString() },
        ])
      : (newServers[name].history = [
          { delay: 0, time: new Date().toLocaleString() },
        ]);
    return newServers;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

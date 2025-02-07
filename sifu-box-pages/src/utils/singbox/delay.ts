import axios from "axios";
import { cloneDeep } from "lodash";
import { Outbound, OutboundGroup } from "@/types/singbox/outbound";
export async function GroupDelay(
  group: string,
  url: string,
  secret: string,
  outbounds: {
    [key: string]: OutboundGroup | Outbound;
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
    const newOutbounds = cloneDeep(outbounds);
    res.status === 200
      ? Object.keys(res.data).forEach((outbound) => {
          newOutbounds[outbound].history = [
            { delay: res.data[outbound], time: new Date().toLocaleString() },
          ];
        })
      : (newOutbounds[group].history = [
          { delay: 0, time: new Date().toLocaleString() },
        ]);
    return newOutbounds;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function OutboundDelay(
  name: string,
  url: string,
  secret: string,
  outbounds: {
    [key: string]: Outbound | OutboundGroup;
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
    const newOutbounds = cloneDeep(outbounds);
    res.status === 200
      ? (newOutbounds[name].history = [
          { delay: res.data.delay, time: new Date().toLocaleString() },
        ])
      : (newOutbounds[name].history = [
          { delay: 0, time: new Date().toLocaleString() },
        ]);
    return newOutbounds;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

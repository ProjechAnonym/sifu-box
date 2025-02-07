import axios from "axios";
import { cloneDeep } from "lodash";
import { Outbound, OutboundGroup } from "@/types/singbox/outbound";
export async function SwitchOutbound(
  group: string,
  name: string,
  url: string,
  secret: string,
  outbounds: {
    [key: string]: OutboundGroup | Outbound;
  }
) {
  try {
    const res = await axios.put(
      `${url}/proxies/${encodeURI(group)}`,
      { name: name },
      { headers: { Authorization: `Bearer ${secret}` } }
    );
    if (res.status === 204) {
      const newOutbounds = cloneDeep(outbounds);
      (newOutbounds[group] as OutboundGroup).now = name;
      return newOutbounds;
    }
    return outbounds;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

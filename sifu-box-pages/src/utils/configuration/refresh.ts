import axios from "axios";
export async function Refresh(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/execute/refresh", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200 ? true : isRefreshRes(res.data) ? res.data : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

function isRefreshRes(res: any): res is { message: {message: string}[]} {
  return (
    typeof res === 'object' && res !== null && 'message' in res && Array.isArray(res.message) && 
      res.message.every((item: unknown): item is {message: string} => typeof item === 'object' && item !== null && 'message' in item && typeof item.message === 'string')
  );
}

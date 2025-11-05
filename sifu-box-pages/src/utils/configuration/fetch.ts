import axios from "axios";

export async function FetchYacd(token: string) {
    try {
    const res = await axios.get("http://192.168.10.6:9090/api/configuration/yacd", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200
      ? { status: true, msg: res.data.message }
      : { status: false, msg: res.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function FetchConfiguration(token: string) {
  try {
    const res = await axios.get("http://192.168.10.6:9090/api/configuration/fetch", {
      headers: { Authorization: token },
    });
    return res.status === 207 && isConfiguration(res.data)
      ? res.data.message
      : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

function isConfiguration(obj: any): obj is {message: Array<{message: Array<any>,status: boolean, type:string}>} {
  return (
    typeof obj === "object" && "message" in obj && Array.isArray(obj.message) &&
    obj.message.every((item: any): item is {message: Array<any>,status: boolean, type:string} => typeof item === "object" 
    && "message" in item && "status" in item && "type" in item && Array.isArray(item.message) && typeof item.status === "boolean" && typeof item.type === "string") 
  );
}
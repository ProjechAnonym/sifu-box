import axios from "axios";
import { FileData } from "@/types/hosting/file";
export async function FetchFile(token: string) {
  try {
    const res = await axios.get("/api/files/list", {
      headers: {
        Authorization: token,
      },
    });
    return res.status === 200 ? (isFileRes(res.data) && res.data) : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
function isFileRes(res: any): res is { message: FileData[] } {
  return (typeof res === "object" && res !== null && "message" in res && Array.isArray(res.message) && res.message.every(isFileData));
}
function isFileData(item: any): item is FileData {
  return typeof item === 'object' 
    && item !== null 
    && typeof item.name === 'string'
    && typeof item.expire_time === 'string' 
    && typeof item.signature === 'string'
    && typeof item.path === 'string';
}
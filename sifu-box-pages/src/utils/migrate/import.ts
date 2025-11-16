import axios from "axios";
export async function Import(token: string, file: File) {
  const data = new FormData();
  data.append("file", file);

  try {
    const res = await axios.post("/api/migrate/import", data, {
      headers: { Authorization: token },
    });
    return res.status === 200 ? true : isImportRes(res.data) ? res.data : false;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

function isImportRes(res: any): res is Array<{message: string, status: boolean}> {
  return typeof res === 'object' && res !== null && Array.isArray(res) && 
      res.every((item: unknown): item is {message: string, status: boolean} => typeof item === 'object' && item !== null && 'message' in item && 'status' in item && typeof item.message === 'string' && typeof item.status === 'boolean')
}
import axios from "axios";
import { saveAs } from "file-saver";
export async function Export(token: string) {
  try {
    const res = await axios.get("/api/migrate/export", {
      headers: {
        Authorization: token,
      },
      responseType: "blob",
    });
    console.log((res.headers["content-disposition"] as string).match(
          /filename="?(.+)"?/
        )![1]);
    res.status === 200 &&
      typeof res.headers["content-disposition"] === "string" &&
      (res.headers["content-disposition"] as string).match(
        /filename="?([^"]+)"?/
      ) &&
      (res.headers["content-disposition"] as string).match(/filename="?([^"]+)"?/)!
        .length > 0 &&
      saveAs(
        res.data,
        (res.headers["content-disposition"] as string).match(
          /filename="?([^"]+)"?/
        )![1]
      );
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
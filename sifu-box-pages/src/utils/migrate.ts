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
    res.status === 200 &&
      typeof res.headers["content-disposition"] === "string" &&
      (res.headers["content-disposition"] as string).match(
        /filename="?(.+)"?/
      ) &&
      (res.headers["content-disposition"] as string).match(/filename="?(.+)"?/)!
        .length > 0 &&
      saveAs(
        res.data,
        (res.headers["content-disposition"] as string).match(
          /filename="?(.+)"?/
        )![1]
      );
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

export async function Migrate(token: string, file: File) {
  const data = new FormData();
  data.append("file", file);

  try {
    const res = await axios.post("/api/migrate/import", data, {
      headers: { Authorization: token },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

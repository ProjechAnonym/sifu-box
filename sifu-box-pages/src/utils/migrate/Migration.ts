import axios from "axios";
import fileDownload from "js-file-download";
export async function GetConfig(secret: string) {
  try {
    const res = await axios.get(
      "http://192.168.213.128:8080/api/migrate/export",
      {
        headers: { Authorization: secret },
        responseType: "blob",
      }
    );
    fileDownload(res.data, "sifu-box-export.yaml");
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function ImportConfig(secret: string, files: FileList) {
  const formData = new FormData();
  formData.append("file", files.item(0)!, files.item(0)!.name);
  try {
    const res = await axios.post(
      "http://192.168.213.128:8080/api/migrate/import",
      formData,
      {
        headers: { Authorization: secret },
      }
    );
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

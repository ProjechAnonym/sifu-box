import axios from "axios";
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
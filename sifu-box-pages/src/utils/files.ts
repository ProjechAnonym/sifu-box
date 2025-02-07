import axios from "axios";
export async function FetchFiles(token: string) {
  try {
    const response = await axios.get("/api/files/fetch", {
      headers: {
        Authorization: token,
      },
    });
    return response.status !== 200
      ? { status: false, message: "Failed to fetch files" }
      : { status: true, message: response.data.message };
  } catch (e) {
    console.error(e);
    throw e;
  }
}

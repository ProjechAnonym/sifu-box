import axios from "axios";
export async function DeleteTemplate(secret: string, templates: Array<string>) {
  const formdata = new FormData();
  templates.forEach((value) => {
    formdata.append("names", value);
  });
  try {
    const res = await axios.delete("/api/templates/delete", {
      headers: { Authorization: secret },
      data: formdata,
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

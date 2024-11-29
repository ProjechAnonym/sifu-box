import axios from "axios";
import { stringify } from "yaml";
export async function ModifyTemplate(
  secret: string,
  name: string,
  jsonString: string
) {
  const content = JSON.parse(jsonString);
  const contentYaml = stringify(content);
  try {
    const res = await axios.post("/api/templates/set", contentYaml, {
      headers: { Authorization: secret },
      params: { name: name },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

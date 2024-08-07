import axios from "axios";
export function SelectTemplate(
  template: string,
  templateItems: Array<{
    template: string;
    key: string;
    files: Array<{ label: string; path: string }>;
  }>
) {
  const templateFiles = templateItems.find(
    (item) => item.template === template
  );
  return templateFiles ? templateFiles.files : null;
}
export async function SwitchFile(secret: string, addr: string, config: string) {
  const data = new FormData();
  data.append("addr", addr);
  data.append("config", config);
  try {
    const res = await axios.post("/api/exec/update", data, {
      headers: {
        Authorization: secret,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

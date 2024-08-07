import axios from "axios";
export async function FetchConfigs(secret: string) {
  try {
    const response = await axios.get("/api/files/fetch", {
      headers: { Authorization: secret },
    });

    return response.status === 200
      ? Object.entries(response.data).map((files, i) => {
          return {
            template: files[0],
            files: (files[1] as { label: string; path: string }[]).map(
              (file) => {
                return {
                  label: file.label,
                  path: window.location.href + file.path,
                };
              }
            ),
            key: files[0].concat(`-${i}`),
          };
        })
      : null;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

import axios from "axios";
import { SharedSelection } from "@nextui-org/react";

export async function SwitchTemplate(
  secret: string,
  template: string,
  hosts: SharedSelection
) {
  const formdata = new FormData();
  Array.from(hosts).forEach((value) => {
    formdata.append("urls", value.toString().split("-")[0]);
  });
  formdata.append("template", template);
  try {
    const res = await axios.post("/api/host/switch", formdata, {
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

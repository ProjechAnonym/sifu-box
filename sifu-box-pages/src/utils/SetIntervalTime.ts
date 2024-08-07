import axios from "axios";
import { TimeInputValue } from "@nextui-org/react";
function sortTime(time: TimeInputValue | null) {
  return time
    ? [
        time.minute ? time.minute.toString() : "30",
        time.hour ? time.hour.toString() : "4",
      ]
    : ["30", "4"];
}
export async function SetIntervalTime(
  secret: string,
  time: TimeInputValue | null,
  day: number
) {
  const data = new FormData();
  const [minute, hour] = sortTime(time);
  switch (day) {
    case 7:
      data.append("span", minute);
      data.append("span", hour);
      break;
    case 8:
      break;
    default:
      data.append("span", minute);
      data.append("span", hour);
      data.append("span", day.toString());
      break;
  }
  try {
    const res = await axios.post("/api/exec/interval", data, {
      headers: {
        Authorization: `${secret}`,
      },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}

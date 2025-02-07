import { cloneDeep } from "lodash";
import { parseISO, differenceInMinutes } from "date-fns";
import { timeCalc, sizeCalc } from "./calc";
import { logsColumns } from "@/types/singbox/logs";
import { ConnectionColumns, MetaData } from "@/types/singbox/connections";
export function ParseLog(
  data: { type: string; payload: string },
  prevDatas: Array<logsColumns>
) {
  const currentTime = new Date().toTimeString();
  const newDatas = prevDatas.length <= 500 ? prevDatas : prevDatas.slice(1, -1);
  newDatas.push({
    time: currentTime,
    key: `${currentTime}-${
      newDatas[newDatas.length - 1]
        ? Number(newDatas[newDatas.length - 1].key.split("-")[1]) + 1
        : 1
    }`,
    ...data,
  });
  return cloneDeep(newDatas);
}
export function ParseConnection(
  data: Array<{
    chains: Array<string>;
    download: number;
    upload: number;
    id: string;
    rule: string;
    rulePayload: string;
    start: string;
    metadata: MetaData;
  }>,
  prevDatas: Array<ConnectionColumns>
) {
  const lastEleKey = prevDatas[prevDatas.length - 1]
    ? Number(
        prevDatas[prevDatas.length - 1].key.split("-")[
          prevDatas[prevDatas.length - 1].key.split("-").length - 1
        ]
      )
    : 1;
  data.forEach((item, i) => {
    const existData = prevDatas.find((data) => data.id === item.id);
    const startTime = parseISO(item.start);
    const currentTime = new Date();
    const timeDiffMinutes = differenceInMinutes(currentTime, startTime);
    existData
      ? (prevDatas[
          prevDatas.findIndex(
            (existData) => existData.key === `${item.id}-${lastEleKey + 1 + i}`
          )
        ] = {
          ...item.metadata,
          download: sizeCalc(item.download),
          upload: sizeCalc(item.upload),
          id: item.id,
          rule: item.rule,
          rulePayload: item.rulePayload,
          start: timeCalc(timeDiffMinutes),
          chains: item.chains.reverse().join(" -> "),
          key: `${item.id}-${lastEleKey + 1 + i}`,
        })
      : prevDatas.push({
          ...item.metadata,
          download: sizeCalc(item.download),
          upload: sizeCalc(item.upload),
          id: item.id,
          rule: item.rule,
          rulePayload: item.rulePayload,
          start: timeCalc(timeDiffMinutes),
          chains: item.chains.reverse().join(" -> "),
          key: `${item.id}-${lastEleKey + 1 + i}`,
        });
  });

  const deleteDatas: Array<ConnectionColumns> = [];
  const newDatas = prevDatas
    .filter((existData) => {
      if (data.find((item) => item.id === existData.id)) {
        return true;
      } else {
        deleteDatas.push({ ...existData });
        return false;
      }
    })
    .map((existData) => existData);
  const newDeleteDatas =
    deleteDatas.length <= 500
      ? deleteDatas
      : deleteDatas.slice(prevDatas.length - 500);
  return {
    aliveConnections:
      newDatas.length <= 500 ? newDatas : newDatas.slice(newDatas.length - 500),
    deadConnections: newDeleteDatas,
  };
}

import { cloneDeep } from "lodash";
import { parseISO, format, differenceInMinutes } from "date-fns";

import { MetaData, ConnectionLog } from "@/types/singbox";
export function ParseLog(
  data: { type: string; payload: string },
  prevDatas: Array<{ type: string; payload: string; key: string; time: string }>
) {
  const labels = [
    { key: "time", label: "时间", allowSort: true, initShow: true },
    { key: "type", label: "等级", allowSort: true, initShow: true },
    { key: "payload", label: "信息", allowSort: true, initShow: true },
  ];
  const currentTime = new Date().toString();
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
  return { labels, values: cloneDeep(newDatas) };
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
  prevDatas: Array<ConnectionLog>
) {
  const labels = [
    { key: "id", label: "ID", allowSort: true, initShow: false },
    { key: "rule", label: "命中规则", allowSort: true, initShow: true },
    { key: "rulePayload", label: "规则信息", allowSort: true, initShow: false },
    { key: "start", label: "时间", allowSort: true, initShow: true },
    { key: "upload", label: "上传", allowSort: true, initShow: true },
    { key: "download", label: "下载", allowSort: true, initShow: true },
    { key: "chains", label: "节点链", allowSort: true, initShow: true },
    { key: "destinationIP", label: "目标IP", allowSort: true, initShow: true },
    {
      key: "destinationPort",
      label: "目标端口",
      allowSort: true,
      initShow: false,
    },
    {
      key: "dnsMode",
      label: "dns模式",
      allowSort: false,
      initShow: false,
    },
    {
      key: "network",
      label: "网络",
      allowSort: false,
      initShow: false,
    },
    {
      key: "processPath",
      label: "进程",
      allowSort: true,
      initShow: false,
    },
    {
      key: "sourceIP",
      label: "源IP",
      allowSort: true,
      initShow: false,
    },
    {
      key: "sourcePort",
      label: "源端口",
      allowSort: true,
      initShow: false,
    },
    {
      key: "type",
      label: "类型",
      allowSort: false,
      initShow: false,
    },
    { key: "host", label: "嗅探域名", allowSort: true, initShow: true },
  ];
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
    const formatTime = format(startTime, "yyyy-MM-dd HH:mm:ss");
    const currentTime = new Date();
    const timeDiffMinutes = differenceInMinutes(currentTime, formatTime);
    existData
      ? (prevDatas[
          prevDatas.findIndex(
            (existData) => existData.key === `${item.id}-${lastEleKey + 1 + i}`
          )
        ] = {
          ...item.metadata,
          download: SizeCalc(item.download),
          upload: SizeCalc(item.download),
          id: item.id,
          rule: item.rule,
          rulePayload: item.rulePayload,
          start: TimeCalc(timeDiffMinutes),
          chains: item.chains.reverse().join(" -> "),
          key: `${item.id}-${lastEleKey + 1 + i}`,
        })
      : prevDatas.push({
          ...item.metadata,
          download: SizeCalc(item.download),
          upload: SizeCalc(item.download),
          id: item.id,
          rule: item.rule,
          rulePayload: item.rulePayload,
          start: TimeCalc(timeDiffMinutes),
          chains: item.chains.reverse().join(" -> "),
          key: `${item.id}-${lastEleKey + 1 + i}`,
        });
  });
  const deleteDatas: Array<ConnectionLog> = [];
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
      : deleteDatas.slice(prevDatas.length - 500, -1);
  return {
    labels,
    aliveConnections: cloneDeep(
      newDatas.length <= 500
        ? newDatas
        : newDatas.slice(newDatas.length - 500, -1)
    ),
    deadConnections: newDeleteDatas,
  };
}
function SizeCalc(size: number) {
  return size < 1024
    ? `${size} B`
    : size < 1024 * 1024
    ? `${(size / 1024).toFixed(2)} KB`
    : size < 1024 * 1024 * 1024
    ? `${(size / (1024 * 1024)).toFixed(2)} MB`
    : `${(size / (1024 * 1024 * 1024)).toFixed()}`;
}
function TimeCalc(size: number) {
  return size < 60
    ? `${size}分钟`
    : size < 1440
    ? `${(size / 60).toFixed(2)} 小时`
    : size < 60 * 24 * 7
    ? `${(size / (60 * 24)).toFixed(2)}天`
    : `${(size / (60 * 24 * 7)).toFixed(2)}周`;
}

import { useEffect, useState } from "react";
import { Input, Switch } from "@nextui-org/react";
import { RulesetValue } from "@/types/proxy";
export default function AddRuleset(props: {
  onSubmit: (ruleset: RulesetValue) => void;
}) {
  const { onSubmit } = props;
  const [tag, setTag] = useState("");
  const [path, setPath] = useState("");
  const [dnsRule, setDnsRule] = useState("");
  const [group, setGroup] = useState("");
  const [interval, setInterval] = useState("");
  const [detour, setDetour] = useState("");
  const [domestic, setDomestic] = useState(false);
  const [remote, setRemote] = useState(false);
  const [format, setFormat] = useState(false);
  useEffect(() => {
    tag !== "" &&
      path != "" &&
      remote &&
      interval !== "" &&
      detour !== "" &&
      onSubmit({
        id: 0,
        label: group,
        tag,
        url: path,
        path: "",
        update_interval: interval.toString() + "d",
        download_detour: detour,
        china: domestic,
        format: format ? "binary" : "source",
        type: remote ? "remote" : "local",
        dnsRule: dnsRule,
      });
    tag !== "" &&
      path != "" &&
      !remote &&
      onSubmit({
        label: group,
        id: 0,
        tag,
        url: "",
        path,
        type: "local",
        update_interval: "",
        download_detour: "",
        china: domestic,
        format: format ? "binary" : "source",
        dnsRule: dnsRule,
      });
  }, [tag, path, interval, detour, domestic, remote, format]);
  return (
    <div className="flex flex-col gap-y-1">
      <Input
        size="sm"
        label="path"
        isRequired
        isClearable
        onValueChange={setPath}
        validationBehavior="native"
        value={path}
        type={remote ? "url" : "text"}
      />
      <div className="flex flex-row gap-x-2">
        <Input
          className="w-32"
          size="sm"
          label="DNS"
          isClearable
          onValueChange={setDnsRule}
          validationBehavior="native"
          value={dnsRule}
          type={"text"}
        />
        <Input
          className="w-32"
          size="sm"
          label="group"
          isClearable
          onValueChange={setGroup}
          validationBehavior="native"
          value={group}
          type={"text"}
        />
      </div>
      <div className="flex flex-row gap-1">
        <Input
          size="sm"
          className="w-1/3"
          label="label"
          isRequired
          onValueChange={setTag}
          validationBehavior="native"
          value={tag}
        />
        {remote && (
          <Input
            value={interval}
            size="sm"
            className="w-1/3"
            label="interval"
            isRequired
            type="number"
            onValueChange={(e) => setInterval(e.toString())}
            validationBehavior="native"
          />
        )}
        {remote && (
          <Input
            value={detour}
            size="sm"
            className="w-1/3"
            label="detour"
            isRequired
            onValueChange={setDetour}
            validationBehavior="native"
          />
        )}
      </div>
      <div className="flex flex-row justify-start items-center">
        <div className="flex flex-row gap-x-2 items-center">
          <label className="text-foreground font-black">国内</label>
          <Switch
            size="lg"
            onValueChange={(e) => setDomestic(e)}
            startContent={
              <span>
                <i className="bi bi-send-check" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi-send-slash" />
              </span>
            }
          />
        </div>
        <div className="flex flex-row gap-x-2 items-center">
          <label className="text-foreground font-black">类型</label>
          <Switch
            size="lg"
            onValueChange={(e) => setFormat(e)}
            startContent={
              <span>
                <i className="bi bi-file-earmark-binary" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi-filetype-json" />
              </span>
            }
          />
        </div>
        <div className="flex flex-row gap-x-2 items-center">
          <label className="text-foreground font-black">远程</label>
          <Switch
            size="lg"
            onValueChange={(e) => setRemote(e)}
            startContent={
              <span>
                <i className="bi bi-cloud-fill" />
              </span>
            }
            endContent={
              <span>
                <i className="bi bi-cloud-slash-fill" />
              </span>
            }
          />
        </div>
      </div>
    </div>
  );
}

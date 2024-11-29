import { useState, useEffect, useRef } from "react";
import {
  Button,
  Select,
  SelectItem,
  Skeleton,
  Modal,
  ModalContent,
  ModalBody,
  ModalHeader,
  ModalFooter,
  useDisclosure,
  Autocomplete,
  AutocompleteItem,
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@nextui-org/react";
import toast from "react-hot-toast";
import HostAdd from "./addHost";
import { SwitchTemplate } from "@/utils/host/SwitchTemplate";
import { UpgradeApp } from "@/utils/upgrade/upgrade";
import { HostValue } from "@/types/host";
import { SharedSelection } from "@nextui-org/react";
export default function SetHost(props: {
  secret: string;
  hosts: Array<HostValue> | null;
  templates: Array<{ Name: string; Template: Object }> | null;
  setUpdateHosts: (updateHosts: boolean) => void;
  dark: boolean;
}) {
  const { secret, hosts, setUpdateHosts, dark, templates } = props;
  const singboxInput = useRef<HTMLInputElement>(null);
  const [selectHost, setSelectHost] = useState<SharedSelection>(new Set());
  const [selectedTemplate, setSelectedTemplate] = useState<string>("");
  const [errors, setErrors] = useState<Array<string> | null>(null);
  const [submit, setSubmit] = useState(false);
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const clearFiles = () => {
    toast.error("请选择服务器");
    singboxInput.current ? (singboxInput.current.value = "") : null;
  };
  useEffect(() => {
    errors && errors.map((error) => toast.error(error));
  }, [secret, errors]);

  return (
    <div className="flex flex-wrap gap-2 items-center">
      <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
        <ModalContent
          className={`${
            dark ? "sifudark" : "sifulight"
          } bg-background text-foreground`}
        >
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">
                添加主机
              </ModalHeader>
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  setSubmit(true);
                }}
              >
                <ModalBody>
                  <HostAdd
                    submit={submit}
                    setUpdateHosts={setUpdateHosts}
                    secret={secret}
                    setSubmit={setSubmit}
                    onClose={onClose}
                  />
                </ModalBody>
                <ModalFooter className="w-full">
                  <Button color="danger" onPress={onClose} type="button">
                    <span className="text-lg font-black">关闭</span>
                  </Button>
                  <Button color="primary" type="submit">
                    <span className="text-lg font-black">提交</span>
                  </Button>
                </ModalFooter>
              </form>
            </>
          )}
        </ModalContent>
      </Modal>
      {hosts ? (
        <Select
          selectionMode="multiple"
          label="Select a server"
          size="sm"
          radius="sm"
          className="w-56 h-12"
          classNames={{
            popoverContent: `${
              dark ? "sifudark" : "sifulight"
            } bg-default-100 text-foreground`,
          }}
          selectedKeys={selectHost}
          onSelectionChange={setSelectHost}
        >
          {hosts.map((host) => (
            <SelectItem
              key={host.key}
              startContent={
                host.localhost && <i className="bi bi-pc-display-horizontal" />
              }
            >
              {host.url}
            </SelectItem>
          ))}
        </Select>
      ) : (
        <Skeleton className="w-56 h-12 bg-zinc-400 rounded-lg" />
      )}
      <Button color="primary" size="sm">
        <label htmlFor="singboxFile-upload" className="font-black text-md">
          上传sing-box
        </label>
      </Button>
      <input
        type="file"
        id="singboxFile-upload"
        className="hidden"
        ref={singboxInput}
        onChange={(e) => {
          Array.from(selectHost).length !== 0
            ? e.target.files &&
              toast.promise(UpgradeApp(secret, e.target.files, selectHost), {
                loading: "loading",
                success: (res) => {
                  if (singboxInput.current) {
                    singboxInput.current.value = "";
                  }
                  return res ? "升级成功" : "升级失败";
                },
                error: (err) => {
                  err.code != "ERR_NETWORK" &&
                    setErrors(err.response.data.message);
                  if (singboxInput.current) {
                    singboxInput.current.value = "";
                  }
                  return err.code === "ERR_NETWORK"
                    ? "网络错误"
                    : "升级singbox出现错误";
                },
              })
            : clearFiles();
        }}
      />
      {templates ? (
        <Autocomplete
          selectedKey={selectedTemplate}
          onSelectionChange={(e) => e && setSelectedTemplate(e.toString())}
          label={<span className="font-black">模板</span>}
          className="w-36 h-12"
          classNames={{
            popoverContent: `${
              dark ? "sifudark" : "sifulight"
            } bg-default-100 text-foreground`,
          }}
        >
          {templates &&
            templates.map((template) => (
              <AutocompleteItem key={template.Name}>
                {template.Name}
              </AutocompleteItem>
            ))}
        </Autocomplete>
      ) : (
        <Skeleton className="w-36 h-12" />
      )}
      <Popover
        classNames={{
          content: `${
            dark ? "sifudark" : "sifulight"
          } bg-default-100 text-foreground`,
        }}
        placement="bottom"
        showArrow={true}
      >
        <PopoverTrigger>
          <Button color="primary" size="sm">
            <span className="font-black text-md">更改模板</span>
          </Button>
        </PopoverTrigger>
        <PopoverContent>
          <div className="w-44 p-2">
            {Array.from(selectHost).map((value) => (
              <p className="font-black" key={value.toString()}>
                {value.toString().split("-")[0]}
              </p>
            ))}
            <p className="font-black">
              {Array.from(selectHost).length === 0
                ? `请选择要更改模板的主机`
                : `${
                    selectedTemplate
                      ? `模板将更改为"${selectedTemplate}"`
                      : "没有选择要更改的模板"
                  }`}
            </p>
          </div>
          <footer className="flex justify-end w-44 p-2">
            <Button
              size="sm"
              color="primary"
              onPress={() =>
                selectedTemplate && Array.from(selectHost).length !== 0
                  ? toast.promise(
                      SwitchTemplate(secret, selectedTemplate, selectHost),
                      {
                        loading: "loading",
                        success: (res) =>
                          res ? "模板更改成功" : "模板更改失败",
                        error: (err) =>
                          err.code === "ERR_NETWORK"
                            ? "网络错误"
                            : err.response.data.message,
                      }
                    )
                  : toast.error("请选择要更改的模板和主机")
              }
            >
              <span className="font-black text-lg">确认</span>
            </Button>
          </footer>
        </PopoverContent>
      </Popover>
      <Button onPress={onOpen} color="primary" size="sm">
        <span className="font-black text-md">添加主机</span>
      </Button>
    </div>
  );
}

import { useState, useMemo, useRef } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Input } from "@heroui/input";
import { Badge } from "@heroui/badge";
import { Button } from "@heroui/button";
import { Checkbox, CheckboxGroup } from "@heroui/checkbox";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/modal";
import { Divider } from "@heroui/divider";
import { Popover, PopoverTrigger, PopoverContent } from "@heroui/popover";
import { Switch } from "@heroui/switch";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { SetProvider } from "@/utils/application";
import { DeleteItems, AddProviderFiles, AddItems } from "@/utils/configuration";
import { provider } from "@/types/configuration";
export default function Provider(props: {
  theme: string;
  providers: Array<provider>;
  currentProvider: string;
  setCurrentUpdate: (update: boolean) => void;
  setUpdateProviders: (update: boolean) => void;
  token: string;
}) {
  const {
    providers,
    currentProvider,
    token,
    setCurrentUpdate,
    setUpdateProviders,
    theme,
  } = props;
  const fileInput = useRef<HTMLInputElement>(null);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [checked, setChecked] = useState<Array<string>>([]);
  const [search, setSearch] = useState("");
  const [filteredProviders, setFilteredProviders] = useState(providers);
  const [files, setFiles] = useState<FileList | null>(null);

  useMemo(() => {
    providers &&
      search &&
      setFilteredProviders(
        providers.filter((provider) => provider.name.includes(search))
      );
    providers && search === "" && setFilteredProviders(providers);
  }, [search, providers]);
  return (
    <ScrollShadow className="h-1/2 w-full">
      <AddProviders
        isOpen={isOpen}
        onClose={onClose}
        theme={theme}
        setUpdate={setUpdateProviders}
        token={token}
      />
      <header className="flex flex-wrap gap-1 p-1 items-center">
        <Input
          size="sm"
          label={<span className="text-md font-black">机场</span>}
          variant="underlined"
          value={search}
          onValueChange={setSearch}
          className="w-24"
        />
        <Button
          color="danger"
          size="sm"
          variant="shadow"
          onPress={() =>
            toast.promise(DeleteItems(token, checked, "providers"), {
              loading: `删除机场中...`,
              success: (res) => {
                setUpdateProviders(true);
                return res.status ? `删除指定机场成功` : `删除指定机场失败`;
              },
              error: (e) => {
                if (e.code === "ERR_NETWORK") {
                  return "请检查网络连接";
                }
                if (Array.isArray(e.response.data.message)) {
                  (e.response.data.message as Array<string>).map((err) =>
                    toast.error(err)
                  );
                  return "";
                }
                return e.response.data;
              },
            })
          }
        >
          <span className="text-xl font-black">删除</span>
        </Button>
        <Button color="primary" size="sm" variant="shadow" onPress={onOpen}>
          <span className="text-xl font-black">提交</span>
        </Button>
        <form
          className="flex flex-row gap-1"
          onSubmit={(e) => {
            e.preventDefault();
            files
              ? toast.promise(AddProviderFiles(token, files), {
                  loading: "上传Yaml文件...",
                  success: (res) => {
                    res && setUpdateProviders(true);
                    setFiles(null);
                    fileInput.current!.value = "";
                    return "成功添加Yaml文件";
                  },
                  error: (err) => {
                    setFiles(null);
                    fileInput.current!.value = "";
                    setUpdateProviders(true);
                    if (err.code === "ERR_NETWORK") {
                      return "请检查网络连接";
                    }
                    if (Array.isArray(err.response.data.message)) {
                      (err.response.data.message as Array<string>).map((err) =>
                        toast.error(err)
                      );
                      return "";
                    }
                    return `${err.response.data}`;
                  },
                })
              : toast.error("请选择文件");
          }}
        >
          <Badge
            content={files ? files.length : 0}
            color="success"
            placement="bottom-right"
          >
            <label
              htmlFor="file-upload"
              className="px-2 bg-primary h-8 rounded-md hover:cursor-pointer hover:bg-opacity-85 transition-all text-white"
            >
              <i className="bi bi-filetype-yml text-2xl" />
            </label>
          </Badge>
          <input
            type="file"
            id="file-upload"
            className="hidden"
            onChange={(e) => setFiles(e.target.files)}
            multiple
            ref={fileInput}
          />
          <Button size="sm" color="primary" variant="shadow" type="submit">
            <span className="font-black text-xl">上传</span>
          </Button>
        </form>
      </header>
      <CheckboxGroup value={checked} onValueChange={setChecked}>
        <div className={`p-1 flex flex-wrap gap-2`}>
          {filteredProviders && filteredProviders.map((provider) => (
            <Badge
              key={provider.name}
              content={<i className="bi bi-trash-fill" />}
              placement="bottom-left"
              color="danger"
              shape="rectangle"
              className="hover:cursor-pointer"
              onClick={() =>
                toast.promise(
                  DeleteItems(token, [provider.name], "providers"),
                  {
                    loading: `删除${provider.name}机场中...`,
                    success: (res) => {
                      setUpdateProviders(true);
                      return res.status
                        ? `删除${provider.name}机场成功`
                        : `删除${provider.name}机场失败`;
                    },
                    error: (e) => {
                      if (e.code === "ERR_NETWORK") {
                        return "请检查网络连接";
                      }
                      if (Array.isArray(e.response.data.message)) {
                        (e.response.data.message as Array<string>).map((err) =>
                          toast.error(err)
                        );
                        return "";
                      }
                      return e.response.data;
                    },
                  }
                )
              }
            >
              <div className="flex flex-row justify-center h-fit gap-1">
                <Popover
                  classNames={{
                    content: `${theme} bg-content1 text-foreground`,
                  }}
                >
                  <PopoverTrigger>
                    <Button variant="shadow" size="md">
                      <span
                        className={`text-md w-28 text-wrap h-fit font-black select-none ${currentProvider === provider.name && "text-blue-500"}`}
                      >
                        {provider.name}
                      </span>
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent>
                    <p className="text-md font-black w-36 p-1">
                      是否将"{provider.name}"机场设置为活动机场
                    </p>
                    <p className="w-full justify-end flex p-1">
                      <Button
                        size="sm"
                        color="primary"
                        variant="shadow"
                        onPress={() =>
                          toast.promise(SetProvider(token, provider.name), {
                            loading: "设置机场中...",
                            success: (res) => {
                              setCurrentUpdate(true);
                              return res.status
                                ? `设置${provider.name}机场成功`
                                : `设置${provider.name}机场失败`;
                            },
                            error: (e) =>
                              e.code === "ERR_NETWORK"
                                ? "请检查网络连接"
                                : e.response.data.message
                                  ? e.response.data.message
                                  : e.response.data,
                          })
                        }
                      >
                        <span className="text-xl font-black">确认</span>
                      </Button>
                    </p>
                  </PopoverContent>
                </Popover>

                <Checkbox value={provider.name} />
              </div>
            </Badge>
          ))}
        </div>
      </CheckboxGroup>
    </ScrollShadow>
  );
}
function AddProviders(props: {
  isOpen: boolean;
  onClose: () => void;
  theme: string;
  token: string;
  setUpdate: (update: boolean) => void;
}) {
  const { isOpen, onClose, theme, token, setUpdate } = props;
  const [providers, setProviders] = useState<Array<provider>>([
    { name: "", detour: "", remote: false, path: "" },
  ]);
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      backdrop="blur"
      size="xl"
      classNames={{ base: `${theme} bg-content1 text-foreground` }}
    >
      <ModalContent>
        {(onClose) => (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              toast.promise(AddItems(token, providers, "providers"), {
                loading: "添加新机场中...",
                success: (res) => {
                  setUpdate(true);
                  onClose();
                  return res ? "添加成功" : "添加失败";
                },
                error: (e) => {
                  if (e.code === "ERR_NETWORK") {
                    return "请检查网络连接";
                  }
                  if (e.response.data.message) {
                    (e.response.data.message as Array<string>).map((err) =>
                      toast.error(err)
                    );
                    return "";
                  }
                  return e.response.data;
                },
              });
            }}
          >
            <ModalHeader className="gap-1 items-center">
              <span className="font-black text-xl">添加机场</span>
              <Button
                size="sm"
                isIconOnly
                onPress={() => {
                  if (providers.length == 1) {
                    toast.error("无法继续删除");
                    return;
                  }
                  const tempProviders = cloneDeep(providers);
                  tempProviders.pop();
                  setProviders(tempProviders);
                }}
              >
                <i className="bi bi-dash text-3xl" />
              </Button>
              <Button
                size="sm"
                isIconOnly
                onPress={() => {
                  const tempProviders = cloneDeep(providers);
                  tempProviders.push({
                    name: "",
                    detour: "",
                    remote: false,
                    path: "",
                  });
                  setProviders(tempProviders);
                }}
              >
                <i className="bi bi-plus text-3xl" />
              </Button>
            </ModalHeader>
            <ModalBody>
              <ScrollShadow className="flex flex-col gap-4 h-96">
                {providers &&
                  providers.map((provider, i) => (
                    <div key={`provider-${i}`} className="flex flex-col gap-2">
                      <Input
                        variant="flat"
                        size="sm"
                        label={<span className="text-lg font-black">Path</span>}
                        className="w-full"
                        isRequired
                        isClearable
                        value={provider.path}
                        onValueChange={(value) => {
                          const tempProviders = cloneDeep(providers);
                          tempProviders[i].path = value;
                          setProviders(tempProviders);
                        }}
                      />
                      <div className="flex flex-row gap-2 items-center">
                        <Input
                          variant="flat"
                          size="sm"
                          label={
                            <span className="text-lg font-black">名称</span>
                          }
                          className="w-36"
                          isRequired
                          isClearable
                          value={provider.name}
                          onValueChange={(value) => {
                            const tempProviders = cloneDeep(providers);
                            tempProviders[i].name = value;
                            setProviders(tempProviders);
                          }}
                        />
                        <Input
                          variant="flat"
                          size="sm"
                          label={
                            <span className="text-lg font-black">出站</span>
                          }
                          className="w-36"
                          isRequired
                          isClearable
                          value={provider.detour}
                          onValueChange={(value) => {
                            const tempProviders = cloneDeep(providers);
                            tempProviders[i].detour = value;
                            setProviders(tempProviders);
                          }}
                        />
                        <Switch
                          isSelected={provider.remote}
                          onValueChange={(value) => {
                            const tempProviders = cloneDeep(providers);
                            tempProviders[i].remote = value;
                            setProviders(tempProviders);
                          }}
                        >
                          <span className="text-xl font-black">
                            {provider.remote ? "远程" : "本地"}
                          </span>
                        </Switch>
                      </div>
                      <Divider />
                    </div>
                  ))}
              </ScrollShadow>
            </ModalBody>
            <ModalFooter>
              <Button
                size="sm"
                variant="shadow"
                color="danger"
                onPress={onClose}
                type="button"
              >
                <span className="text-lg font-black">关闭</span>
              </Button>
              <Button size="sm" variant="shadow" color="primary" type="submit">
                <span className="text-lg font-black">确认</span>
              </Button>
            </ModalFooter>
          </form>
        )}
      </ModalContent>
    </Modal>
  );
}

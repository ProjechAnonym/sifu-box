import { useState, useMemo } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Input } from "@heroui/input";
import { Badge } from "@heroui/badge";
import { Button } from "@heroui/button";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/modal";
import { Switch } from "@heroui/switch";
import { Tooltip } from "@heroui/tooltip";
import { Divider } from "@heroui/divider";
import { Checkbox, CheckboxGroup } from "@heroui/checkbox";
import toast from "react-hot-toast";
import { cloneDeep } from "lodash";
import { DeleteItems, AddItems } from "@/utils/configuration";
import { ruleset } from "@/types/configuration";
export default function Ruleset(props: {
  rulesets: Array<ruleset>;
  setUpdateRuleset: (update: boolean) => void;
  token: string;
  theme: string;
}) {
  const { rulesets, token, setUpdateRuleset, theme } = props;
  const [checked, setChecked] = useState<Array<string>>([]);
  const [search, setSearch] = useState("");
  const [filteredRulesets, setFilteredRulesets] = useState(rulesets);
  const { isOpen, onOpen, onClose } = useDisclosure();
  useMemo(() => {
    rulesets &&
      search &&
      setFilteredRulesets(
        rulesets.filter((ruleset) => ruleset.tag.includes(search))
      );
    rulesets && search === "" && setFilteredRulesets(rulesets);
  }, [search, rulesets]);
  return (
    <ScrollShadow className="h-1/2 w-full">
      <AddRulesets
        isOpen={isOpen}
        onClose={onClose}
        theme={theme}
        token={token}
        setUpdate={setUpdateRuleset}
      />

      <header className="flex flex-wrap gap-1 p-1 items-center">
        <Input
          size="sm"
          label={<span className="text-md font-black">规则集</span>}
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
            toast.promise(DeleteItems(token, checked, "rulesets"), {
              loading: `删除规则集中...`,
              success: (res) => {
                setUpdateRuleset(true);
                return res.status ? `删除指定规则集成功` : `删除指定规则集失败`;
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
      </header>
      <CheckboxGroup value={checked} onValueChange={setChecked}>
        <div className={`p-1 flex flex-wrap gap-2`}>
          {filteredRulesets.map((ruleset) => (
            <Badge
              key={ruleset.tag}
              content={<i className="bi bi-trash-fill" />}
              placement="bottom-left"
              color="danger"
              shape="rectangle"
              className="hover:cursor-pointer"
              onClick={() =>
                toast.promise(DeleteItems(token, [ruleset.tag], "ruletsets"), {
                  loading: `删除${ruleset.tag}规则集中...`,
                  success: (res) => {
                    setUpdateRuleset(true);
                    return res.status
                      ? `删除${ruleset.tag}规则集成功`
                      : `删除${ruleset.tag}规则集失败`;
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
              <div className="flex flex-row justify-center h-fit gap-1">
                <Button variant="shadow" size="md" isDisabled>
                  <span
                    className={`text-md w-28 text-wrap h-fit font-black select-none`}
                  >
                    {ruleset.tag}
                  </span>
                </Button>
                <Checkbox value={ruleset.tag} />
              </div>
            </Badge>
          ))}
        </div>
      </CheckboxGroup>
    </ScrollShadow>
  );
}
function AddRulesets(props: {
  isOpen: boolean;
  onClose: () => void;
  theme: string;
  token: string;
  setUpdate: (update: boolean) => void;
}) {
  const { isOpen, onClose, theme, token, setUpdate } = props;
  const [rulesets, setRulesets] = useState<Array<ruleset>>([
    {
      tag: "",
      type: "remote",
      path: "",
      format: "binary",
      china: false,
      name_server: "",
      label: "",
      download_detour: "",
      update_interval: "",
    },
  ]);
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      backdrop="blur"
      size="2xl"
      classNames={{ base: `${theme} bg-content1 text-foreground` }}
    >
      <ModalContent>
        {(onClose) => (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              toast.promise(AddItems(token, rulesets, "rulesets"), {
                loading: "添加新规则集中...",
                success: (res) => {
                  setUpdate(true);
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
              <span className="font-black text-xl">添加规则集</span>
              <Button
                size="sm"
                isIconOnly
                onPress={() => {
                  if (rulesets.length == 1) {
                    toast.error("无法继续删除");
                    return;
                  }
                  const tempRulesets = cloneDeep(rulesets);
                  tempRulesets.pop();
                  setRulesets(tempRulesets);
                }}
              >
                <i className="bi bi-dash text-3xl" />
              </Button>
              <Button
                size="sm"
                isIconOnly
                onPress={() => {
                  const tempRulesets = cloneDeep(rulesets);
                  tempRulesets.push({
                    tag: "",
                    type: "remote",
                    path: "",
                    format: "binary",
                    china: false,
                    name_server: "",
                    label: "",
                    download_detour: "",
                    update_interval: "",
                  });
                  setRulesets(tempRulesets);
                }}
              >
                <i className="bi bi-plus text-3xl" />
              </Button>
            </ModalHeader>
            <ModalBody>
              <ScrollShadow className="flex flex-col gap-4 h-96">
                {rulesets &&
                  rulesets.map((ruleset, i) => (
                    <div key={`ruleset-${i}`} className="flex flex-col gap-2">
                      <div className="flex flex-wrap gap-2 items-center">
                        <Input
                          labelPlacement="outside"
                          variant="flat"
                          size="sm"
                          type={ruleset.type === "remote" ? "url" : "text"}
                          label={
                            <span className="text-md font-black">
                              路径
                              <Tooltip
                                placement="top"
                                offset={5}
                                classNames={{
                                  content: [`${theme} bg-content1`],
                                }}
                                content={
                                  <span
                                    className={`${theme} text-foreground font-black w-48`}
                                  >
                                    路径表示规则集的位置,
                                    如果存储按钮选中则表示规则集位于云端,此时仅接受url格式字符串
                                  </span>
                                }
                              >
                                <i className="bi bi-question-circle-fill mx-1" />
                              </Tooltip>
                            </span>
                          }
                          isRequired
                          isClearable
                          value={ruleset.path}
                          onValueChange={(value) => {
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].path = value;
                            setRulesets(tempRulesets);
                          }}
                        />
                        <Input
                          labelPlacement="outside"
                          variant="flat"
                          size="sm"
                          label={
                            <span className="text-md font-black">名称</span>
                          }
                          className="w-36"
                          isRequired
                          isClearable
                          value={ruleset.tag}
                          onValueChange={(value) => {
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].tag = value;
                            setRulesets(tempRulesets);
                          }}
                        />
                        <Input
                          labelPlacement="outside"
                          variant="flat"
                          size="sm"
                          label={
                            <span className="text-md font-black">
                              标签
                              <Tooltip
                                placement="top"
                                offset={5}
                                classNames={{
                                  content: [`${theme} bg-content1`],
                                }}
                                content={
                                  <span
                                    className={`${theme} text-foreground font-black w-48`}
                                  >
                                    标签不同于名称, 标签用于表示该规则集的用途,
                                    比如有网飞IP和网飞域名共同组合的规则集,则这两个规则集会组合为一个网飞的selector出站节点
                                  </span>
                                }
                              >
                                <i className="bi bi-question-circle-fill mx-1" />
                              </Tooltip>
                            </span>
                          }
                          className="w-36"
                          isRequired
                          isClearable
                          value={ruleset.label}
                          onValueChange={(value) => {
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].label = value;
                            setRulesets(tempRulesets);
                          }}
                        />
                        <div className="flex flex-row gap-2">
                          <p className="flex flex-col gap-1">
                            <span className="text-sm font-black">
                              格式
                              <Tooltip
                                placement="top"
                                offset={5}
                                classNames={{
                                  content: [`${theme} bg-content1`],
                                }}
                                content={
                                  <span
                                    className={`${theme} text-foreground font-black w-48`}
                                  >
                                    该按钮用于确定该规则集是用于Json文件还是编译后的二进制文件
                                  </span>
                                }
                              >
                                <i className="bi bi-question-circle-fill mx-1" />
                              </Tooltip>
                            </span>
                            <Switch
                              size="sm"
                              startContent={
                                <span>
                                  <i className="bi bi-filetype-raw" />
                                </span>
                              }
                              endContent={
                                <span>
                                  <i className="bi bi-filetype-json" />
                                </span>
                              }
                              isSelected={ruleset.format === "binary"}
                              onValueChange={(value) => {
                                const tempRulesets = cloneDeep(rulesets);
                                tempRulesets[i].format = value
                                  ? "binary"
                                  : "source";
                                setRulesets(tempRulesets);
                              }}
                            />
                          </p>
                          <p className="flex flex-col gap-1">
                            <span className="text-sm font-black">
                              区域
                              <Tooltip
                                placement="top"
                                offset={5}
                                classNames={{
                                  content: [`${theme} bg-content1`],
                                }}
                                content={
                                  <span
                                    className={`${theme} text-foreground font-black w-48`}
                                  >
                                    该按钮用于确定该规则集是用于国内还是国外
                                  </span>
                                }
                              >
                                <i className="bi bi-question-circle-fill mx-1" />
                              </Tooltip>
                            </span>
                            <Switch
                              size="sm"
                              isSelected={ruleset.china}
                              onValueChange={(value) => {
                                const tempRulesets = cloneDeep(rulesets);
                                tempRulesets[i].china = value;
                                setRulesets(tempRulesets);
                              }}
                              startContent={
                                <span>
                                  <i className="bi bi-bricks" />
                                </span>
                              }
                              endContent={
                                <span>
                                  <i className="bi bi-globe-americas" />
                                </span>
                              }
                            />
                          </p>
                          <Input
                            labelPlacement="outside"
                            variant="flat"
                            size="sm"
                            label={
                              <span className="text-md font-black">
                                NS
                                <Tooltip
                                  placement="top"
                                  offset={5}
                                  classNames={{
                                    content: [`${theme} bg-content1`],
                                  }}
                                  content={
                                    <span
                                      className={`${theme} text-foreground font-black w-48`}
                                    >
                                      NS为用于匹配该规则集时的DNS出站标签
                                    </span>
                                  }
                                >
                                  <i className="bi bi-question-circle-fill mx-1" />
                                </Tooltip>
                              </span>
                            }
                            className="w-32"
                            isClearable
                            value={ruleset.name_server}
                            onValueChange={(value) => {
                              const tempRulesets = cloneDeep(rulesets);
                              tempRulesets[i].name_server = value;
                              setRulesets(tempRulesets);
                            }}
                          />
                        </div>
                        <div className="flex flex-row gap-2">
                          <p className="flex flex-col gap-1">
                            <span className="text-sm font-black">
                              存储
                              <Tooltip
                                placement="top"
                                offset={5}
                                classNames={{
                                  content: [`${theme} bg-content1`],
                                }}
                                content={
                                  <span
                                    className={`${theme} text-foreground font-black w-48`}
                                  >
                                    该按钮用于确定该规则集是位于本地还是云端
                                  </span>
                                }
                              >
                                <i className="bi bi-question-circle-fill mx-1" />
                              </Tooltip>
                            </span>
                            <Switch
                              size="sm"
                              isSelected={ruleset.type === "remote"}
                              startContent={
                                <span>
                                  <i className="bi bi-cloud-arrow-down-fill" />
                                </span>
                              }
                              endContent={
                                <span>
                                  <i className="bi bi-floppy-fill" />
                                </span>
                              }
                              onValueChange={(value) => {
                                const tempRulesets = cloneDeep(rulesets);
                                tempRulesets[i].type = value
                                  ? "remote"
                                  : "local";
                                setRulesets(tempRulesets);
                              }}
                            />
                          </p>
                          {ruleset.type === "remote" && (
                            <Input
                              className="w-16"
                              labelPlacement="outside"
                              variant="flat"
                              size="sm"
                              type="number"
                              label={
                                <span className="text-md font-black">
                                  间隔
                                  <Tooltip
                                    placement="top"
                                    offset={5}
                                    classNames={{
                                      content: [`${theme} bg-content1`],
                                    }}
                                    content={
                                      <span
                                        className={`${theme} text-foreground font-black w-48`}
                                      >
                                        该规则集定期下载的时间间隔
                                      </span>
                                    }
                                  >
                                    <i className="bi bi-question-circle-fill mx-1" />
                                  </Tooltip>
                                </span>
                              }
                              isClearable
                              value={ruleset.update_interval}
                              onValueChange={(value) => {
                                const tempRulesets = cloneDeep(rulesets);
                                tempRulesets[i].update_interval = value;
                                setRulesets(tempRulesets);
                              }}
                            />
                          )}
                          {ruleset.type === "remote" && (
                            <Input
                              className="w-24"
                              labelPlacement="outside"
                              variant="flat"
                              size="sm"
                              label={
                                <span className="text-md font-black">
                                  下载出站
                                  <Tooltip
                                    placement="top"
                                    offset={5}
                                    classNames={{
                                      content: [`${theme} bg-content1`],
                                    }}
                                    content={
                                      <span
                                        className={`${theme} text-foreground font-black w-48`}
                                      >
                                        该规则集定期下载时的出站节点
                                      </span>
                                    }
                                  >
                                    <i className="bi bi-question-circle-fill mx-1" />
                                  </Tooltip>
                                </span>
                              }
                              isClearable
                              value={ruleset.download_detour}
                              onValueChange={(value) => {
                                const tempRulesets = cloneDeep(rulesets);
                                tempRulesets[i].download_detour = value;
                                setRulesets(tempRulesets);
                              }}
                            />
                          )}
                        </div>
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

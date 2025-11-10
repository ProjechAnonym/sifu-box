import { useState, useEffect } from "react";
import { Modal, ModalContent, ModalBody, ModalFooter, ModalHeader } from "@heroui/modal";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { Switch } from "@heroui/switch";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Divider } from "@heroui/divider";
import toast from "react-hot-toast";
import { AddRulesetMsg, EditRuleset } from "@/utils/configuration/ruleset";
import { cloneDeep } from "lodash";
import { DEFAULT_RULESET } from "@/types/setting/ruleset";
export default function AddRuleset(props: {edit:boolean; isOpen: boolean; onClose: () => void; theme: string; token: string; setUpdate: (update: boolean) => void; initial_value: {name: string, path: string, remote: boolean, binary: boolean, download_detour: string, update_interval: string}}) {
  const { isOpen, onClose, theme, token, setUpdate, initial_value, edit } = props;
  const [rulesets, setRulesets] = useState<Array<{name: string, path: string, remote: boolean, binary: boolean, download_detour: string, update_interval: string}>>([initial_value])
  useEffect(() => {
    setRulesets([initial_value]);
  }, [initial_value]);
  const AddItems = () => toast.promise(
    AddRulesetMsg(token, rulesets),
    {
      loading: "正在添加...",
      success: (res) => {
        res ? res.map(item => item.status ? toast.success(item.message) : toast.error(item.message)) : toast.error("添加失败, 未知错误")
        setUpdate(true);
        onClose();
        return "添加操作完成";
      },
      error:(e) => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                e.response.data.message ? e.response.data.message : e.response.data
    }
  );
  const EditItems = () => toast.promise(
    EditRuleset(token, rulesets[0]),
    {
      loading: "正在修改...",
      success: (res) => {
        res ? res.message ? toast.success(res.message) : toast.error("未知错误") : toast.error("修改失败, 未知错误")
        setUpdate(true);
        onClose();
        return "修改操作完成";
      },
      error:(e) => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                e.response.data.message ? e.response.data.message : e.response.data
    }
  )
  return (
    <Modal isOpen={isOpen} onClose={onClose} backdrop="blur" size="xl" classNames={{ base: `${theme} bg-content1 text-foreground` }}>
      <ModalContent>
        {(onClose) => (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              edit ? EditItems() : AddItems();
            }}
          >
            <ModalHeader className="gap-2 items-center">
              <span className="font-black text-xl">添加规则集</span>
              {!edit &&
              <div className="flex flex-row gap-1">
                <Button size="sm" isIconOnly 
                  onPress={() => {
                    if (rulesets.length == 1) {
                      toast.error("无法继续删除");
                      return;
                    }
                    const temp_rulesets = cloneDeep(rulesets);
                    temp_rulesets.pop();
                    setRulesets(temp_rulesets);
                }}>
                  <i className="bi bi-dash text-3xl" />
                </Button>
                <Button size="sm" isIconOnly 
                  onPress={() => {
                  const temp_rulesets = cloneDeep(rulesets);
                    temp_rulesets.push(DEFAULT_RULESET);
                    setRulesets(temp_rulesets);
                  }}
                >
                  <i className="bi bi-plus text-3xl" />
                </Button>
              </div>}
            </ModalHeader>
            <ModalBody>
              <ScrollShadow className="flex flex-col gap-4 h-96">
                {rulesets &&
                  rulesets.map((ruleset, i) => (
                    <div key={`ruleset-${i}`} className="flex flex-wrap gap-2">
                      <Input
                        variant="flat"
                        size="sm"
                        label={<span className="text-lg font-black">Path</span>}
                        className="w-full"
                        isRequired
                        isClearable
                        value={ruleset.path}
                        onValueChange={(value) => {
                          const temp_rulesets = cloneDeep(rulesets);
                          temp_rulesets[i].path = value;
                          setRulesets(temp_rulesets);
                        }}
                      />
                      {ruleset.remote &&
                        <div className="flex flex-row gap-2">
                          <Input
                            variant="flat"
                            size="sm"
                            label={<span className="text-lg font-black">Detour</span>}
                            className="w-1/2"
                            isRequired
                            isClearable
                            value={ruleset.download_detour}
                            onValueChange={(value) => {
                              const temp_rulesets = cloneDeep(rulesets);
                              temp_rulesets[i].download_detour = value;
                              setRulesets(temp_rulesets);
                            }}
                          />
                          <Input
                            variant="flat"
                            size="sm"
                            label={<span className="text-lg font-black">Interval</span>}
                            className="w-1/2"
                            isRequired
                            isClearable
                            value={ruleset.update_interval}
                            onValueChange={(value) => {
                              const temp_rulesets  = cloneDeep(rulesets);
                              temp_rulesets[i].update_interval = value;
                              setRulesets(temp_rulesets);
                            }}
                          />
                        </div>
                      }
                      <div className="flex flex-row gap-2 items-center">
                        <Input
                          variant="flat"
                          size="sm"
                          label={
                            <span className="text-lg font-black">名称</span>
                          }
                          className="w-32"
                          isRequired
                          isClearable
                          value={ruleset.name}
                          onValueChange={(value) => {
                            const temp_rulesets = cloneDeep(rulesets);
                            temp_rulesets[i].name = value;
                            setRulesets(temp_rulesets);
                          }}
                        />
                        <Switch
                          isSelected={ruleset.remote}
                          onValueChange={(value) => {
                            const temp_rulesets = cloneDeep(rulesets);
                            temp_rulesets[i].remote = value;
                            setRulesets(temp_rulesets);
                          }}
                          size="sm"
                        >
                          <span className="text-sm font-black">
                            {ruleset.remote ? "远程" : "本地"}
                          </span>
                        </Switch>
                        <Switch
                          isSelected={!ruleset.binary}
                          onValueChange={(value) => {
                            const temp_rulesets = cloneDeep(rulesets);
                            temp_rulesets[i].binary = !value;
                            setRulesets(temp_rulesets);
                          }}
                          size="sm"
                        >
                          <span className="text-sm font-black">
                            {ruleset.binary ? "Raw" : "Json"}
                          </span>
                        </Switch>
                      </div>
                      <Divider />
                    </div>
                  ))}
              </ScrollShadow>
            </ModalBody>
            <ModalFooter>
              <Button size="sm" variant="shadow" color="danger" onPress={()=>{setRulesets([initial_value]); onClose()}} type="button">
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
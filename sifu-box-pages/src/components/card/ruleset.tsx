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
export default function AddRuleset(props: {edit:boolean; isOpen: boolean; onClose: () => void; theme: string; token: string; setUpdate: (update: boolean) => void; initial_value: {name: string, path: string, remote: boolean}}) {
  const { isOpen, onClose, theme, token, setUpdate, initial_value, edit } = props;
  const [rulesets, setRulesets] = useState<Array<{name: string, path: string, remote: boolean}>>([initial_value])
  useEffect(() => {
    setRulesets([initial_value]);
  }, [initial_value]);
  // const addItems = () => toast.promise(
  //   AddProviderMsg(token, edit, providers),
  //   {
  //     loading: "正在添加...",
  //     success: (res) => {
  //       res ? res.map(item => item.status ? toast.success(item.message) : toast.error(item.message)) : toast.error("添加失败, 未知错误")
  //       setUpdate(true);
  //       onClose();
  //       return "添加操作完成";
  //     },
  //     error:(e) => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
  //               e.response.data.message ? e.response.data.message : e.response.data
  //   }
  // );
  // const editItems = () => toast.promise(
  //   EditProvider(token, providers[0]),
  //   {
  //     loading: "正在修改...",
  //     success: (res) => {
  //       res ? res.message ? toast.success(res.message) : toast.error("未知错误") : toast.error("修改失败, 未知错误")
  //       setUpdate(true);
  //       onClose();
  //       return "修改操作完成";
  //     },
  //     error:(e) => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
  //               e.response.data.message ? e.response.data.message : e.response.data
  //   }
  // )
  return (
    <Modal isOpen={isOpen} onClose={onClose} backdrop="blur" size="xl" classNames={{ base: `${theme} bg-content1 text-foreground` }}>
      <ModalContent>
        {(onClose) => (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              // edit ? editItems() : addItems();
            }}
          >
            <ModalHeader className="gap-2 items-center">
              <span className="font-black text-xl">添加机场</span>
              {!edit &&
              <div className="flex flex-row gap-1">
                <Button size="sm" isIconOnly 
                  onPress={() => {
                    if (rulesets.length == 1) {
                      toast.error("无法继续删除");
                      return;
                    }
                    const tempRulesets = cloneDeep(rulesets);
                    tempRulesets.pop();
                    setRulesets(tempRulesets);
                }}>
                  <i className="bi bi-dash text-3xl" />
                </Button>
                <Button size="sm" isIconOnly 
                  onPress={() => {
                  const tempRulesets = cloneDeep(rulesets);
                    tempRulesets.push({
                      name: "",
                      remote: false,
                      path: "",
                    });
                    setRulesets(tempRulesets);
                  }}
                >
                  <i className="bi bi-plus text-3xl" />
                </Button>
              </div>}
            </ModalHeader>
            <ModalBody>
              <ScrollShadow className="flex flex-col gap-4 h-96">
                {rulesets &&
                  rulesets.map((provider, i) => (
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
                          const tempRulesets  = cloneDeep(rulesets);
                          tempRulesets [i].path = value;
                          setRulesets (tempRulesets);
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
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].name = value;
                            setRulesets(tempRulesets);
                          }}
                        />
                        <Switch
                          isSelected={provider.remote}
                          onValueChange={(value) => {
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].remote = value;
                            setRulesets(tempRulesets);
                          }}
                        >
                          <span className="text-xl font-black">
                            {provider.remote ? "远程" : "本地"}
                          </span>
                        </Switch>
                        <Switch
                          isSelected={provider.remote}
                          onValueChange={(value) => {
                            const tempRulesets = cloneDeep(rulesets);
                            tempRulesets[i].remote = value;
                            setRulesets(tempRulesets);
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
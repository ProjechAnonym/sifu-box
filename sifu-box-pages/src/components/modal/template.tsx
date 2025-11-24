import { useState, useMemo } from "react";
import { Modal, ModalBody, ModalContent, ModalFooter, ModalHeader } from "@heroui/modal";
import { Button } from "@heroui/button";
import { Input, Textarea } from "@heroui/input";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { CheckboxGroup, Checkbox } from "@heroui/checkbox";
import {Breadcrumbs, BreadcrumbItem} from "@heroui/breadcrumbs";
import OutboundsGroup from "../outboundsGroup";
import toast from "react-hot-toast";
import { AddTemplateMsg, EditTemplateMsg } from "@/utils/configuration/template";
import { Provider } from "@/types/setting/provider";
import { RuleSet } from "@/types/setting/ruleset";
export default function AddTemplate(props: {
  isOpen: boolean;
  onClose: () => void;
  edit: boolean;
  token: string;
  template: {name: string, [key: string]: any};
  theme: string;
  setUpdate: (update: boolean) => void;
  providers:Array<Provider>;
  rulesets:Array<RuleSet>,
}) {
  const { isOpen, onClose, edit, token, template, theme, setUpdate, providers, rulesets } = props;
  const [content, setContent] = useState("");
  const [name, setName] = useState("");
  const [current_modal, setCurrentModal] = useState("providers");
  const [selected_providers, setSelectedProviders] = useState<Array<string>>([]);
  const [selected_rulesets, setSelectedRulesets] = useState<Array<string>>([]);
  const [outbounds_groups, setOutboundsGroups] = useState<Array<{type: string, tag: string, providers: string[], tag_groups: string[]}>>([{type: "direct", tag: "direct", providers: [], tag_groups: []}]);
  useMemo(() => {
    template && setContent(JSON.stringify(template, null, 4));
    setSelectedProviders(template.providers ?  template.providers : [])
    setSelectedRulesets(template.route.rule_set ? template.route.rule_set.map((rule_set: Record<string, any>) => rule_set.tag) : [])
    setOutboundsGroups(template.outbounds_group ? template.outbounds_group : [{type: "direct", tag: "direct", providers: [], tag_groups: []}])
  }, [template]);
  const AddItem = (content: string, name: string) => {
    try {
      const template_msg = JSON.parse(content);
      template_msg.name = name
      template_msg.providers = selected_providers;
      const ruleset_list = rulesets.filter(ruleset => selected_rulesets.includes(ruleset.name)).map(ruleset => {
        const ruleset_msg: Record<string, any>  = {tag: ruleset.name, type: ruleset.remote ? "remote" : "local", format: ruleset.binary ? "binary" : "source"}
        ruleset.remote ? ruleset_msg.url = ruleset.path : ruleset_msg.path = ruleset.path
        ruleset.remote ? ruleset_msg.update_interval = ruleset.update_interval : null
        ruleset.remote ? ruleset_msg.download_detour = ruleset.download_detour : null
        return ruleset_msg
      })
      template_msg.route.rule_set = ruleset_list
      template_msg.outbounds_group = outbounds_groups
      content && toast.promise(AddTemplateMsg(token, template_msg),
        {
          loading: "正在添加...",
          success: (res) => {
            setUpdate(true);
            return res
          },
          error: e => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
              e.response.data.message ? e.response.data.message : e.response.data
        }
      )
      
    } catch (e) {
      console.error(e)
      toast.error("模板格式错误");
    }
  }
  const EditItem = (content: string, name: string) => {
    try {
      const template_msg = JSON.parse(content);
      template_msg.name = name
      template_msg.providers = selected_providers;
      const ruleset_list = rulesets.filter(ruleset => selected_rulesets.includes(ruleset.name)).map(ruleset => {
        const ruleset_msg: Record<string, any>  = {tag: ruleset.name, type: ruleset.remote ? "remote" : "local", format: ruleset.binary ? "binary" : "source"}
        ruleset.remote ? ruleset_msg.url = ruleset.path : ruleset_msg.path = ruleset.path
        ruleset.remote ? ruleset_msg.update_interval = ruleset.update_interval : null
        ruleset.remote ? ruleset_msg.download_detour = ruleset.download_detour : null
        return ruleset_msg
      })
      template_msg.route.rule_set = ruleset_list
      template_msg.outbounds_group = outbounds_groups
      content && toast.promise(EditTemplateMsg(token, template_msg),
        {
          loading: "正在修改...",
          success: (res) => {
            setUpdate(true);
            return res
          },
          error: e => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
              e.response.data.message ? e.response.data.message : e.response.data
        }
      )
      
    } catch (e) {
      console.error(e)
      toast.error("模板格式错误");
    }
  }

  return (
    <Modal isOpen={isOpen} backdrop="blur" onClose={onClose} size="3xl" classNames={{ base: `${theme} bg-content1 text-foreground h-96` }}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader>
              {edit ? (template.name) : (<Input className="w-24" size="sm" variant="underlined" label={<span className="font-black">模板名称</span>} value={name} onValueChange={setName}/>)}
            </ModalHeader>
            <ModalBody>
              <Breadcrumbs onAction={(key) => typeof key === "string" && setCurrentModal(key)}>
                <BreadcrumbItem key="providers" isCurrent={current_modal === "providers"}>
                  选择机场
                </BreadcrumbItem>
                <BreadcrumbItem key="rulesets" isCurrent={current_modal === "rulesets"}>
                  选择规则集
                </BreadcrumbItem>
                <BreadcrumbItem key="outbounds_group" isCurrent={current_modal === "outbounds_group"}>
                  设置出站集
                </BreadcrumbItem>
                <BreadcrumbItem key="template" isCurrent={current_modal === "template"}>
                  编写模板
                </BreadcrumbItem>
              </Breadcrumbs>
              {current_modal === "providers" && (
                <ScrollShadow className="w-full h-full">
                  <CheckboxGroup value={selected_providers} onValueChange={setSelectedProviders}>
                    <div className="w-full flex flex-wrap gap-2 h-full">
                      {providers.map((provider) => (
                        <div key={provider.name} className="border-1.5 p-1 rounded-lg">
                          <Checkbox value={provider.name}>
                            {provider.name}
                          </Checkbox>
                        </div>
                      ))}
                    </div>
                  </CheckboxGroup>
                </ScrollShadow>)}
              {current_modal === "rulesets" && (
                <ScrollShadow className="w-full h-full">
                  <CheckboxGroup value={selected_rulesets} onValueChange={setSelectedRulesets}>
                    <div className="flex flex-wrap gap-2 h-full">
                      {rulesets.map((ruleset) => (
                        <div key={ruleset.name} className="border-1.5 p-1 rounded-lg">
                          <Checkbox value={ruleset.name}>
                            {ruleset.name}
                          </Checkbox>
                        </div>
                      ))}
                    </div>
                  </CheckboxGroup>
                </ScrollShadow>)}
              {current_modal === "template" && (<Textarea label="模板内容" value={content} onValueChange={setContent}/>)}
              {current_modal === "outbounds_group" && <OutboundsGroup theme={theme} providers={selected_providers} outbounds_group={outbounds_groups} setOutboundsGroup={setOutboundsGroups}/>}
            </ModalBody>
            <ModalFooter>
              <Button size="sm" color="danger" variant="shadow" onPress={onClose}>
                <span className="font-black text-xl">关闭</span>
              </Button>
              <Button size="sm" color="danger" variant="shadow" onPress={() => template && setContent(JSON.stringify(template, null, 4))}>
                <span className="font-black text-xl">恢复</span>
              </Button>
              <Button size="sm" color="primary" variant="shadow" onPress={() => edit ? EditItem(content, template.name): AddItem(content, name)}>
                <span className="font-black text-xl">提交</span>
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
import { useState, useMemo } from "react";
import { Modal, ModalBody, ModalContent, ModalFooter, ModalHeader } from "@heroui/modal";
import { Button } from "@heroui/button";
import { Input, Textarea } from "@heroui/input";
import toast from "react-hot-toast";
import { AddTemplateMsg, EditTemplateMsg } from "@/utils/configuration/template";
export default function AddTemplate(props: {
  isOpen: boolean;
  onClose: () => void;
  edit: boolean;
  token: string;
  template: {name: string, [key: string]: any};
  theme: string;
  setUpdate: (update: boolean) => void;
}) {
  const { isOpen, onClose, edit, token, template, theme, setUpdate } =
    props;
  const [content, setContent] = useState("");
  const [name, setName] = useState("");
  useMemo(() => {
    template && setContent(JSON.stringify(template, null, 4));
  }, [template]);
  const addItem = (content: string, name: string) => {
    try {
      const template_msg = JSON.parse(content);
      template_msg.name = name
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
  const editItem = (content: string, name: string) => {
    try {
      const template_msg = JSON.parse(content);
      template_msg.name = name
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
    <Modal
      isOpen={isOpen}
      backdrop="blur"
      onClose={onClose}
      size="3xl"
      classNames={{ base: `${theme} bg-content1 text-foreground` }}
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader>
              {edit ? (
                template.name
              ) : (
                <Input
                  className="w-24"
                  size="sm"
                  variant="underlined"
                  label={<span className="font-black">模板名称</span>}
                  value={name}
                  onValueChange={setName}
                />
              )}
            </ModalHeader>
            <ModalBody>
              <Textarea
                label="模板内容"
                value={content}
                onValueChange={setContent}
              />
            </ModalBody>
            <ModalFooter>
              <Button
                size="sm"
                color="danger"
                variant="shadow"
                onPress={onClose}
              >
                <span className="font-black text-xl">关闭</span>
              </Button>
              <Button
                size="sm"
                color="danger"
                variant="shadow"
                onPress={() =>
                  template && setContent(JSON.stringify(template, null, 4))
                }
              >
                <span className="font-black text-xl">恢复</span>
              </Button>
              <Button
                size="sm"
                color="primary"
                variant="shadow"
                onPress={() => edit ? editItem(content, template.name): addItem(content, name)}
              >
                <span className="font-black text-xl">提交</span>
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
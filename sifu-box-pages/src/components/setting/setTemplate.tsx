import { useState, useMemo } from "react";
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Textarea,
  Button,
  Input,
} from "@nextui-org/react";
import toast from "react-hot-toast";
import { ModifyTemplate } from "@/utils/template/ModifyTemplate";
export default function SetTemplate(props: {
  secret: string;
  dark: boolean;
  isOpen: boolean;
  template: { Name: string; Template: Object } | null;
  newTemplate: boolean;
  setUpdateTemplate: (updateTemplate: boolean) => void;
  onClose: () => void;
}) {
  const {
    secret,
    isOpen,
    onClose,
    template,
    dark,
    newTemplate,
    setUpdateTemplate,
  } = props;
  const [name, setName] = useState("");
  const [content, setContent] = useState("");
  useMemo(() => {
    template && setContent(JSON.stringify(template.Template, null, 4));
    template && setName(template.Name);
  }, [template]);
  return (
    <Modal isOpen={isOpen} backdrop="blur" onClose={onClose} size="3xl">
      <ModalContent
        className={`${
          dark ? "sifudark" : "sifulight"
        } bg-background text-foreground`}
      >
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-row gap-x-6 items-center">
              {template ? template.Name : "模板失效"}
              {newTemplate && (
                <Input
                  className="w-28"
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
                value={content}
                onValueChange={setContent}
                label="模板内容"
              />
            </ModalBody>
            <ModalFooter>
              <Button size="sm" color="danger" onPress={onClose}>
                <span className="font-black text-lg">关闭</span>
              </Button>
              <Button
                size="sm"
                color="danger"
                onPress={() =>
                  template &&
                  setContent(JSON.stringify(template.Template, null, 4))
                }
              >
                <span className="font-black text-lg">恢复</span>
              </Button>
              <Button
                size="sm"
                color="primary"
                onPress={() => {
                  toast.promise(ModifyTemplate(secret, name, content), {
                    loading: "loading",
                    success: (res) => {
                      setUpdateTemplate(true);
                      return res
                        ? `设置"${name}"模板成功`
                        : `设置"${name}"模板失败`;
                    },
                    error: (err) =>
                      err.code === "ERR_NETWORK"
                        ? "网络错误"
                        : err.response.data.message,
                  });
                }}
              >
                <span className="font-black text-lg">提交</span>
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}

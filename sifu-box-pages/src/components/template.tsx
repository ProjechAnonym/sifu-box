import { useState, useMemo, useEffect } from "react";
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
} from "@heroui/modal";
import { CheckboxGroup, Checkbox } from "@heroui/checkbox";
import { Card, CardBody, CardFooter, CardHeader } from "@heroui/card";
import { Button } from "@heroui/button";
import { Input, Textarea } from "@heroui/input";
import { ScrollShadow } from "@heroui/scroll-shadow";
import toast from "react-hot-toast";
import { SetTemplate } from "@/utils/application";
import { DeleteItems, ModifyTemplate } from "@/utils/configuration";
export default function Template(props: {
  token: string;
  template?: { [key: string]: Object } | null;
  setUpdate: (update: boolean) => void;
  defaultTemplate: Object | null;
  theme: string;
}) {
  const { template, token, setUpdate, defaultTemplate, theme } = props;
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [edit, setEdit] = useState(false);
  const [templateName, setTemplateName] = useState("");
  const [editTemplate, setEditTemplate] = useState<Object | null>(null);
  const [selected, setSelected] = useState<Array<string>>([]);
  const [search, setSearch] = useState<string>("");
  const [errors, setErrors] = useState<Array<string>>([]);
  const [displayTemplate, setDisplayTemplate] = useState<{
    [key: string]: Object;
  } | null>(null);
  useEffect(() => {
    errors && errors.map((error) => toast.error(error));
    setErrors([]);
  }, [errors]);
  useMemo(() => {
    template &&
      search &&
      setDisplayTemplate(
        Object.entries(template!)
          .filter(([key]) => key.includes(search))
          .reduce(
            (acc, [key, value]) => {
              acc[key] = value;
              return acc;
            },
            {} as { [key: string]: Object }
          )
      );
    template && !search && setDisplayTemplate(template);
  }, [search, template]);
  return (
    <div>
      <EditTemplate
        isOpen={isOpen}
        onClose={onClose}
        edit={edit}
        token={token}
        template={editTemplate}
        name={templateName}
        theme={theme}
        setUpdate={(update) => setUpdate(update)}
      />
      <header className="flex gap-2 flex-row items-center p-2">
        <Input
          size="sm"
          variant="underlined"
          className="w-56"
          type="text"
          label={<span className="font-black">模板名称</span>}
          isClearable
          value={search}
          onValueChange={setSearch}
        />
        <Button
          size="sm"
          color="primary"
          variant="shadow"
          onPress={() => {
            setEdit(false);
            defaultTemplate && setEditTemplate(defaultTemplate);
            onOpen();
          }}
        >
          <span className="font-black text-xl">添加</span>
        </Button>
        <Button
          size="sm"
          color="danger"
          variant="shadow"
          onPress={() =>
            toast.promise(DeleteItems(token, selected, "templates"), {
              loading: `删除选中模板中...`,
              success: (res) => {
                setUpdate(true);
                return res.status ? `删除选中模板成功` : `删除选中模板失败`;
              },
              error: (e) => {
                if (e.code === "ERR_NETWORK") {
                  return "请检查网络连接";
                }
                if (e.response.data.message) {
                  setErrors(e.response.data.message);
                  return "更新配置文件失败";
                }
                return e.response.data;
              },
            })
          }
        >
          <span className="font-black text-xl">删除</span>
        </Button>
      </header>
      <CheckboxGroup value={selected} onValueChange={setSelected}>
        <ScrollShadow
          className="w-full flex flex-row gap-4 p-2"
          orientation="horizontal"
        >
          {displayTemplate &&
            Object.entries(displayTemplate!).map(([key, value]) => (
              <div key={key} className="w-72">
                <Card>
                  <CardHeader className="justify-between">
                    <span
                      className="text-xl font-black hover:cursor-pointer hover:bg-gray-500 transition-all p-1 rounded-md"
                      onClick={() =>
                        toast.promise(SetTemplate(token, key), {
                          loading: "设置模板中...",
                          success: (res) =>
                            res.status
                              ? `设置${key}模板成功`
                              : `设置${key}模板失败`,
                          error: (e) =>
                            e.code === "ERR_NETWORK"
                              ? "请检查网络连接"
                              : e.response.data.message
                                ? e.response.data.message
                                : e.response.data,
                        })
                      }
                    >
                      {key}
                    </span>
                    <Checkbox value={key} />
                  </CardHeader>
                  <CardBody>
                    <ScrollShadow className="w-full h-52">
                      {JSON.stringify(value, null, 4)}
                    </ScrollShadow>
                  </CardBody>
                  <CardFooter className="justify-end gap-2">
                    <Button
                      color="danger"
                      size="sm"
                      variant="shadow"
                      onPress={() =>
                        toast.promise(DeleteItems(token, [key], "templates"), {
                          loading: `删除${key}模板中...`,
                          success: (res) => {
                            setUpdate(true);
                            return res.status
                              ? `删除${key}模板成功`
                              : `删除${key}模板失败`;
                          },
                          error: (e) => {
                            if (e.code === "ERR_NETWORK") {
                              return "请检查网络连接";
                            }
                            if (e.response.data.message) {
                              setErrors(e.response.data.message);
                              return "更新配置文件失败";
                            }
                            return e.response.data;
                          },
                        })
                      }
                    >
                      <span className="text-xl font-black">删除</span>
                    </Button>
                    <Button
                      color="primary"
                      size="sm"
                      variant="shadow"
                      onPress={() => {
                        setEdit(true);
                        setEditTemplate(value);
                        setTemplateName(key);
                        onOpen();
                      }}
                    >
                      <span className="text-xl font-black">修改</span>
                    </Button>
                  </CardFooter>
                </Card>
              </div>
            ))}
        </ScrollShadow>
      </CheckboxGroup>
    </div>
  );
}

function EditTemplate(props: {
  isOpen: boolean;
  onClose: () => void;
  edit: boolean;
  token: string;
  template: Object | null;
  name: string;
  theme: string;
  setUpdate: (update: boolean) => void;
}) {
  const { isOpen, onClose, edit, token, template, name, theme, setUpdate } =
    props;
  const [content, setContent] = useState("");
  const [errors, setErrors] = useState<Array<string>>([]);
  const [templateName, setTemplateName] = useState("");
  useEffect(() => {
    errors && errors.map((error) => toast.error(error));
    edit && setTemplateName(name);
  }, [errors, edit]);
  useMemo(() => {
    template && setContent(JSON.stringify(template, null, 4));
  }, [template]);
  return (
    <Modal isOpen={isOpen} backdrop="blur" onClose={onClose} size="3xl">
      <ModalContent className={`${theme} bg-background text-foreground`}>
        {(onClose) => (
          <>
            <ModalHeader>
              {edit ? (
                templateName
              ) : (
                <Input
                  className="w-24"
                  size="sm"
                  variant="underlined"
                  label={<span className="font-black">模板名称</span>}
                  value={templateName}
                  onValueChange={setTemplateName}
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
                onPress={() => {
                  try {
                    JSON.parse(content);
                  } catch (e) {
                    console.error(e);
                    toast.error("模板内容解析失败");
                    return;
                  }
                  templateName === ""
                    ? toast.error("没有设置模板名称")
                    : toast.promise(
                        ModifyTemplate(token, templateName, content),
                        {
                          loading: "loading",
                          success: (res) => {
                            setUpdate(true);
                            return res
                              ? `设置"${templateName}"模板成功`
                              : `设置"${templateName}"模板失败`;
                          },
                          error: (err) => {
                            if (err.code === "ERR_NETWORK") {
                              return "请检查网络连接";
                            }
                            setUpdate(true);
                            if (err.response.data.message) {
                              setErrors(err.response.data.message);
                              return "更新配置文件失败";
                            }
                            return err.response.data;
                          },
                        }
                      );
                }}
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
